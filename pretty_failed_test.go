package godogx_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bool64/godogx"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func TestRegisterPrettyFailedFormatter(t *testing.T) {
	godogx.RegisterPrettyFailedFormatter()

	out := bytes.NewBuffer(nil)

	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			s.Step("I pass", func() {})
			s.Step("I fail", func() error { return errors.New("failed") })
		},
		Options: &godog.Options{
			Format:   "pretty-failed",
			Output:   out,
			NoColors: true,
			Paths:    []string{"_testdata"},
		},
	}

	st := suite.Run()
	assert.Equal(t, 1, st) // Failed.

	assert.Equal(t, `Feature: that fails

  Scenario: pass then fail # _testdata/Failed.feature:3
    When I pass            # pretty_failed_test.go:20 -> github.com/bool64/godogx_test.TestRegisterPrettyFailedFormatter.func1.1
    Then I fail            # pretty_failed_test.go:21 -> github.com/bool64/godogx_test.TestRegisterPrettyFailedFormatter.func1.2
    failed

Feature: that passes

--- Failed steps:

  Scenario: pass then fail # _testdata/Failed.feature:3
    Then I fail # _testdata/Failed.feature:7
      Error: failed


2 scenarios (1 passed, 1 failed)
5 steps (3 passed, 1 failed, 1 skipped)
`, out.String()[0:584])
}
