package main

import (
	"fmt"
	"regexp"
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

var buildMapper = map[Environment]Build{
	continuous:  {},
	integration: {},
	staging:     {},
	review:      {},
}

func webhookHandler(context ExpoBuild) error {
	templateTableKey := string(context.Metadata.BuildProfile) + "-" + string(context.Platform)
	build := Build{
		Platform:    string(context.Metadata.BuildProfile),
		Key:         templateTableKey,
		Id:          context.Id,
		Version:     context.Metadata.AppVersion,
		Sdk:         context.Metadata.SdkVersion,
		CompletedAt: context.CompletedAt,
		ExpiresAt:   context.ExpirationDate,
	}

	buildTemplate := generateBuildTemplate(build)

	previousPage, err := getConfluencePage(CONFLUENCE_PAGE_ID)
	if err != nil {
		return err
	}
	fmt.Printf("Got '%v' Page from Confluence\n", previousPage.Title)

	messageVersion := fmt.Sprintf("EAS build %v finished", context.Id)
	storageValue := updateStorageValueWithNewBuildTemplate(previousPage.Body.Storage.Value, buildTemplate, templateTableKey)
	if storageValue == "" {
		return fmt.Errorf("Did not find any valid <table> tag.")
	}
	nextPage := generateConfluenceUpdatePagePayload(previousPage, messageVersion, storageValue)

	return putConfluencePage(CONFLUENCE_PAGE_ID, nextPage)
}

func generateBuildTemplate(build Build) string {
	buildURL := "https://expo.dev/accounts/guay/projects/guay/builds/" + build.Id
	return getBuildTemplate(build.Key, build.Platform, build.Version, build.Sdk, buildURL, build.CompletedAt, build.ExpiresAt)
}

func updateStorageValueWithNewBuildTemplate(storageValue string, buildTemplate string, tableKey string) string {
	selector := regexp.MustCompile(fmt.Sprintf(`<table data-layout="default" ac:local-id="%v">.*?</table>`, tableKey))
	parts := selector.Split(storageValue, 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[0] + buildTemplate + parts[1]
}
