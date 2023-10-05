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

type Environment string

// List of available environments
const (
	Review      Environment = "reviewapp"
	Continuous  Environment = "continuous"
	Integration Environment = "integration"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

type Metadata struct {
	AppVersion      string      `json:"appVersion"`
	AppBuildVersion string      `json:"appBuildVersion"`
	BuildProfile    Environment `json:"buildProfile"`
	SdkVersion      string      `json:"sdkVersion"`
}
