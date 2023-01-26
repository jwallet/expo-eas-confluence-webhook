package main

import (
	"fmt"
	"regexp"
)

var ExpoBuildProfilePageId = map[string]PageId{
	"continuous:android":  continuous_android,
	"continuous:ios":      continuous_ios,
	"integration:android": integration_android,
	"integration:ios":     integration_ios,
	"staging:android":     staging_android,
	"staging:ios":         staging_ios,
	"review:android":      review_android,
	"review:ios":          review_ios,
}

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
	profileKey := string(context.Metadata.BuildProfile) + ":" + string(context.Platform)
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

	pageId := ExpoBuildProfilePageId[profileKey]
	buildTemplate := generateBuildTemplate(build)

	previousPage, err := getConfluencePage(pageId)
	if err != nil {
		return err
	}
	fmt.Printf("Got '%v' Page from Confluence\n", previousPage.Title)

	storageValue := updateStorageValueWithNewBuildTemplate(previousPage.Body.Storage.Value, buildTemplate, templateTableKey)

	nextPage := ConfluencePage{
		Version: PageVersion{
			Number:  previousPage.Version.Number + 1,
			Message: fmt.Sprintf("EAS build %v finished", context.Id),
		},
		PageType: "page",
		Status:   "current",
		Title:    previousPage.Title,
		Space: PageSpace{
			Key: "BLOG",
		},
		Body: PageBody{
			Storage: PageStorage{
				Value:          storageValue,
				Representation: "storage",
			},
		},
	}

	return putConfluencePage(pageId, &nextPage)
}

func generateBuildTemplate(build Build) string {
	buildURL := "https://expo.dev/accounts/guay/projects/guay/builds/" + build.Id
	return getBuildTemplate(build.Key, build.Platform, build.Version, build.Sdk, buildURL, build.CompletedAt, build.ExpiresAt)
}

func updateStorageValueWithNewBuildTemplate(storageValue string, buildTemplate string, tableKey string) string {
	selector := regexp.MustCompile(fmt.Sprintf(`<table data-layout="default" ac:local-id="%v">.*</table>`, tableKey))
	parts := selector.Split(storageValue, 2)
	return parts[0] + buildTemplate + parts[1]
}
