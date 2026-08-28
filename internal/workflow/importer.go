package workflow

import (
	"bookstore/recommendation/internal/model"
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ImportResult struct {
	Created []string `json:"created"`
	Failed  []string `json:"failed"`
}

func (s *Service) ImportCSV(reader io.Reader, actor string, at int64) ImportResult {
	result := ImportResult{Created: []string{}, Failed: []string{}}
	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || lineNumber == 1 && strings.HasPrefix(strings.ToLower(line), "id,") {
			continue
		}
		input, err := parseLine(line)
		if err != nil {
			result.Failed = append(result.Failed, fmt.Sprintf("line %d: %v", lineNumber, err))
			continue
		}
		record, err := s.Register(input, at+int64(lineNumber), actor)
		if err != nil {
			result.Failed = append(result.Failed, fmt.Sprintf("line %d: %v", lineNumber, err))
			continue
		}
		result.Created = append(result.Created, record.ID)
	}
	if err := scanner.Err(); err != nil {
		result.Failed = append(result.Failed, err.Error())
	}
	return result
}

func parseLine(line string) (model.RecordInput, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 7 {
		return model.RecordInput{}, fmt.Errorf("expected seven columns")
	}
	score, err := strconv.Atoi(strings.TrimSpace(parts[5]))
	if err != nil {
		return model.RecordInput{}, fmt.Errorf("invalid score")
	}
	return model.RecordInput{ID: parts[0], StoreID: parts[1], Title: parts[2], Author: parts[3], Genre: parts[4], Score: score, Summary: parts[6]}, nil
}
