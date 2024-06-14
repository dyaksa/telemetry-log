# Log Telemetry

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

	l.Log.WithTrace(errors.New("error message")).Error("error message")
}
```

