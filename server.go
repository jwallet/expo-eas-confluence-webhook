package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func SetupServer() {
	log.Println("Setting up server")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("OK"))
		log.Println("healthcheck")
	})

	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/text")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var expoContext ExpoBuild
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		digest := hmac.New(sha1.New, []byte(EXPO_HMAC_SECRET))
		digest.Write(body)
		receivedSignature := r.Header.Get("expo-signature")
		log.Printf("Received signature: %v\n", receivedSignature)
		expectedSignature := fmt.Sprintf("sha1=%v", hex.EncodeToString(digest.Sum(nil)))
		if expectedSignature != receivedSignature {
			log.Printf("Invalid HMAC, received %v, expected %v\n", receivedSignature, expectedSignature)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		json.Unmarshal(body, &expoContext)

		if expoContext.Status == "finished" {
			log.Println("Received build")
			err := Webhook(expoContext)
			if err != nil {
				log.Printf("An error occured: %v\n", err)
				w.WriteHeader(http.StatusConflict)
			} else {
				log.Println("Expo build published on Confluence")
				w.WriteHeader(http.StatusOK)
			}
		} else {
			log.Println("Received unfinished build")
			w.WriteHeader(http.StatusOK)
		}
	})

	mux.HandleFunc("/inject", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/text")
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Println("Injecting a build manually")

		var build Build
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		json.Unmarshal(body, &build)

		err = PushBuild(build)
		if err != nil {
			log.Printf("An error occured %v\n", err)
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusOK)
		}

	})

	mux.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/text")
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := Init()
		if err != nil {
			log.Printf("An error occured %v\n", err)
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	log.Println("Starting server...")

	// Determine port for HTTP service.
	port := PORT
	log.Printf("Using port %v\n", port)
	if port == 0 {
		port = 8080
		log.Printf("Defaulting to port %v\n", port)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}

	server.SetKeepAlivesEnabled(false)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server listening on localhost:%v\n", port)
}
