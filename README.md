# Echo CBOR Middleware

[![Go Version][GoVer-Image]][GoDoc-Url] [![License][License-Image]][License-Url] [![GoDoc][GoDoc-Image]][GoDoc-Url] [![Go Report Card][ReportCard-Image]][ReportCard-Url]

[GoVer-Image]: https://img.shields.io/badge/Go-1.24%2B-blue
[GoDoc-Url]: https://pkg.go.dev/github.com/itpey/echocbor
[GoDoc-Image]: https://pkg.go.dev/badge/github.com/itpey/echocbor.svg
[ReportCard-Url]: https://goreportcard.com/report/github.com/itpey/echocbor
[ReportCard-Image]: https://goreportcard.com/badge/github.com/itpey/echocbor?style=flat
[License-Url]: https://github.com/itpey/echocbor/blob/main/LICENSE
[License-Image]: https://img.shields.io/github/license/itpey/echocbor

Echo CBOR is a middleware for the Echo web framework that enables encoding and decoding of CBOR (Concise Binary Object Representation) data in HTTP requests and responses using the [fxamacker/cbor](https://github.com/fxamacker/cbor) library.

## Features

- Adds a `Cbor` method to the Echo context for sending CBOR-encoded responses.
- Provides a custom `Binder` to decode CBOR request bodies.
- Middleware to wrap the Echo context for CBOR support.
- Lightweight and easy to integrate with existing Echo applications.

## Installation

To install the package, use the following command:

```bash
go get -u github.com/itpey/echocbor
```

## Usage

### Setting Up the Middleware

To use the CBOR middleware, wrap your Echo handlers with the `ContextWrapper` middleware and set the custom `Binder`. Here's an example:

```go
package main

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/itpey/echocbor"
)

func main() {
    e := echo.New()

    // Set the custom CBOR binder
    e.Binder = &echocbor.Binder{}

    // Use the CBOR context wrapper middleware
    e.Use(echocbor.ContextWrapper)

    // Example route
    e.POST("/example", func(c echo.Context) error {
        // Define a struct to bind the request body
        type Payload struct {
            Name string `json:"name" cbor:"name"`
            Age  int    `json:"age" cbor:"age"`
        }

        // Bind the CBOR request body
        var payload Payload
        if err := c.Bind(&payload); err != nil {
            return err
        }

        // Send a CBOR response
        return c.(echocbor.Context).Cbor(http.StatusOK, map[string]interface{}{
            "message": "Received",
            "data":    payload,
        })
    })

    e.Logger.Fatal(e.Start(":8080"))
}
```

### Sending CBOR Responses

The `Cbor` method allows you to send a CBOR-encoded response with a specified HTTP status code. For example:

```go
e.GET("/data", func(c echo.Context) error {
    data := map[string]interface{}{
        "id":   1,
        "name": "John Doe",
    }
    return c.(echocbor.Context).Cbor(http.StatusOK, data)
})
```

## Feedback and Contributions

If you encounter any issues or have suggestions for improvement, please [open an issue](https://github.com/itpey/echocbor/issues) on GitHub.

We welcome contributions! Fork the repository, make your changes, and submit a pull request.

## Support

If you enjoy using Echo CBOR, please consider giving it a star! Your support helps others discover the project and encourages further development.

## License

Echo CBOR is open-source software released under the Apache License, Version 2.0. You can find a copy of the license in the [LICENSE](https://github.com/itpey/echocbor/blob/main/LICENSE) file.
