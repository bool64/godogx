package allure_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bool64/godogx/allure"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	allure.RegisterAllureFormatter("./allure-results")

	out := bytes.NewBuffer(nil)

	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			s.Step("I pass", func() {})
			s.Step("I fail", func() error { return errors.New("failed") })
		},
		Options: &godog.Options{
			Format:   "allure",
			Output:   out,
			NoColors: true,
			Paths:    []string{"../_testdata"},
		},
	}

	st := suite.Run()
	assert.Equal(t, 1, st) // Failed.
}

// rm -rf ./allure-results/* && cp -r ./allure-report/history ./allure-results/history  && make integration-test; allure generate --clean
