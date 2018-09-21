# go-cloudlogzap
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/anexia-it/go-cloudlogzap/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/anexia-it/go-cloudlogzap?status.svg)](https://godoc.org/github.com/anexia-it/go-cloudlogzap)
[![Build Status](https://travis-ci.org/anexia-it/go-cloudlogzap.svg?branch=master)](https://travis-ci.org/anexia-it/go-cloudlogzap)
[![codecov](https://codecov.io/gh/anexia-it/go-cloudlogzap/branch/master/graph/badge.svg)](https://codecov.io/gh/anexia-it/go-cloudlogzap)
[![Go Report Card](https://goreportcard.com/badge/github.com/anexia-it/go-cloudlogzap)](https://goreportcard.com/report/github.com/anexia-it/go-cloudlogzap)

go-cloudlogzap implements a `zapcore` to hook log message from a `zap.Logger` to Anexia's CloudLog.  

## Motivation
The primary motivation for go-cloudlogzap is the need to send log output from `go.uber.org/zap`'s logger to ANEXIA's CloudLog infrastructure.

## Install
`go get -u github.com/anexia-it/go-cloudlogzap`

## Quickstart
Use zap's `zapcore.NewTee` func to tee any log message to multiple zapcores. Use `logger.WithOptions` to create a new logger from the multi-core.
```
opts := []cloudlog.Option{...}
multiCore, err := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
  cloudlogCore, err := NewCloudlogCore(core, indexName, opts)
  if err != nil {
    return core
  }
  return zapcore.NewTee(core, cloudlogCore)
})
logger = logger.WithOptions(multiCore)
```

## Custom CloudLog Options
Important:  
To create a functional cloudlog core, pass the following `cloudlog.Option` slice to the core initialization:
* `cloudlog.OptionCACertificateFile`
* `cloudlog.OptionClientCertificateFile`

## Issue tracker
Issues in go-cloudlogzap are tracked using the corresponding Github [issue tracker](https://github.com/anexia-it/go-cloudlogzap/issues).

## Status
The current release is **v1.0.0**
Changes to go-cloudlogzap are subject to [semantic versioning](http://semver.org/).
The [ChangeLog](https://github.com/anexia-it/go-cloudlogzap/blob/master/CHANGELOG.md) provides information on releases and changes.

## license
go-cloudlogzap is licensed under the terms of the [MIT license](https://github.com/anexia-it/go-cloudlogzap/LICENSE).
