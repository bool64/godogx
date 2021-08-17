package allure

import (
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
)

type Suite struct {
	Name      string      `json:"name"`
	Start     TimestampMs `json:"start"`
	Stop      TimestampMs `json:"stop"`
	Version   string      `json:"version"`
	TestCases []TestCase  `json:"testCases,omitempty"`
	Labels    []Label     `json:"labels,omitempty"`
}

// TestCase is the top level report object for a test.
type TestCase struct {
	UUID          string         `json:"uuid,omitempty"`
	Name          string         `json:"name,omitempty"`
	Description   string         `json:"description,omitempty"`
	Status        Status         `json:"status,omitempty"`
	StatusDetails *StatusDetails `json:"statusDetails,omitempty"`
	Stage         string         `json:"stage,omitempty"` // "finished"
	Steps         []Step         `json:"steps,omitempty"`
	Attachments   []Attachment   `json:"attachments,omitempty"`
	Parameters    []Parameter    `json:"parameters,omitempty"`
	Start         TimestampMs    `json:"start,omitempty"`
	Stop          TimestampMs    `json:"stop,omitempty"`
	Children      []string       `json:"children,omitempty"`
	FullName      string         `json:"fullName,omitempty"`
	Labels        []Label        `json:"labels,omitempty"`
	Links         []Link         `json:"links,omitempty"`
}

const (
	Broken  = Status("broken")
	Passed  = Status("passed")
	Failed  = Status("failed")
	Skipped = Status("skipped")
	Unknown = Status("unknown")
)

type TimestampMs int64

type LinkType string

const (
	Issue  LinkType = "issue"
	TMS    LinkType = "tms"
	Custom LinkType = "custom"
)

type Link struct {
	Name string   `json:"name,omitempty"`
	Type LinkType `json:"type,omitempty"`
	URL  string   `json:"url,omitempty"`
}

type Status string

type StatusDetails struct {
	Known   bool   `json:"known,omitempty"`
	Muted   bool   `json:"muted,omitempty"`
	Flaky   bool   `json:"flaky,omitempty"`
	Message string `json:"message,omitempty"`
	Trace   string `json:"trace,omitempty"`
}

type Step struct {
	Name          string         `json:"name,omitempty"`
	Status        Status         `json:"status,omitempty"`
	StatusDetails *StatusDetails `json:"statusDetails,omitempty"`
	Stage         string         `json:"stage"`
	ChildrenSteps []Step         `json:"steps"`
	Attachments   []Attachment   `json:"attachments"`
	Parameters    []Parameter    `json:"parameters"`
	Start         TimestampMs    `json:"start"`
	Stop          TimestampMs    `json:"stop"`
}

type Attachment struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Size   int    `json:"size"`
}

func NewAttachment(name string, mimeType string, resultsPath string, content []byte) (*Attachment, error) {
	a := &Attachment{
		Name: name,
		Type: mimeType,
	}

	id := uuid.New().String()

	ext := ".txt"
	if mimeType == "application/json" {
		ext = ".json"
	} else if mimeType == "image/png" {
		ext = ".png"
	}
	a.Source = fmt.Sprintf("%s-Attachment%s", id, ext)
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s", resultsPath, a.Source), content, 0o600)
	if err != nil {
		return nil, err
	}

	return a, nil
}

type Parameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Label struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
