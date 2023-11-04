// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package gen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// CreateUser defines model for CreateUser.
type CreateUser struct {
	Name string `json:"name"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message string `json:"message"`
}

// User defines model for User.
type User struct {
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
}

// Voucher defines model for Voucher.
type Voucher struct {
	CreatedAt time.Time `json:"created_at"`
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	Value     string    `json:"value"`
}

// SearchVoucherParams defines parameters for SearchVoucher.
type SearchVoucherParams struct {
	// UserId User ID
	UserId int64 `form:"userId" json:"userId"`
}

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = CreateUser

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create voucher
	// (POST /v1/users)
	CreateUser(ctx echo.Context) error
	// Get voucher by id
	// (GET /v1/users/{id})
	GetUserById(ctx echo.Context, id int64) error
	// Search voucher
	// (GET /v1/vouchers)
	SearchVoucher(ctx echo.Context, params SearchVoucherParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateUser converts echo context to params.
func (w *ServerInterfaceWrapper) CreateUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateUser(ctx)
	return err
}

// GetUserById converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUserById(ctx, id)
	return err
}

// SearchVoucher converts echo context to params.
func (w *ServerInterfaceWrapper) SearchVoucher(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params SearchVoucherParams
	// ------------- Required query parameter "userId" -------------

	err = runtime.BindQueryParameter("form", true, true, "userId", ctx.QueryParams(), &params.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SearchVoucher(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/v1/users", wrapper.CreateUser)
	router.GET(baseURL+"/v1/users/:id", wrapper.GetUserById)
	router.GET(baseURL+"/v1/vouchers", wrapper.SearchVoucher)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xWXW/bNhT9KwS3hw1QJCXphkJPW9eu8MOwoe72khgFTV1bbMSPXl45FQL/9+FSlm2l",
	"DpICxYAAezJFHd2Pcw4vfSe1t8E7cBRldSejbsCqtPwNQRH8HQH5KaAPgGQgvXPKAv9SH0BWMhIat5bb",
	"bSYRPnUGoZbV1YBaZCPKLz+CJrnN5BtEj+8gBu8ifBndQoxq/YQEI/BUjtOF69RU/UERP608Wl7JWhGc",
	"kbEgs/spM2nqCdY4+vnFAWccwRqQgU9jxdRyB82OyznVwz++081/3kYXAT98HfqBzjO5UW33VE4G7CH/",
	"cexHqOJ48JkAnWpfe504qiFqNIGMd7KSvxtXC9+RsB5BqCUv57dqzV1kssNWVrIhClVRxGE7Nz7R5lY+",
	"Ue4dKZ34BqsM461qb2BjdJNb+GXNm7n2lj+a5n7fmChMFHPAjdEglipCLbwT1ID4M4D79a+ZuMxLEQNo",
	"szJa8YfitjG6EQH9xtQgVp3TvK1aQ70gL6y6AeGxBoz5tbt2c29BdBFWXSta425ide3OxNX7BsTaUNMt",
	"BULw0ZDHfvED9xqrohhecd3F0M65bn5koQ21zPCOo33tZ8cFs2iAceiyzMv8nJv3AZwKRlbyMi/zS5nJ",
	"oKhJmhSb84JFTQ/BR/pSqGHscCMobg3tCagFK5HLFB8TQ7N6j0+nffAURHrl637UDFxKokJod8QWHyNn",
	"Gocdr75HWMlKflccpmGxG4XFUYLks2m5f/gaWpZDH+qWx+Ym7CC5fZh2qfOLsvxm5T1U2LzTGmJkN+wJ",
	"G6y5Ul1L3yz/dJafKCQB0omPnbUK+6nI7DW1jjwFBmcsGLo3SnFn6i3XsIYTZnkH1KGLg1tO++MtEFP0",
	"qp/VyYmoLFBy4NX9aIwTs9eSD72skmvHUV0NI2oqa3ZEEXxWNqQTc3F5fvHiJ5k9Pjy3i/99MfXFW6BB",
	"y2UvEuEPWGMzXIzxUWPsgA94Yw4KdTPesl/njk8dYH+wBxc4e2YWGRt/bi4ZZBu1PXLJuLNIgSLgZpRy",
	"csHzTVXmZfWyfFlKZnj3/X3F32wAe2qMW+/+MLDI+VTzKLeL7b8BAAD//whlccpGCwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
