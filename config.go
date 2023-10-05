package main

import (
	"os"
)

var EXPO_HMAC_SECRET = os.Getenv("EXPO_HMAC_SECRET")

var PORT int = getEnvInt("PORT")

var CONFLUENCE_CLOUD_DOMAIN = os.Getenv("CONFLUENCE_CLOUD_DOMAIN")
var CONFLUENCE_TOKEN = os.Getenv("CONFLUENCE_TOKEN")
var CONFLUENCE_USER = os.Getenv("CONFLUENCE_USER")
var CONFLUENCE_PAGE_ID int = getEnvInt("CONFLUENCE_PAGE_ID")
var CONFLUENCE_SPACE = os.Getenv("CONFLUENCE_SPACE")

// Titles related to environments -- Used for Confluence
var environments = map[Environment]string{
	Review:      "Review App",
	Continuous:  "Continuous",
	Integration: "Integration",
	Staging:     "Staging",
	Production:  "Production",
}
