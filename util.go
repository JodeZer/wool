package wool

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func decodeContentEncoding(reader io.Reader, encoding string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	if encoding == "gzip" {
		greader, err := gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(&b, greader); err != nil {
			return nil, err
		}
		return &b, nil
	}
	return nil, fmt.Errorf("unknown content-encoding %s", encoding)
}
