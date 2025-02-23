package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/skibish/ddns/updater"

	log "github.com/sirupsen/logrus"
	"github.com/skibish/ddns/conf"
	"github.com/skibish/ddns/notifier"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(stdout)

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		confFile    = flags.String("conf-file", "", "Location of the configuration file. If not provided, searches current directory, then $HOME for ddns.yml file")
		showVersion = flags.Bool("ver", false, "Show version")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if *showVersion {
		fmt.Printf("Version: %s\nCommit: %s\nBuild date: %s\n", version, commit, date)
		return nil
	}

	cf, err := conf.NewConfiguration(*confFile)
	if err != nil {
		log.Fatal(err)
	}

	// try to register all provided hooks
	for _, v := range cf.Notifications {
		hook, err := notifier.GetHook(v)
		if err != nil {
			log.Debugf("failed to add a notifier: %s", err)
			continue
		}
		log.AddHook(hook)
	}

	upd := updater.New(cf)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ctx.Done()
		log.Debug("shutdown")
		upd.Stop()
		stop()
		os.Exit(0)
	}()

	return upd.Start(ctx)
}
