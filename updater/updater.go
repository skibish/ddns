package updater

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/mitchellh/copystructure"
	"github.com/skibish/ddns/conf"

	log "github.com/sirupsen/logrus"
	"github.com/skibish/ddns/do"
	"github.com/skibish/ddns/ipprovider"
)

// Updater is responsible for updating DNS records if IP has changed
type Updater struct {
	ip           string
	updateTick   time.Duration
	digitalOcean do.DigitalOceanInterface
	ipprovider   *ipprovider.IPProvider
	storage      *conf.Configuration
	config       *conf.Configuration
}

// New return new Updater.
func New(hc *http.Client, ipprovider *ipprovider.IPProvider, cfg *conf.Configuration, updateTick time.Duration) (u *Updater, err error) {

	u = &Updater{
		updateTick:   updateTick,
		digitalOcean: do.New(cfg.Domain, cfg.Token, hc),
		ipprovider:   ipprovider,
	}

	// configuration and storage have same structure for simplicity
	copyForStorage, err := copystructure.Copy(cfg)
	if err != nil {
		return
	}

	var ok bool
	u.storage, ok = copyForStorage.(*conf.Configuration)
	if !ok {
		return nil, errors.New("Failed to convert interface{} to conf.Configuration")
	}
	u.config = cfg

	return
}

// Start starts the updater process goroutine
func (u *Updater) Start() (err error) {
	// get current IP
	u.ip = u.ipprovider.GetIP()
	if u.ip == "" {
		return errors.New("IP can't be empty in the beginning... Do you have internet connection?")
	}
	log.Infof("Current IP is %q", u.ip)

	// do request to the digital ocean API for list of records
	allRecords, err := u.digitalOcean.GetDomainRecords()
	if err != nil {
		return err
	}

	// do initial sync of records
	err = u.syncRecords(allRecords)
	if err != nil {
		return err
	}

	periodC := time.NewTicker(u.updateTick).C

	// start main proceess
	go func() {
		// for defined period of time, perform IP check
		for {
			select {
			case <-periodC:
				errCheck := u.checkAndUpdate()
				if errCheck != nil {
					log.Errorf("failed to update: %s", errCheck.Error())
				}
			}
		}

	}()

	return
}

// syncRecords perform initial sync between what we provided
// in configuration and what already exist in DNS records
func (u *Updater) syncRecords(allRecords []do.Record) error {
	cRec := len(u.storage.Records)
	cAllRec := len(allRecords)
	for i := 0; i < cRec; i++ {
		for j := 0; j < cAllRec; j++ {

			// we are only interested in those who have full match
			// by `type AND name`
			if u.storage.Records[i].Type == allRecords[j].Type &&
				u.storage.Records[i].Name == allRecords[j].Name {
				u.storage.Records[i] = allRecords[j]
				break
			}
		}

		// if there was no match, we should create new DNS record
		// and updatee current configuration
		if u.storage.Records[i].ID == 0 {
			// if there is not template in configuration, set current IP as data,
			// otherwise parse data and fill template with provided params
			errUpdStorage := u.updateStorage(&u.storage.Records[i], &u.config.Records[i], u.config.Params)
			if errUpdStorage != nil {
				return errUpdStorage
			}

			newR, errCreate := u.digitalOcean.CreateRecord(u.storage.Records[i])
			if errCreate != nil {
				return errCreate
			}

			u.storage.Records[i] = *newR
		}

		// if IPs are different, update record
		if u.storage.Records[i].Data != u.ip {
			errUpdStorage := u.updateStorage(&u.storage.Records[i], &u.config.Records[i], u.config.Params)
			if errUpdStorage != nil {
				return errUpdStorage
			}

			newR, errUpdate := u.digitalOcean.UpdateRecord(u.storage.Records[i])
			if errUpdate != nil {
				return errUpdate
			}

			u.storage.Records[i] = *newR
		}
	}

	return nil
}

// checkAndUpdate check for new IP and if it has been changed,
// trigger the update of the DNS records
func (u *Updater) checkAndUpdate() error {
	log.Debug("IP check")
	newIP := u.ipprovider.GetIP()

	if u.ip != newIP {
		log.Infof("IP has changed from %q to %q", u.ip, newIP)
		u.ip = newIP

		cRec := len(u.storage.Records)
		for i := 0; i < cRec; i++ {
			errUpdStorage := u.updateStorage(&u.storage.Records[i], &u.config.Records[i], u.storage.Params)
			if errUpdStorage != nil {
				return errUpdStorage
			}

			newR, errUpdate := u.digitalOcean.UpdateRecord(u.storage.Records[i])
			if errUpdate != nil {
				return errUpdate
			}

			u.storage.Records[i] = *newR
		}
	}

	return nil
}

// updateStorage updates the storage based on data in configuration
func (u *Updater) updateStorage(storageRecord, configRecord *do.Record, params map[string]string) (err error) {
	if configRecord.Data == "" {
		storageRecord.Data = u.ip
	} else {
		params["IP"] = u.ip
		t := template.Must(template.New("t1").Parse(configRecord.Data))
		buf := new(bytes.Buffer)
		err = t.Execute(buf, params)
		if err != nil {
			return
		}
		storageRecord.Data = buf.String()
	}
	return
}
