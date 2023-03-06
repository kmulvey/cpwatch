package main

import (
	"net/http"
	"time"

	"encoding/json"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func webServer(addr string, db *gorm.DB) {

	http.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		var ds, err = getLatest(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("error getting latest status, err: " + err.Error()))
			if err != nil {
				log.Errorf("error sending http error: %s", err.Error())
			}
		}
		js, err := json.Marshal(ds)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("error marshalling latest status to json, err: " + err.Error()))
			if err != nil {
				log.Errorf("error sending http error: %s", err.Error())
			}
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(js)
		if err != nil {
			log.Errorf("error sending http response, error: %s", err.Error())
		}
	})

	log.Infof("Server is running on port: %s", addr)

	var server = &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("http server error: ", err)
	}
}
