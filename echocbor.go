package echocbor

import (
	"net/http"

	"github.com/fxamacker/cbor/v2"
	"github.com/labstack/echo/v4"
)

const (
	mimeType = "application/cbor"
)

// Context extends echo.Context to include Cbor response functionality.
type Context interface {
	echo.Context
	Cbor(code int, i interface{}) error
}

// context wraps echo.Context to implement the Cbor method.
type context struct {
	echo.Context
}

// Cbor sends a Cbor response with status code.
func (c context) Cbor(code int, i interface{}) error {
	// encode as Cbor
	b, err := cbor.Marshal(i)
	if err != nil {
		return err
	}

	return c.Blob(code, mimeType, b)
}

// Binder for Cbor.
type Binder struct{}

// Bind binds a request body to given interface.
func (b *Binder) Bind(i interface{}, c echo.Context) error {
	if err := new(echo.DefaultBinder).Bind(i, c); err != echo.ErrUnsupportedMediaType {
		return err
	}

	// check content type
	if c.Request().Header.Get(echo.HeaderContentType) != mimeType {
		return echo.ErrUnsupportedMediaType
	}

	// decode Cbor body
	if err := cbor.NewDecoder(c.Request().Body).Decode(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}
	return nil
}

// ContextWrapper wraps context for sending a Cbor response.
func ContextWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := &context{c} // Create a new context that embeds echo.Context
		return next(ctx)   // Pass the wrapped context to the next handler
	}
}
