package signature

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
)

type Length = uint8
type Algo = string

const (
	Short Length = 8
	Long  Length = 32
)

const (
	SHA1   Algo = "sha1"
	SHA256 Algo = "sha256"
)

func Sign(content string, secret string, algo Algo) ([]byte, error) {
	var hashFunc hash.Hash
	switch algo {
	case SHA1:
		hashFunc = sha1.New()
	case SHA256:
		hashFunc = sha256.New()
	default:
		return nil, errors.New("unsupported signature algorithm")
	}

	hashFunc.Write([]byte(content + secret))

	return hashFunc.Sum(nil), nil
}

func SignURL(content string, secret string, algo Algo, length uint8) string {
	rawSignature, _ := Sign(content, secret, algo)
	signature := base64.RawURLEncoding.EncodeToString(rawSignature)

	return fmt.Sprintf("s--%s--", signature[:length])
}
