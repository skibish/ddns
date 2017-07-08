package main

import (
	"bytes"
	"flag"
	"html/template"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skibish/ddns/conf"
	"github.com/skibish/ddns/do"
	"github.com/skibish/ddns/ipprovider"
	"github.com/skibish/ddns/notifier"
)

var (
	digio   do.DigitalOceanInterface
	cf      *conf.Configuration
	storage *conf.Configuration
	periodC <-chan time.Time
)

var (
	reqTimeouts = flag.Duration("req-timeout", 10*time.Second, "Request timeout to external resources")
	checkPeriod = flag.Duration("check-period", 5*time.Minute, "Check if IP has been changed period")
	confFile    = flag.String("conf-file", "$HOME/.ddns.yml", "Location of the configuration file")
)

// current remembered IP
var currentIP string

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	// initialize storage
	storage = &conf.Configuration{}
}

func main() {
	flag.Parse()

	// read configuration
	var errConf error
	cf, errConf = conf.NewConfiguration(*confFile)
	if errConf != nil {
		log.Fatal(errConf.Error())
	}

	// try to register all provided hooks
	for k, v := range cf.Notify {
		hook, errGet := notifier.GetHook(k, v)
		if errGet != nil {
			log.Debugf("Notifier %q not added: %s", k, errGet.Error())
			continue
		}
		log.AddHook(hook)
	}

	// setup http client
	hc := &http.Client{
		Timeout: *reqTimeouts,
	}

	// initialte digital ocean client
	digio = do.NewDigitalOcean(cf.Domain, cf.Token, hc)

	// register providers
	ipprovider.Register(hc)

	// get current IP
	currentIP = ipprovider.GetIP()
	if currentIP == "" {
		log.Fatal("IP can't be empty in the beginning... Do you have internet connection?")
	}
	log.Infof("Current IP is %q", currentIP)

	// do request to the digital ocean API for list of records
	allRecords, errGetDR := digio.GetDomainRecords()
	if errGetDR != nil {
		log.Fatal(errGetDR.Error())
	}

	// do initial sync of records
	var errSync error
	errSync = syncRecords(storage, cf, allRecords)
	if errSync != nil {
		log.Fatal(errSync.Error())
	}

	periodC = time.NewTicker(*checkPeriod).C

	// start main proceess
	go func(storage *conf.Configuration) {
		// for defined period of time, perform IP check
		for {
			select {
			case <-periodC:
				errCheck := checkAndUpdate(storage, cf, ipprovider.GetIP)
				if errCheck != nil {
					log.Errorf("Failed to update: %s", errCheck.Error())
				}
			}
		}

	}(storage)

	select {}
}

// syncRecords perform initial sync between what we provided
// in configuration and what already exist in DNS records
func syncRecords(storage *conf.Configuration, cf *conf.Configuration, allRecords []do.Record) error {
	cRec := len(storage.Records)
	cAllRec := len(allRecords)
	for i := 0; i < cRec; i++ {
		for j := 0; j < cAllRec; j++ {

			// we are only interested in those who have full match
			// by `type AND name`
			if storage.Records[i].Type == allRecords[j].Type &&
				storage.Records[i].Name == allRecords[j].Name {
				storage.Records[i] = allRecords[j]
				break
			}
		}

		// if there was no match, we should create new DNS record
		// and updatee current configuration
		if storage.Records[i].ID == 0 {
			// if there is not template in configuration, set current IP as data,
			// otherwise parse data and fill template with provided params
			if cf.Records[i].Data == "" {
				storage.Records[i].Data = currentIP
			} else {
				storage.Params["IP"] = currentIP
				t := template.Must(template.New("t1").Parse(cf.Records[i].Data))
				buf := new(bytes.Buffer)
				errExec := t.Execute(buf, storage.Params)
				if errExec != nil {
					return errExec
				}
				storage.Records[i].Data = buf.String()
			}

			newR, errCreate := digio.CreateRecord(storage.Records[i])
			if errCreate != nil {
				return errCreate
			}

			storage.Records[i] = *newR
		}

		// if IPs are different, update record
		if storage.Records[i].Data != currentIP {
			if cf.Records[i].Data == "" {
				storage.Records[i].Data = currentIP
			} else {
				storage.Params["IP"] = currentIP
				t := template.Must(template.New("t1").Parse(cf.Records[i].Data))
				buf := new(bytes.Buffer)
				errExec := t.Execute(buf, storage.Params)
				if errExec != nil {
					return errExec
				}
				storage.Records[i].Data = buf.String()
			}

			newR, errUpdate := digio.UpdateRecord(storage.Records[i])
			if errUpdate != nil {
				return errUpdate
			}

			storage.Records[i] = *newR
		}
	}

	return nil
}

// checkAndUpdate check for new IP and if it has been changed,
// trigger the update of the DNS records
func checkAndUpdate(storage *conf.Configuration, cf *conf.Configuration, getIP ipprovider.FGetIP) error {
	log.Debug("IP check")
	newIP := getIP()

	if currentIP != newIP {
		log.Infof("IP has changed from %q to %q", currentIP, newIP)
		currentIP = newIP

		cRec := len(storage.Records)
		for i := 0; i < cRec; i++ {
			if cf.Records[i].Data == "" {
				storage.Records[i].Data = currentIP
			} else {
				storage.Params["IP"] = currentIP
				t := template.Must(template.New("t1").Parse(cf.Records[i].Data))
				buf := new(bytes.Buffer)
				errExec := t.Execute(buf, storage.Params)
				if errExec != nil {
					return errExec
				}
				storage.Records[i].Data = buf.String()
			}

			newR, errUpdate := digio.UpdateRecord(storage.Records[i])
			if errUpdate != nil {
				return errUpdate
			}

			storage.Records[i] = *newR
		}
	}

	return nil
}
