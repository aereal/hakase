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
	repoPath   string
	files      filesList
	maxCommits int
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

	cpus := runtime.NumCPU()
	scanner := newConcurrentRepoScanner(cpus*2, args.maxCommits)
	res := scanner.scan(args.repoPath, args.files)
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
	flag.Parse()
	if len(args.repoPath) == 0 {
		return nil, fmt.Errorf("repo-path required")
	}
	return args, nil
}
