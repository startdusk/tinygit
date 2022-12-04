package tinygit

import (
	"crypto/sha1"
	"fmt"
)

func sha1Hash(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}
