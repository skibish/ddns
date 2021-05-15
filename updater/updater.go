package updater

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/skibish/ddns/conf"

	log "github.com/sirupsen/logrus"
	"github.com/skibish/ddns/do"
	"github.com/skibish/ddns/ipprovider"
)

//go:generate moq -out do_moq_test.go -pkg updater ../do DomainsService
//go:generate moq -out ipprovider_moq_test.go -pkg updater ../ipprovider Provider

// Updater is responsible for DNS records updates.
type Updater struct {
	ip         string
	ticker     *time.Ticker
	do         do.DomainsService
	ipprovider ipprovider.Provider
	config     *conf.Configuration
	shutdown   chan bool
}

// New return new Updater.
func New(cfg *conf.Configuration) *Updater {
	return &Updater{
		ticker:     time.NewTicker(cfg.CheckPeriod),
		do:         do.New(cfg.Token, cfg.RequestTimeout),
		ipprovider: ipprovider.New(cfg.IPv6, cfg.RequestTimeout),
		shutdown:   make(chan bool),
		config:     cfg,
	}
}

// Start starts the updater process.
func (u *Updater) Start(ctx context.Context) (err error) {
	log.Debug("initializing ip")

	if _, err := u.ipUpdated(ctx); err != nil {
		return fmt.Errorf("failed to initialize ip: %w", err)
	}

	log.Infof("current ip is %s", u.ip)

	log.Debug("syncing dns records")
	if err := u.sync(ctx); err != nil {
		return fmt.Errorf("failed to sync dns records: %w", err)
	}
	log.Debug("done")

	// perform IP checks in intervals
	for {
		select {
		case <-u.ticker.C:
			log.Debugf("checking if ip (%s) has been updated", u.ip)

			updated, err := u.ipUpdated(ctx)
			if err != nil {
				return fmt.Errorf("failed to get ip: %w", err)
			}

			if !updated {
				continue
			}
			log.Infof("ip has been updated to %s", u.ip)

			log.Debug("updating dns records", u.ip)
			if err := u.sync(ctx); err != nil {
				return fmt.Errorf("failed to update dns records: %w", err)
			}
			log.Debug("done")
		case <-u.shutdown:
			return nil
		}
	}
}

// Stops stops the updater.
func (u *Updater) Stop() {
	u.ticker.Stop()
	u.shutdown <- true
}

// ipUpdated returns true and updates IP to a new value if IP changed.
func (u *Updater) ipUpdated(ctx context.Context) (bool, error) {
	newIP, err := u.ipprovider.GetIP(ctx)
	if err != nil {
		return false, err
	}

	if u.ip == newIP {
		return false, nil
	}

	u.ip = newIP

	return true, nil
}

// match checks if records are the same
func match(a, b do.Record) bool {
	return a.Type == b.Type && a.Name == b.Name
}

// search searches for a record in records.
// If success, returns record ID which is not 0.
func (u *Updater) search(records []do.Record, record do.Record) uint64 {
	for _, r := range records {
		if match(record, r) {
			return r.ID
		}
	}

	return 0
}

// sync syncs DNS records.
func (u *Updater) sync(ctx context.Context) error {
	for domain := range u.config.Domains {
		records, err := u.do.List(ctx, domain)
		if err != nil {
			return fmt.Errorf("failed to get the records for the domain %s: %w", domain, err)
		}

		for _, r := range u.config.Domains[domain] {
			r.Data, err = u.prepareData(r, u.config.Params)
			if err != nil {
				return fmt.Errorf("failed to set data to the record %s of the domain %s: %w", domain, r.Type, err)
			}

			recordID := u.search(records, r)
			if recordID == 0 {
				if err := u.do.Create(ctx, domain, r); err != nil {
					return fmt.Errorf("failed to create a record for the domain %s: %w", domain, err)
				}
				continue
			}

			r.ID = recordID
			if err := u.do.Update(ctx, domain, r); err != nil {
				return fmt.Errorf("failed to update a record for the domain %s: %w", domain, err)
			}
		}
	}

	return nil
}

// prepareData executes template and return what should be set in the DNS record data field.
// It can be just an IP or some string.
func (u *Updater) prepareData(configRecord do.Record, params map[string]string) (string, error) {
	if configRecord.Data == "" {
		return u.ip, nil
	}

	params["IP"] = u.ip

	t, err := template.New("t1").Parse(configRecord.Data)
	if err != nil {
		return "", fmt.Errorf("failed to parse the template: %w", err)
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, params); err != nil {
		return "", fmt.Errorf("failed to execute a template: %w", err)
	}

	return buf.String(), nil
}
