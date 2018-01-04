package main

import (
	"reflect"
	"testing"
)

func TestNewPullRequest(t *testing.T) {
	type expectation struct {
		field    string
		expected string
	}

	expectations := map[string][]expectation{
		"https://github.com/aereal/hakase/pulls/1": []expectation{
			expectation{
				field:    "owner",
				expected: "aereal",
			},
			expectation{
				field:    "repo",
				expected: "hakase",
			},
			expectation{
				field:    "number",
				expected: "1",
			},
		},
	}
	for url, expcts := range expectations {
		pr, err := newPullRequest(url)
		if err != nil || pr == nil {
			t.Fatalf("PR can be created by github.com URL but: %v", err)
		}
		for _, exp := range expcts {
			v := reflect.ValueOf(pr)
			got := reflect.Indirect(v).FieldByName(exp.field)
			if got.String() != exp.expected {
				t.Errorf("%s must be %s; but %s", exp.field, exp.expected, got)
			}
		}
	}
}

func TestGetGitHubAPIBase(t *testing.T) {
	type expectation struct {
		input    string
		expected string
	}
	expectations := []expectation{
		expectation{
			input:    "https://github.com",
			expected: "https://api.github.com",
		},
		expectation{
			input:    "https://github.com/aereal/hakase/pulls/1",
			expected: "https://api.github.com",
		},
		expectation{
			input:    "https://ghe.example.com",
			expected: "https://ghe.example.com/api/v3",
		},
	}
	for _, expct := range expectations {
		got, err := getGitHubAPIBase(expct.input)
		if err != nil {
			t.Errorf("%s must be determined but error occurred: %s", expct.input, err)
		}
		if got != expct.expected {
			t.Errorf("baseURL from %s must be %s but %s got", expct.input, expct.expected, got)
		}
	}
}
