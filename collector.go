package main

import (
	"log"
)

type candidateFile string

type candidatesCollector interface {
	collectsCandidates() []candidateFile
}

type argsCollector struct {
	files filesList
}

func newArgsCollector(files filesList) argsCollector {
	return argsCollector{files: files}
}

func (c argsCollector) collectsCandidates() []candidateFile {
	cs := make([]candidateFile, len(c.files))
	for i, f := range c.files {
		cs[i] = candidateFile(f)
	}
	return cs
}

type pullRequestCollector struct {
	pr           *pullRequest
	gitHubClient *gitHubClient
}

func newPullRequestCollector(gitHubClient *gitHubClient, pr *pullRequest) pullRequestCollector {
	return pullRequestCollector{
		gitHubClient: gitHubClient,
		pr:           pr,
	}
}

func (c pullRequestCollector) collectsCandidates() []candidateFile {
	cs := []candidateFile{}
	prFiles, err := c.gitHubClient.getPullRequestFiles(c.pr)
	if err != nil {
		log.Fatalf("failed: %#v", err)
	}
	for _, f := range prFiles {
		cs = append(cs, candidateFile(f.Filename))
	}
	return cs
}
