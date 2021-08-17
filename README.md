# godogx

[![Build Status](https://github.com/bool64/godogx/workflows/test-unit/badge.svg)](https://github.com/bool64/godogx/actions?query=branch%3Amaster+workflow%3Atest-unit)
[![Coverage Status](https://codecov.io/gh/bool64/godogx/branch/master/graph/badge.svg)](https://codecov.io/gh/bool64/godogx)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/bool64/godogx)
<!--[![Time Tracker](https://wakatime.com/badge/github/bool64/godogx.svg)](https://wakatime.com/badge/github/bool64/godogx)
![Code lines](https://sloc.xyz/github/bool64/godogx/?category=code)
![Comments](https://sloc.xyz/github/bool64/godogx/?category=comments)--->

A library of [`godog`](https://github.com/cucumber/godog) tools and extensions.

## Pretty Failed Formatter

When running a bit test suite, where most of the scenarios pass, the output becomes less helpful if you try to check the
failing scenarios.

`pretty-failed` formatter extends `pretty` formatter, but discards output of successful scenarios and also does not show
skipped steps after the failure was encountered.

You can enable it by calling `godogx.RegisterPrettyFailedFormatter()`.

## Allure Formatter

[Allure](https://github.com/allure-framework/allure2) is convenient UI to expose test results.

You can enable it by calling `allure.RegisterFormatter()`.

Additional configuration can be added with env vars before test run.

`ALLURE_ENV_*` are added to allure environment report.

`ALLURE_EXECUTOR_*` configure `Executor` info.

`ALLURE_RESULTS_PATH` can change default `./allure-results` destination.

Example:
```bash
export ALLURE_ENV_TICKET=JIRA-1234
export ALLURE_ENV_APP=todo-list
export ALLURE_EXECUTOR_NAME=IntegrationTest
export ALLURE_EXECUTOR_TYPE=github
```

Then you can run test with 
```bash
# Optionally clean up current result (if you have it).
rm -rf ./allure-results/*
# Optionally copy history from previous report.
cp -r ./allure-report/history ./allure-results/history
# Run suite with godog CLI tool or with go test.
godog -f allure
# Generate report with allure CLI tool.
allure generate --clean
```