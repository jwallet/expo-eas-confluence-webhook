package main

type PlatformType string

const (
	ios     PlatformType = "ios"
	android PlatformType = "android"
)

type StatusType string

const (
	finished  StatusType = "finished"
	errored   StatusType = "errored"
	cancelled StatusType = "cancelled"
)

type ExpoBuild struct {
	Id                  string       `json:"id"`
	Platform            PlatformType `json:"platform"`
	Status              StatusType   `json:"status"`
	BuildDetailsPageUrl string       `json:"buildDetailsPageUrl"`
	Artifacts           Artifact     `json:"artifacts"`
	Metadata            Metadata     `json:"metadata"`
	CompletedAt         string       `json:"completedAt"`
	ExpirationDate      string       `json:"expirationDate"`
}

type Artifact struct {
	Build string `json:"build"`
}

type Metadata struct {
	AppVersion      string      `json:"appVersion"`
	AppBuildVersion string      `json:"appBuildVersion"`
	BuildProfile    Environment `json:"buildProfile"`
	SdkVersion      string      `json:"sdkVersion"`
}

type Environment string

const (
	review      Environment = "review"
	continuous  Environment = "continuous"
	integration Environment = "integration"
	staging     Environment = "staging"
	production  Environment = "production"
)
