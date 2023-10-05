package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PageVersion struct {
	Message string `json:"message"`
	Number  int16  `json:"number"`
}

type PageStorage struct {
	Representation string `json:"representation"`
	Value          string `json:"value"`
}

type PageBody struct {
	Storage PageStorage `json:"storage"`
}

type PageSpace struct {
	Key string `json:"key"`
}

type ConfluencePage struct {
	Body     PageBody    `json:"body"`
	PageType string      `json:"type"`
	Space    PageSpace   `json:"space"`
	Status   string      `json:"status"`
	Title    string      `json:"title"`
	Version  PageVersion `json:"version"`
}

func GetConfluencePage(pageId int) (*ConfluencePage, error) {
	client := GetClient()
	var currentPage ConfluencePage

	url := fmt.Sprintf("https://%s.atlassian.net/wiki/rest/api/content/%v?expand=version,body.storage", CONFLUENCE_CLOUD_DOMAIN, pageId)
	log.Printf("Sending GET confluence page request %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return currentPage, err
	}
	req.SetBasicAuth(CONFLUENCE_USER, CONFLUENCE_TOKEN)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return currentPage, fmt.Errorf("GET confluence page Failed %v", resp.StatusCode)
	}
	if err != nil {
		return currentPage, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return currentPage, err
	}

	client.CloseIdleConnections()
	json.Unmarshal(body, currentPage)

	return currentPage, nil
}

func PutConfluencePage(pageId int, content *ConfluencePage) error {
	client := GetClient()

	var payload bytes.Buffer

	enc := json.NewEncoder(&payload)
	enc.SetEscapeHTML(false)
	enc.Encode(&content)

	url := fmt.Sprintf("https://%s.atlassian.net/wiki/rest/api/content/%v", CONFLUENCE_CLOUD_DOMAIN, pageId)

	req, err := http.NewRequest("PUT", url, &payload)
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth(CONFLUENCE_USER, CONFLUENCE_TOKEN)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("PUT confluence page Failed (%v)\n", resp.StatusCode)
	}

	defer resp.Body.Close()
	client.CloseIdleConnections()

	return nil
}

func GenerateConfluenceUpdatePagePayload(previousPage *ConfluencePage, messageVersion string, storageValue string) *ConfluencePage {
	return &ConfluencePage{
		Version: PageVersion{
			Number:  previousPage.Version.Number + 1,
			Message: messageVersion,
		},
		PageType: "page",
		Status:   "current",
		Title:    previousPage.Title,
		Space: PageSpace{
			Key: CONFLUENCE_SPACE,
		},
		Body: PageBody{
			Storage: PageStorage{
				Value:          storageValue,
				Representation: "storage",
			},
		},
	}
}
