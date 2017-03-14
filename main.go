package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/skibish/ddns/conf"
	"github.com/skibish/ddns/do"
	"github.com/skibish/ddns/ipprovider"
)

var (
	digio   do.DigitalOceanInterface
	cf      *conf.Configuration
	periodC <-chan time.Time
)

var (
	reqTimeouts = flag.Duration("req-timeout", 1*time.Second, "Request timeout to external resources")
	checkPeriod = flag.Duration("check-period", 5*time.Minute, "Check if IP has been changed period")
	confFile    = flag.String("conf-file", "$HOME/.ddns.yml", "Location of the configuration file")
)

// current remembered IP
var currentIP string

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()

	// read configuration
	var errConf error
	cf, errConf = conf.NewConfiguration(*confFile)
	if errConf != nil {
		log.Fatal(errConf.Error())
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
	errSync = syncRecords(cf, allRecords)
	if errSync != nil {
		log.Fatal(errSync.Error())
	}

	periodC = time.NewTicker(*checkPeriod).C

	// start main proceess
	go func(cf *conf.Configuration) {
		// for defined period of time, perform IP check
		for {
			select {
			case <-periodC:
				errCheck := checkAndUpdate(cf, ipprovider.GetIP)
				if errCheck != nil {
					log.Errorf("Failed to update: %s", errCheck.Error())
				}
			}
		}

	}(cf)

	select {}
}

// syncRecords perform initial sync between what we provided
// in configuration and what already exist in DNS records
func syncRecords(cf *conf.Configuration, allRecords []do.Record) error {
	cRec := len(cf.Records)
	cAllRec := len(allRecords)
	for i := 0; i < cRec; i++ {
		for j := 0; j < cAllRec; j++ {

			// we are only interested in those who have full match
			// by `type AND name`
			if cf.Records[i].Type == allRecords[j].Type &&
				cf.Records[i].Name == allRecords[j].Name {
				cf.Records[i] = allRecords[j]
				break
			}
		}

		// if there was no match, we should create new DNS record
		// and updatee current configuration
		if cf.Records[i].ID == 0 {
			cf.Records[i].Data = currentIP

			newR, errCreate := digio.CreateRecord(cf.Records[i])
			if errCreate != nil {
				return errCreate
			}

			cf.Records[i] = *newR
		}

		// if IPs are different, update record
		if cf.Records[i].Data != currentIP {
			cf.Records[i].Data = currentIP

			newR, errUpdate := digio.UpdateRecord(cf.Records[i])
			if errUpdate != nil {
				return errUpdate
			}

			cf.Records[i] = *newR
		}
	}

	return nil
}

// checkAndUpdate check for new IP and if it has been changed,
// trigger the update of the DNS records
func checkAndUpdate(cf *conf.Configuration, getIP ipprovider.FGetIP) error {
	log.Info("IP check")
	newIP := getIP()

	if currentIP != newIP {
		log.Infof("IP has changed from %q to %q", currentIP, newIP)
		currentIP = newIP

		cRec := len(cf.Records)
		for i := 0; i < cRec; i++ {
			cf.Records[i].Data = currentIP

			newR, errUpdate := digio.UpdateRecord(cf.Records[i])
			if errUpdate != nil {
				return errUpdate
			}

			cf.Records[i] = *newR
		}
	}

	return nil
}
