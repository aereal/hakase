package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

type ownershipScore struct {
	Count int     `json:"count"`
	Score float32 `json:"score"`
}

type ownershipSurvey map[string]ownershipScore

type fileSurveyResult struct {
	file   string
	result ownershipSurvey
}

type surveyResult struct {
	Files map[string]ownershipSurvey `json:"files"`
}

type repoScanner interface {
	scan(repo string, files filesList) surveyResult
}

type concurrentRepoScanner struct {
	maxConcurrency int
	maxCommits     int
}

func newConcurrentRepoScanner(maxConcurrency int, maxCommits int) concurrentRepoScanner {
	scanner := concurrentRepoScanner{
		maxConcurrency: maxConcurrency,
		maxCommits:     maxCommits,
	}
	return scanner
}

func (s concurrentRepoScanner) scan(repo string, files filesList) surveyResult {
	semaphore := make(chan int, s.maxConcurrency)
	ch := make(chan *fileSurveyResult, len(files))
	var wg sync.WaitGroup
	go func(ch chan *fileSurveyResult) {
		for _, argFile := range files {
			wg.Add(1)
			go func(ch chan *fileSurveyResult, file string) {
				defer wg.Done()
				semaphore <- 1
				ret, err := s.scanFile(repo, file)
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

func (s concurrentRepoScanner) scanFile(repoPath string, file string) (ownershipSurvey, error) {
	outBuf := new(bytes.Buffer)
	cmd := exec.Command("git", "-C", repoPath, "log", "-C", "-M", "--no-merges", "-n", fmt.Sprintf("%d", s.maxCommits), "--format=%aN", file)
	cmd.Stdout = outBuf
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(outBuf)
	totalChangesCount := 0
	changesByAuthor := map[string]int{}
	for scanner.Scan() {
		author := scanner.Text()
		changesByAuthor[author]++
		totalChangesCount++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	res := ownershipSurvey{}
	for author, count := range changesByAuthor {
		res[author] = ownershipScore{
			Count: count,
			Score: float32(count) / float32(totalChangesCount),
		}
	}
	return res, nil
}
