package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

type Build struct {
	Platform    string `json:"platform,omitempty"`
	Key         string `json:"key,omitempty"`
	Id          string `json:"id,omitempty"`
	Version     string `json:"version,omitempty"`
	Sdk         string `json:"sdk,omitempty"`
	CompletedAt string `json:"completedAt,omitempty"`
	ExpiresAt   string `json:"expiresAt,omitempty"`
}

func Init() error {
	log.Println("Init confluence page")

	previousPage, err := GetConfluencePage(CONFLUENCE_PAGE_ID)
	if err != nil {
		log.Printf("An error occured %v\n", err)
		return err
	}

	template := GetDefaultTemplate()
	minifier := strings.NewReplacer("\n", "", "\t", "")
	minifiedTemplate := minifier.Replace(template)
	page := GenerateConfluenceUpdatePagePayload(previousPage, "Init EAS builds template", minifiedTemplate)

	return PutConfluencePage(CONFLUENCE_PAGE_ID, page)
}

func PushBuild(build Build) error {
	json, _ := json.Marshal(build)
	log.Println(string(json))

	buildTemplate := generateBuildTemplate(build)

	log.Printf("GET confluence page for build #%v profile '%v'\n", build.Id, build.Key)
	previousPage, err := GetConfluencePage(CONFLUENCE_PAGE_ID)
	if err != nil {
		return err
	}
	log.Printf("Got '%v' Page from Confluence\n", previousPage.Title)

	messageVersion := fmt.Sprintf("EAS build %v finished", build.Id)
	storageValue := updateStorageValueWithNewBuildTemplate(previousPage.Body.Storage.Value, buildTemplate, build.Key)
	if storageValue == "" {
		return fmt.Errorf("Did not find any valid <table> tag.")
	}
	nextPage := GenerateConfluenceUpdatePagePayload(previousPage, messageVersion, storageValue)

	log.Printf("PUT confluence page for build #%v profile '%v'\n", build.Id, build.Key)
	return PutConfluencePage(CONFLUENCE_PAGE_ID, nextPage)
}

func Webhook(context ExpoBuild) error {
	templateTableKey := string(context.Metadata.BuildProfile) + "-" + string(context.Platform)
	build := Build{
		Platform:    string(context.Platform),
		Key:         templateTableKey,
		Id:          context.Id,
		Version:     context.Metadata.AppVersion,
		Sdk:         context.Metadata.SdkVersion,
		CompletedAt: context.CompletedAt,
		ExpiresAt:   context.ExpirationDate,
	}

	return retry(1, 1000, func() error {
		return PushBuild(build)
	})
}

func generateBuildTemplate(build Build) string {
	buildURL := "https://expo.dev/accounts/guay/projects/guay/builds/" + build.Id
	return GetBuildTemplate(build.Key, build.Platform, build.Version, build.Sdk, buildURL, build.CompletedAt, build.ExpiresAt)
}

func updateStorageValueWithNewBuildTemplate(storageValue string, buildTemplate string, tableKey string) string {
	selector := regexp.MustCompile(fmt.Sprintf(`<table data-layout="default" ac:local-id="%v">.*?</table>`, tableKey))
	parts := selector.Split(storageValue, 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[0] + buildTemplate + parts[1]
}

func retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("Retrying after error:", err)
			time.Sleep(sleep)
			sleep *= 2
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("After %d attempts, last error: %s", attempts, err)
}
