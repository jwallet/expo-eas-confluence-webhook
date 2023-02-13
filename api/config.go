package main

import "os"

type Environment string

var EXPO_HMAC_SECRET = os.Getenv("EXPO_HMAC_SECRET")

var PORT int = GetEnvInt("PORT")

var CONFLUENCE_CLOUD_DOMAIN = os.Getenv("CONFLUENCE_CLOUD_DOMAIN")
var CONFLUENCE_TOKEN = os.Getenv("CONFLUENCE_TOKEN")
var CONFLUENCE_USER = os.Getenv("CONFLUENCE_USER")
var CONFLUENCE_PAGE_ID int = GetEnvInt("CONFLUENCE_PAGE_ID")
var CONFLUENCE_SPACE = os.Getenv("CONFLUENCE_SPACE")

// List of available environments
const (
	review      Environment = "reviewapp"
	continuous  Environment = "continuous"
	integration Environment = "integration"
	staging     Environment = "staging"
	production  Environment = "production"
)

// Titles related to environments -- Used for Confluence
var environments = map[Environment]string{
	review:      "Review App",
	continuous:  "Continuous",
	integration: "Integration",
	staging:     "Staging",
	production:  "Production",
}
