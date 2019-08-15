package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/skibish/ddns/ipprovider/ifconfig"
	"github.com/skibish/ddns/ipprovider/ipify"
	"github.com/skibish/ddns/ipprovider/wtfismyip"
	"github.com/skibish/ddns/updater"

	log "github.com/sirupsen/logrus"
	"github.com/skibish/ddns/conf"
	"github.com/skibish/ddns/ipprovider"
	"github.com/skibish/ddns/notifier"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	var (
		reqTimeouts = flag.Duration("req-timeout", 10*time.Second, "Request timeout to external resources")
		checkPeriod = flag.Duration("check-period", 5*time.Minute, "Check if IP has been changed period")
		confFile    = flag.String("conf-file", "$HOME/.ddns.yml", "Location of the configuration file")
	)
	flag.Parse()

	// read configuration
	var errConf error
	cf, errConf := conf.NewConfiguration(*confFile)
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
	c := &http.Client{
		Timeout: *reqTimeouts,
	}

	// initialize new ipprovider and register IP providers
	provider := ipprovider.New()

	providerList := []ipprovider.Provider{
		ifconfig.New(c),
		wtfismyip.New(c),
		ipify.New(c),
	}

	if cf.ForceIPV6 {
		for _, p := range providerList {
			p.ForceIPV6()
		}
	}

	provider.Register(providerList...)

	// Initialize and start updaters
	for _, domain := range cf.Domains {
		upd, errUpdater := updater.New(c, provider, cf, domain, *checkPeriod)
		if errUpdater != nil {
			log.Fatal(errUpdater)
		}

		errStart := upd.Start()
		if errStart != nil {
			log.Fatal(errStart)
		}
	}

	select {}
}
