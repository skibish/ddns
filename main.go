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

var (
	reqTimeouts = flag.Duration("req-timeout", 10*time.Second, "Request timeout to external resources")
	checkPeriod = flag.Duration("check-period", 5*time.Minute, "Check if IP has been changed period")
	confFile    = flag.String("conf-file", "$HOME/.ddns.yml", "Location of the configuration file")
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func main() {
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
	var provider *ipprovider.IPProvider
	if cf.IPv6 {
		provider = ipprovider.New()
		provider.Register(
			ifconfig.New(c, true),
			ipify.New(c, true),
		)
	} else {
		provider = ipprovider.New()
		provider.Register(
			ifconfig.New(c, false),
			ipify.New(c, false),
			wtfismyip.New(c),
		)
	}

	// Initialize and start updater
	upd, errUpdater := updater.New(c, provider, cf, *checkPeriod)
	if errUpdater != nil {
		log.Fatal(errUpdater)
	}

	errStart := upd.Start()
	if errStart != nil {
		log.Fatal(errStart)
	}

	select {}
}
