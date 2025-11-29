package main

import (
	"bytes"
	"compress/gzip"
)

func compressToGzip(data []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)
	if _, err := gw.Write(data); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
