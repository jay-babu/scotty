// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7SST2/UMBDFv0oYOEab8E+qfCynHlClAlKl0oPrzDYujseMJwurVb47Gru7FbQsQoLb",
	"KPM87zcvswNHU6KIUTKYHWQ34mRL+R5ztreoZWJKyOKxNKaHhmwTgoEs7OMtLEsLjF9nzziAuToIr9u9",
	"kG7u0AksLXxwJLJ9N1o5i2mWv3FptUI7aWvAtZ2DgBGe8WBzQxTQRpXOGflfkZ7Pcgy11F5wKsULxjUY",
	"eN49BNzdp9vto10OdpbZbn+HlZ/gUqmPa6ohZMc+iaeojs35Bnnj8dvnCC2Il6AP6xrQwgY5V2W/6le9",
	"QlDCaJMHA6/LpxaSlbGs0eXyrtu87Nxo6/KU5bHtRxu+NEJN9SnWmpHV9tlwANAcoW6JWU5p2OooR1Ew",
	"lqk2peBdedbdZR29v8o/pfrrSZWMjkFCiTsnirn+vFd9/x9o7s/mCZyqaS4wBY/DM+V5UxF+1p3aobmo",
	"ia1U9Pby8rHoU8TvCZ3g0CAz8Uod1TMj6z8Hc7WDmQMYGEWS6bpAzoaRspiT/qSH5Xr5EQAA//9piEH/",
	"DwQAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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
