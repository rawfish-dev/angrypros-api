package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	knownEnvironments = []string{"test", "dev", "prod"}
)

type AppConfig struct {
	GoogleConfig             GoogleConfig             `json:"google"`
	PostgresConfig           PostgresConfig           `json:"postgres"`
	DigitalOceanSpacesConfig DigitalOceanSpacesConfig `json:"dospaces"`
	EntryConfig              EntryConfig              `json:"entry"`
	FeedConfig               FeedConfig               `json:"feed"`
	UserConfig               UserConfig               `json:"user"`
}

type GoogleConfig struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"project_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

type PostgresConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}

type DigitalOceanSpacesConfig struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	Endpoint string `json:"endpoint"`
	Region   string `json:"region"`
}

type EntryConfig struct {
	EntryTextContentMaximumLength int `json:"entryTextContentMaximumLength"`
	InitialLoadCommentCount       int `json:"initialLoadCommentCount"`
	SubsequentLoadCommentCount    int `json:"subsequentLoadCommentCount"`
}

type FeedConfig struct {
	DefaultPageSize int `json:"defaultPageSize"`
}

type UserConfig struct {
	PasswordMinimumLength int    `json:"passwordMinimumLength"`
	UsernameMinimumLength int    `json:"usernameMinimumLength"`
	UsernameMaximumLength int    `json:"usernameMaximumLength"`
	UsernameRegex         string `json:"usernameRegex"`
}

func NewAppConfig(env, directoryPrefix string) AppConfig {
	validEnvironment := false
	for _, knownEnvironment := range knownEnvironments {
		if env == knownEnvironment {
			validEnvironment = true
			break
		}
	}
	if !validEnvironment {
		panic(fmt.Sprintf("'%s' is not a known environment!", env))
	}

	configFilePath := fmt.Sprintf("%s/config/app-%s.json", directoryPrefix, env)
	fullConfigFilePath, err := filepath.Abs(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("Unable to load %s file due to %s", configFilePath, err))
	}

	configFile, err := os.Open(fullConfigFilePath)
	if err != nil {
		panic(fmt.Sprintf("unable to open config file at %s due to %s", fullConfigFilePath, err))
	}
	defer configFile.Close()

	configData, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic(fmt.Sprintf("unable to read config file at %s due to %s", fullConfigFilePath, err))
	}

	var appConfig AppConfig
	err = json.Unmarshal(configData, &appConfig)
	if err != nil {
		panic(fmt.Sprintf("unable to unmarshal config file at %s due to %s", fullConfigFilePath, err))
	}

	return appConfig
}
