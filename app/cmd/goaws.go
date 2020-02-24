
package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Freman/eventloghook"
	"github.com/kardianos/service"
	"github.com/p4tin/goaws/app/common"
	"github.com/p4tin/goaws/app/conf"
	"github.com/p4tin/goaws/app/gosqs"
	"github.com/p4tin/goaws/app/router"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc/eventlog"
)

var version string

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	var filename string
	var debug bool
	flag.StringVar(&filename, "config", "", "config file location + name")
	flag.BoolVar(&debug, "debug", false, "debug log level (default Warning)")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	e, err := eventlog.Open("AbsorbLMS.GoAws.Service")
	if err != nil {
		panic(err)
	}
	defer e.Close()

	hook := eventloghook.NewHook(e)
	log.AddHook(hook)

	env := "Local"
	if flag.NArg() > 0 {
		env = flag.Arg(0)
	}

	portNumbers := conf.LoadYamlConfig(filename, env)

	if common.LogMessages {
		file, err := os.OpenFile(common.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			mw := io.MultiWriter(os.Stdout, file)
			log.SetOutput(mw)
		} else {
			log.Warn("Failed to log to file, using default stderr")
		}
	}

	r := router.New()

	quit := make(chan struct{}, 0)
	go gosqs.PeriodicTasks(1*time.Second, quit)

	if len(portNumbers) == 1 {
		log.Infof("GoAws listening on: 0.0.0.0:%s", portNumbers[0])
		err := http.ListenAndServe("0.0.0.0:"+portNumbers[0], r)
		log.Fatal(err)
	} else if len(portNumbers) == 2 {
		go func() {
			log.Infof("GoAws listening on: 0.0.0.0:%s", portNumbers[0])
			err := http.ListenAndServe("0.0.0.0:"+portNumbers[0], r)
			log.Fatal(err)
		}()
		log.Infof("GoAws listening on: 0.0.0.0:%s", portNumbers[1])
		err := http.ListenAndServe("0.0.0.0:"+portNumbers[1], r)
		log.Fatal(err)
	} else {
		log.Fatal("Not enough or too many ports defined to start GoAws.")
	}
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	//by default, windows services start in C:\WINDOWS\System32 folder, so we need to change the working directory
	exePath, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	os.Chdir(filepath.Dir(exePath))

	svcConfig := &service.Config{
		Name:        "AbsorbLMS.GoAws.Service",
		DisplayName: "Absorb GoAws Service",
		Description: "This windows service simulates AWS SNS and SQS services using GoAws project(https://github.com/p4tin/goaws)",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Error(err)
	}
}
