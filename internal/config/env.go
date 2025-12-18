// Package config handles environment variable loading and secret management.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Environment variable keys
const (
	KeyGitHubToken     = "GITHUB_TOKEN"
	KeyVirusTotalToken = "VIRUS_TOTAL_TOKEN"
	KeyDockerHubUser   = "DOCKERHUB_USER"
	KeyDockerHubToken  = "DOCKERHUB_TOKEN"
)

// Secrets holds all loaded API credentials.
type Secrets struct {
	GitHubToken     string
	VirusTotalToken string
	DockerHubUser   string
	DockerHubToken  string
}

// LoadSecrets loads environment variables from .env file and returns Secrets.
// Returns error if .env file cannot be loaded.
func LoadSecrets() (*Secrets, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	return &Secrets{
		GitHubToken:     os.Getenv(KeyGitHubToken),
		VirusTotalToken: os.Getenv(KeyVirusTotalToken),
		DockerHubUser:   os.Getenv(KeyDockerHubUser),
		DockerHubToken:  os.Getenv(KeyDockerHubToken),
	}, nil
}

// MustLoadSecrets loads secrets and panics on failure.
func MustLoadSecrets() *Secrets {
	s, err := LoadSecrets()
	if err != nil {
		panic(err)
	}
	return s
}

// HasGitHub returns true if GitHub token is configured.
func (s *Secrets) HasGitHub() bool {
	return s.GitHubToken != ""
}

// HasVirusTotal returns true if VirusTotal token is configured.
func (s *Secrets) HasVirusTotal() bool {
	return s.VirusTotalToken != ""
}

// HasDockerHub returns true if Docker Hub credentials are configured.
func (s *Secrets) HasDockerHub() bool {
	return s.DockerHubUser != "" && s.DockerHubToken != ""
}
