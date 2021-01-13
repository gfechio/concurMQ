package main

import (
	go_syslog "log/syslog"
	"net/http"
	"os"
	"runtime"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	rus_syslog "github.com/sirupsen/logrus/hooks/syslog"
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

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Get the function calling us if possible.
	var caller string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		caller = details.Name()
	}

	_, err := w.Write(append(response, byte(0x0A)))
	if err == nil {
		if caller == "" {
			logger.Info("Response: " + string(response))
		} else {
			logger.WithField("handler", caller).Info("Response: " + string(response))
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)

		if caller == "" {
			logger.Fatal("Preparing response: " + err.Error())
		} else {
			logger.WithField("handler", caller).Fatal("Preparing response: " + err.Error())
		}
	}
}

func management(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := "OK"
	respondWithJSON(w, http.StatusOK, response)
}

func swagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	http.ServeFile(w, r, "swagger.yml")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, `{"healthy": "true"}`)
}

func main() {
	broker()
	router := mux.NewRouter()
	// Just for testing (and as an example to start from).
	router.HandleFunc("/healthz", healthz).Methods("GET")
	// Endpoints.
	router.HandleFunc("/", management).Methods("GET")

	// Web server, this exits only upon errors.
	logger.Fatal(http.ListenAndServe(":5000", router))
	
}
