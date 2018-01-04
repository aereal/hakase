package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
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

func scanFile(repoPath string, file string, maxCommits int) (ownershipSurvey, error) {
	outBuf := new(bytes.Buffer)
	cmd := exec.Command("git", "-C", repoPath, "log", "-C", "-M", "--no-merges", "-n", fmt.Sprintf("%d", maxCommits), "--format=%aN", file)
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
