package reviews

import (
	"fmt"
	"time"
)

// GenerateID returns the current epoch time as a string to be used as an ID for a review
func GenerateID() string {
	now := time.Now()
	// Using Unix() for seconds since epoch, UnixNano() for nanoseconds since epoch
	epochMillis := now.UnixNano() / int64(time.Millisecond)
	reviewID := fmt.Sprintf("%d", epochMillis)
	return reviewID
}
