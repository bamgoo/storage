package storage

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func hashFile(file string) (hash string, hex string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", "", err
	}
	sum := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum), fmt.Sprintf("%x", sum), nil
}
