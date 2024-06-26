# Log Telemetry

[![Build Status](https://img.shields.io/github/actions/workflow/status/caarlos0/env/build.yml?branch=main&style=for-the-badge)](https://github.com/dyaksa/telemetry-log/actions)
[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://pkg.go.dev/github.com/dyaksa/telemetry-log)
[![codecov](https://codecov.io/github/dyaksa/telemetry-log/graph/badge.svg?token=3M9N4DDKCZ)](https://codecov.io/github/dyaksa/telemetry-log)

A simple dependency for trace log and store trace to `mongo db`

## Installation

<p>
  <a href="https://encore.dev">
    <img src="https://i.ibb.co.com/w7xD7fL/Desain-tanpa-judul.png" width="120px" alt="encore icon"></img>
  </a>
  <br/>
  <br/>
  <br/>
</p>

Get the module with:

```sh
go get github.com/dyaksa/telemetry-log/telemetry@latest
```

The usage looks like this:

```go
package main

import (
	"errors"

	"github.com/dyaksa/telemetry-log/telemetry"
)

func main() {
	l, err := telemetry.New(telemetry.WithJSONFormatter())
	if err != nil {
		panic(err)
	}

	l.Log.WithTrace(errors.New("error etc...")).Info("info message")
}
```
You can run it like this:

With `telemetry.New(telemetry.WithJSONFormatter())`
```text
{"file":"main.go","func":"main.main","level":"info","line":15,"msg":"info message","time":"2024-06-14T22:36:10+07:00","trace":[{"name":"runtime.main","file":"proc.go","line":"271"},{"name":"runtime.goexit","file":"asm_arm64.s","line":"1222"}]}
```

With the default `telemetry.New()`.

```sh
INFO[0000] info message file=main.go func=main.main line=15 trace="[{runtime.main proc.go 271} {runtime.goexit asm_arm64.s 1222}]"
```

#### Logging Method Name

If you wish to add the calling method `WithTracer(errors.New("error etc..."))` you can use the `WithTracer(errors.New("error etc..."))` method.

```go
l, err := telemetry.New(telemetry.WithJSONFormatter())

if err != nil {
    panic(err)
}

l.Log.WithTrace(errors.New("error etc")).Info("info message")
```
This adds the caller as 'method' like so:
```json
{"file":"main.go","func":"main.main","level":"info","line":15,"msg":"info message","time":"2024-06-14T22:51:12+07:00","trace":[{"name":"runtime.main","file":"proc.go","line":"271"},{"name":"runtime.goexit","file":"asm_arm64.s","line":"1222"}]}
```

With the default without `WithTracer(error)`.
```go
l.Log.Info("info message")
```

This adds the caller as 'method' like so
```json
{"file":"main.go","func":"main.main","level":"info","line":13,"msg":"info message","time":"2024-06-14T22:54:41+07:00"}
```

#### Level logging

Logrus has seven logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic.

```go
log.Trace("Something very low level.")
log.Debug("Useful debugging information.")
log.Info("Something noteworthy happened!")
log.Warn("You should probably take a look at this.")
log.Error("Something failed but I'm not quitting.")
// Calls os.Exit(1) after logging
log.Fatal("Bye.")
// Calls panic() after logging
log.Panic("I'm bailing.")
```

#### Environments






