package metadata

import (
	"encoding/base64"
	"strings"
)

const (
	binHdrSuffix = "-bin"
)

func encodeKeyValue(k, v string) (string, string) {
	k = strings.ToLower(k)
	if strings.HasSuffix(k, binHdrSuffix) {
		return k, base64.StdEncoding.EncodeToString([]byte(v))
	}
	return k, v
}
