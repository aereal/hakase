package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

type arguments struct {
	repoPath       string
	files          filesList
	maxCommits     int
	pullRequestURL string
}

type filesList []string

func (f *filesList) String() string {
	buf := ""
	for _, v := range *f {
		buf += v
	}
	return buf
}

func (f *filesList) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	var collector candidatesCollector
	if len(args.pullRequestURL) != 0 {
		pr, err := newPullRequest(args.pullRequestURL)
		if err != nil {
			log.Fatalf("Invalid Pull Request URL: %s", err)
		}
		ghToken, err := getConfig("hakase.token")
		if err != nil || ghToken == "" {
			log.Fatalf("Access token must be stored config with key: hakase.token")
		}
		apiBase, err := getGitHubAPIBase(args.pullRequestURL)
		if err != nil {
			log.Fatalf("Cannot determine GitHub API base: %s", err)
		}
		gh := newGitHubClient(apiBase, ghToken)
		collector = newPullRequestCollector(gh, pr)
	} else {
		collector = newArgsCollector(args.files)
	}
	cpus := runtime.NumCPU()
	scanner := newConcurrentRepoScanner(cpus*2, args.maxCommits)
	res := scanner.scan(args.repoPath, collector)
	jsonRet, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Fprintf(os.Stdout, string(jsonRet))
}

func parseArgs() (*arguments, error) {
	args := &arguments{}
	flag.StringVar(&args.repoPath, "repo", "", "repository path to scan")
	flag.IntVar(&args.maxCommits, "max_commits", 100, "max count of commits to scan")
	flag.Var(&args.files, "file", "file path to scan")
	flag.StringVar(&args.pullRequestURL, "pull_request", "", "Pull Request URL to scan")
	flag.Parse()
	if len(args.repoPath) == 0 {
		return nil, fmt.Errorf("repo-path required")
	}
	return args, nil
}
