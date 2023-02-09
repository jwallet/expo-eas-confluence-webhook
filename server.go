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
	"strings"
)

func setupServer() {
	log.Println("Starting server...")

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("OK"))
		fmt.Println("healthcheck")
	})

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Received signature: %v", receivedSignature)
		expectedSignature := hex.EncodeToString(digest.Sum(nil))
		if expectedSignature != receivedSignature {
			log.Printf("Invalid HMAC, received %v, expected %v", receivedSignature, expectedSignature)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		json.Unmarshal(body, &expoContext)

		if expoContext.Status == "finished" {
			log.Println("Received build")
			err := webhookHandler(expoContext)
			if err != nil {
				fmt.Printf("An error occured: %v\n", err)
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

	http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/text")
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Println("Init confluence page")

		previousPage, err := getConfluencePage(CONFLUENCE_PAGE_ID)
		if err != nil {
			fmt.Printf("An error occured %v\n", err)
			w.WriteHeader(http.StatusForbidden)
		}

		template := getDefaultTemplate()
		minifier := strings.NewReplacer("\n", "", "\t", "")
		minifiedTemplate := minifier.Replace(template)
		page := generateConfluenceUpdatePagePayload(previousPage, "Init EAS builds template", minifiedTemplate)

		err = putConfluencePage(CONFLUENCE_PAGE_ID, page)
		if err != nil {
			fmt.Printf("An error occured %v\n", err)
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	// Determine port for HTTP service.
	port := PORT
	fmt.Print(port)
	if port == 0 {
		port = 8080
		log.Printf("Defaulting to port %v\n", port)
	}

	log.Printf("Server listening on localhost:%v\n", port)
	if err := http.ListenAndServe(":"+fmt.Sprint(port), nil); err != nil {
		log.Fatal(err)
	}
}
