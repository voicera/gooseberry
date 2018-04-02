package rest

import (
	"io"
)

type noopResponseDecoder struct{}

func (noopResponseDecoder) DecodeResponse(body io.ReadCloser, result interface{}) error {
	return nil
}
