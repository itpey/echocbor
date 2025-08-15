package echocbor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

// Sample struct for testing
type Sample struct {
	Name  string `cbor:"name"`
	Value int    `cbor:"value"`
}

// Test Cbor response
func TestCborResponse(t *testing.T) {
	e := echo.New()
	e.Use(ContextWrapper)

	e.GET("/cbor", func(c echo.Context) error {
		return c.(Context).Cbor(http.StatusOK, Sample{Name: "test", Value: 42})
	})

	req := httptest.NewRequest(http.MethodGet, "/cbor", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/cbor", rec.Header().Get(echo.HeaderContentType))

	var result Sample
	err := cbor.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 42, result.Value)
}

// Test Cbor binding
func TestCborBinding(t *testing.T) {
	e := echo.New()
	e.Binder = &Binder{}

	e.POST("/cbor", func(c echo.Context) error {
		var sample Sample
		if err := c.Bind(&sample); err != nil {
			return err
		}
		return c.String(http.StatusOK, sample.Name)
	})

	sample := Sample{Name: "test", Value: 42}
	cborData, err := cbor.Marshal(sample)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/cbor", bytes.NewReader(cborData))
	req.Header.Set(echo.HeaderContentType, "application/cbor")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
}

// Test unsupported media type
func TestUnsupportedMediaType(t *testing.T) {
	e := echo.New()
	e.Binder = &Binder{}

	e.POST("/cbor", func(c echo.Context) error {
		var sample Sample
		if err := c.Bind(&sample); err != nil {
			return err
		}
		return c.String(http.StatusOK, sample.Name)
	})

	req := httptest.NewRequest(http.MethodPost, "/cbor", bytes.NewReader([]byte(`{"name": "test", "value": 42}`)))
	req.Header.Set(echo.HeaderContentType, "application/json") // Unsupported media type
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnsupportedMediaType, echo.ErrUnsupportedMediaType.Code)
}

// Test invalid CBOR data
func TestInvalidCborData(t *testing.T) {
	e := echo.New()
	e.Binder = &Binder{}

	e.POST("/cbor", func(c echo.Context) error {
		var sample Sample
		if err := c.Bind(&sample); err != nil {
			return err
		}
		return c.String(http.StatusOK, sample.Name)
	})

	req := httptest.NewRequest(http.MethodPost, "/cbor", bytes.NewReader([]byte{0xFF})) // Invalid CBOR
	req.Header.Set(echo.HeaderContentType, "application/cbor")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnsupportedMediaType, echo.ErrUnsupportedMediaType.Code)
}

// Test Cbor with empty body
func TestCborEmptyBody(t *testing.T) {
	e := echo.New()
	e.Binder = &Binder{}

	e.POST("/cbor", func(c echo.Context) error {
		var sample Sample
		if err := c.Bind(&sample); err != nil {
			return err
		}
		return c.String(http.StatusOK, sample.Name)
	})

	req := httptest.NewRequest(http.MethodPost, "/cbor", bytes.NewReader([]byte{})) // Empty body
	req.Header.Set(echo.HeaderContentType, "application/cbor")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnsupportedMediaType, echo.ErrUnsupportedMediaType.Code)
}

// Test Cbor with body limit
func TestCborBodyLimit(t *testing.T) {
	e := echo.New()
	e.Binder = &Binder{}

	e.Use(middleware.BodyLimit("1B"))

	e.POST("/cbor", func(c echo.Context) error {
		var sample Sample
		if err := c.Bind(&sample); err != nil {
			return err
		}
		return c.String(http.StatusOK, sample.Name)
	})

	largeBody := bytes.Repeat([]byte{0x00}, 2)
	req := httptest.NewRequest(http.MethodPost, "/cbor", bytes.NewReader(largeBody))
	req.Header.Set(echo.HeaderContentType, "application/cbor")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
}

// Benchmark for Cbor encode
func BenchmarkCborEncode(b *testing.B) {
	e := echo.New()
	data := Sample{Name: "Test", Value: 42}

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c := &context{e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)}

		if err := c.Cbor(http.StatusOK, data); err != nil {
			b.Fatalf("Failed to encode CBOR: %v", err)
		}
	}
}

// Benchmark for decoding a Cbor request body
func BenchmarkCborDecode(b *testing.B) {
	e := echo.New()
	data := Sample{Name: "Test", Value: 42}
	encodedData, _ := cbor.Marshal(data)

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(encodedData))
		req.Header.Set(echo.HeaderContentType, mimeType)

		c := e.NewContext(req, rec)

		var decodedData Sample
		if err := new(Binder).Bind(&decodedData, c); err != nil {
			b.Fatalf("Failed to decode CBOR: %v", err)
		}
	}
}
