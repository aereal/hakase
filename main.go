package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
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
	procs := cpus * 2
	res := run(procs, args)
	jsonRet, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Fprintf(os.Stdout, string(jsonRet))
}

func run(maxConcurrency int, args *arguments) surveyResult {
	semaphore := make(chan int, maxConcurrency)
	ch := make(chan *fileSurveyResult, len(args.files))
	var wg sync.WaitGroup
	go func(ch chan *fileSurveyResult) {
		for _, argFile := range args.files {
			wg.Add(1)
			go func(ch chan *fileSurveyResult, file string) {
				defer wg.Done()
				semaphore <- 1
				ret, err := scanFile(args.repoPath, file, args.maxCommits)
				fileRet := &fileSurveyResult{file: file, result: ret}
				if err != nil {
					log.Printf("error: %s", err)
				}
				ch <- fileRet
				<-semaphore
			}(ch, argFile)
		}
		wg.Wait()
		close(ch)
	}(ch)

	res := surveyResult{
		Files: map[string]ownershipSurvey{},
	}
	for fileRet := range ch {
		res.Files[fileRet.file] = fileRet.result
	}
	return res
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
