package main

import (
	"os"
	go_syslog "log/syslog"
	"net/http"
	"github.com/sirupsen/logrus"
	rus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

const (
//	url = "http://gfechio-k8s-onboarding-team2.stress.thebeat.co"
	url = "http://127.0.0.1/"
)

var logger = logrus.New()

func logg() {
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.InfoLevel
	logger.SetReportCaller(true)

	hook, err := rus_syslog.NewSyslogHook("", "", go_syslog.LOG_INFO, "concurMQ")
	if err == nil {
		logger.Hooks.Add(hook)
	} else {
		logger.Error("Syslog is not available, using concurMQ.log instead")
		f, err := os.OpenFile("concurMQ.log", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			logger.Fatal("Can't fallback to file.")
		}
		logrus.SetOutput(f)
	}
}

func main() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		logger.Error("Could not set HTTP request")
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Could not make request")
	}
	logger.Info("Response value as : %v", resp)
}
