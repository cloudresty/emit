# golog

`golog` is a simple Go logging package that provides two primary logging formats: plain text and JSON. It offers color-coded severity levels for plain text logs and allows for optional application name and version information to be included in both formats.

&nbsp;

## Features

* **Plain Text Logging:**
  * Color-coded severity levels (info, warning, error, debug) for improved readability.
  * Aligned output for consistent formatting.
  * Includes timestamp, severity, optional application name and version, and the log message.
* **JSON Logging:**
  * Structured log output in JSON format.
  * Includes timestamp, severity, message, and optional application name and version.
* **Optional Application Information:**
  * Application name and version can be optionally included in both plain text and JSON logs.
* **Simple API:** Easy-to-use functions for both logging formats.
* **No external dependencies:** It uses the Go's native libraries.

&nbsp;

## Installation

```bash
go get github.com/cloudresty/golog
```

&nbsp;

## Usage

&nbsp;

### Plain Text Logging

The `Plain` function generates color-coded plain text logs.

```go

package main

import (
	"github.com/cloudresty/golog"
)

func main() {

    // Basic usage with only severity and message.
    golog.Plain("info", "This is an informational message.")

    // With application name.
    golog.Plain("warning", "Something might be wrong", "myapp")

    // With application name and version.
    golog.Plain("error", "Something went wrong!", "myapp", "1.2.3")

    //Debug log
    golog.Plain("debug", "This is a debug message", "my-app", "1.0.0")

    //Incorrect log level will default to info:
    golog.Plain("any", "This is an unknown message", "my-app", "1.0.0")

}

```

**Output example:**

```text

2023-11-26 19:00:00 | info    | This is an informational message.
2023-11-26 19:00:00 | warning | myapp: Something might be wrong
2023-11-26 19:00:00 | error   | myapp 1.2.3: Something went wrong!
2023-11-26 19:00:00 | debug   | my-app 1.0.0: This is a debug message
2023-11-26 19:00:00 | info    | my-app 1.0.0: This is an unknown message

```

**Note:** The color codes might not be displayed properly in all terminals. If you are using Windows, color code will not be supported.

&nbsp;

### JSON Logging

The `JSON` function generates JSON-formatted logs.

```go

package main

import (
	"github.com/cloudresty/golog"
)

func main() {

    // Basic usage with only severity and message.
    golog.JSON("info", "This is an informational JSON message.")

    // With application name.
    golog.JSON("warning", "Something might be wrong in JSON", "myapp")

    // With application name and version.
    golog.JSON("error", "Something went wrong in JSON!", "myapp", "1.2.3")

    //Debug log
    golog.JSON("debug", "This is a JSON debug message", "my-app", "1.0.0")

}

```

**Output example:**

```json

{"message":"This is an informational JSON message.","severity":"info","timestamp":"2023-11-26T19:00:00Z"}
{"application_name":"myapp","message":"Something might be wrong in JSON","severity":"warning","timestamp":"2023-11-26T19:00:00Z"}
{"application_name":"myapp","application_version":"1.2.3","message":"Something went wrong in JSON!","severity":"error","timestamp":"2023-11-26T19:00:00Z"}
{"application_name":"my-app","application_version":"1.0.0","message":"This is a JSON debug message","severity":"debug","timestamp":"2023-11-26T19:00:00Z"}

```

&nbsp;

### Functions

#### `Plain(severity, message string, optionalParams ...string)`

* **severity:** The severity level of the log message. Valid values are "info", "warning", "error", and "debug" (case-insensitive). If an invalid value is provided, it defaults to "info".
* **message:** The log message string.
* **optionalParams:** Optional parameters for application information.
  * If one parameter is passed, it's considered the application name.
  * If two parameters are passed, the first is the application name, and the second is the application version.

&nbsp;

#### `JSON(severity, message string, optionalParams ...string)`

* **severity:** The severity level of the log message. Valid values are "info", "warning", "error", and "debug" (case-insensitive). If an invalid value is provided, it defaults to "info".
* **message:**  The log message string.
* **optionalParams:** Optional parameters for application information.
  * If one parameter is passed, it's considered the application name.
  * If two parameters are passed, the first is the application name, and the second is the application version.

&nbsp;

### Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

&nbsp;

---
Made with ❤️ by [Cloudresty](https://cloudresty.com)
