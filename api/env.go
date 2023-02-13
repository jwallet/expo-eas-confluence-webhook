package main

import (
	"log"
	"os"
	"strconv"
)

func GetEnvInt(key string) int {
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

func LogEnvs() {
	log.Printf(`Environment variables
	- PORT: %v
	- EXPO_HMAC_SECRET: %v
	- CONFLUENCE_CLOUD_DOMAIN: %v
	- CONFLUENCE_TOKEN: %v
	- CONFLUENCE_USER: %v
	- CONFLUENCE_SPACE: %v
	- CONFLUENCE_PAGE_ID: %v`,
		PORT,
		EXPO_HMAC_SECRET,
		CONFLUENCE_CLOUD_DOMAIN,
		CONFLUENCE_TOKEN,
		CONFLUENCE_USER,
		CONFLUENCE_SPACE,
		CONFLUENCE_PAGE_ID)
}
