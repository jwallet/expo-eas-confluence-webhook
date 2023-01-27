package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

func setupServer() {
	log.Println("Starting server...")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("Check"))
		fmt.Println("healthcheck")
	})

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// h := hmac.New(sha1.New, []byte(secret))
		// h.Write([]byte(r.Body))
		// fmt.Print(h.Size())
		// webhookHandler(r.Body)

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var expoContext ExpoBuild
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		json.Unmarshal(body, &expoContext)

		// if !hmac.Equal([]byte(r.Header.Get("Expo-Signature")), h.Sum(nil)) {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }

		if expoContext.Status == "finished" {
			log.Println("Received build")
			err := webhookHandler(expoContext)
			if err != nil {
				fmt.Printf("An error occured %v\n", err)
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

	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server listening on localhost:8080")
	log.Fatal(http.Serve(l, nil))
}
