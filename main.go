package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"go.szostok.io/version"
	"go.szostok.io/version/printer"
)

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	var sigChannel = make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	// get user opts
	var cmdPath string
	var pollInterval time.Duration
	var v bool
	flag.StringVar(&cmdPath, "cmd-path", "/usr/sbin/pwrstat", "absolute path to pwstat command")
	flag.DurationVar(&pollInterval, "poll-interval", time.Minute, "time interval to gather power stats")
	flag.BoolVar(&v, "version", false, "print version")
	flag.BoolVar(&v, "v", false, "print version")

	flag.Parse()

	if v {
		var verPrinter = printer.New()
		var info = version.Get()
		if err := verPrinter.PrintInfo(os.Stdout, info); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	var db, err = dbConnect("cp.db")
	if err != nil {
		log.Fatal(err)
	}
	go webServer(":8080", db)

	var ticker = time.NewTicker(pollInterval)
	for {
		select {
		case <-ticker.C:
			log.Info("gathering stats")

			out, err := getPowerStats(cmdPath)
			if err != nil {
				log.Fatal(err)
			}

			status, err := parsePowerStats(out)
			if err != nil {
				log.Fatal(err)
			}

			if err := insert(db, status); err != nil {
				log.Fatal(err)
			}
		case <-sigChannel:
			log.Info("shutting down")
			return
		}
	}
}
