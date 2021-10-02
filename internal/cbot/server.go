package cbot

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func (b *cbot) serve() error {
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 5050),
		ErrorLog:     b.errorLog,
		Handler:      b.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	b.infoLog.Println("Starting server on :5000")
	return srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
}
