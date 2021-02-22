package schedule

import (
	"crypto/sha256"
	"fmt"
)

func generateID(s fmt.Stringer) string {
	return fmt.Sprintf(
		"%x",
		sha256.Sum256([]byte(s.String())),
	)
}
