package gen

import (
	"crypto/rand"
	"fmt"
)

func RandString(length int) string {
	buf := make([]byte, length)

	_, err := rand.Read(buf)

	if err != nil {
		fmt.Errorf("Error generating random string: %s", err)
	}

	return fmt.Sprintf("%x", buf)
}

func RandBytes(length int) []byte {
	buf := make([]byte, length)

	_, err := rand.Read(buf)

	if err != nil {
		fmt.Errorf("Error generating random string: %s", err)
	}

	return buf
}
