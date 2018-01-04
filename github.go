package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type gitHubClient struct {
	token   string
	apiBase string
}

type pullRequestFile struct {
	Filename string `json:"filename"`
}

func newGitHubClient(apiBase string, token string) *gitHubClient {
	return &gitHubClient{apiBase: apiBase, token: token}
}

func (c *gitHubClient) getPullRequestFiles(pr *pullRequest) ([]pullRequestFile, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/%s/pulls/%s/files", c.apiBase, pr.owner, pr.repo, pr.number)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+c.token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get pull request (%s/%s#%s) files: %s", pr.owner, pr.repo, pr.number, err)
	}

	decoder := json.NewDecoder(resp.Body)
	_, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %s", err)
	}

	var prFiles []pullRequestFile
	for decoder.More() {
		var prFile pullRequestFile
		err := decoder.Decode(&prFile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse body: %s", err)
		}
		prFiles = append(prFiles, prFile)
	}

	_, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to parse body: %s", err)
	}

	return prFiles, nil
}
