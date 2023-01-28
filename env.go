package main

import (
	"log"
	"os"
	"strconv"
)

var EXPO_HMAC_SECRET = os.Getenv("EXPO_HMAC_SECRET")

var PORT int = getEnvInt("PORT")

var CONFLUENCE_CLOUD_DOMAIN = os.Getenv("CONFLUENCE_CLOUD_DOMAIN")
var CONFLUENCE_TOKEN = os.Getenv("CONFLUENCE_TOKEN")
var CONFLUENCE_USER = os.Getenv("CONFLUENCE_USER")
var CONFLUENCE_PAGE_ID int = getEnvInt("CONFLUENCE_PAGE_ID")

func getEnvInt(key string) int {
	var variable string = os.Getenv(key)
	if len(variable) == 0 {
		variable = "0"
	}
	val, err := strconv.Atoi(variable)
	if err != nil {
		log.Fatal("Failed to convert env var to integer")
	}
	return val
}

func logEnvs() {
	log.Printf(`Environment variables
	- PORT: %v
	- EXPO_HMAC_SECRET: %v
	- CONFLUENCE_CLOUD_DOMAIN: %v
	- CONFLUENCE_TOKEN: %v
	- CONFLUENCE_USER: %v
	- CONFLUENCE_PAGE_ID: %v`,
		PORT,
		EXPO_HMAC_SECRET,
		CONFLUENCE_CLOUD_DOMAIN,
		CONFLUENCE_TOKEN,
		CONFLUENCE_USER,
		CONFLUENCE_PAGE_ID)
}
