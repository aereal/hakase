package main

import (
	"fmt"
	"net/url"
	"strings"
)

type pullRequest struct {
	owner  string
	repo   string
	number string
}

func newPullRequest(rawurl string) (*pullRequest, error) {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	paths := strings.Split(parsed.Path, "/")
	if len(paths) != 5 {
		return nil, fmt.Errorf("Invalid Pull Request URL: %s", parsed)
	}
	pr := &pullRequest{
		owner:  paths[1],
		repo:   paths[2],
		number: paths[4],
	}
	return pr, nil
}

func getGitHubAPIBase(rawurl string) (string, error) {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	if parsed.Hostname() == "github.com" {
		return "https://api.github.com", nil
	} else {
		return fmt.Sprintf("%s://%s/api/v3", parsed.Scheme, parsed.Hostname()), nil
	}
}
