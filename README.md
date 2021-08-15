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
