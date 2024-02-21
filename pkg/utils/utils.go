package utils

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

func ToPascalCase(s string) string {
	var pascalCase strings.Builder
	nextToUpper := true

	s = strings.TrimSpace(s)

	for _, r := range s {
		if nextToUpper {
			pascalCase.WriteRune(unicode.ToUpper(r))
			nextToUpper = false
		} else if r == ' ' {
			nextToUpper = true
		} else {
			pascalCase.WriteRune(r)
		}
	}

	return pascalCase.String()
}

// CalculateAverageRating takes a slice of float64 values (ratings) and returns the average
func CalculateAverageRating(ratings []float64) float64 {
	if len(ratings) == 0 {
		return 0
	}

	var sum float64
	for _, rating := range ratings {
		sum += rating
	}

	return sum / float64(len(ratings))
}

func GenerateReviewID() string {
	now := time.Now()
	// Using Unix() for seconds since epoch, UnixNano() for nanoseconds since epoch
	epochMillis := now.UnixNano() / int64(time.Millisecond)
	reviewID := fmt.Sprintf("%d", epochMillis)
	return reviewID
}

// Install "firebase.google.com/go"
//    "firebase.google.com/go/auth"
//    "google.golang.org/api/option"
// TODO - Verify token
//func VerifyToken(token string) (*jwt.Token, error) {
//	opt := option.WithCredentialsFile("path/to/your/firebase-adminsdk.json")
//	app, err := firebase.NewApp(context.Background(), nil, opt)
//	if err != nil {
//		return nil, fmt.Errorf("error initializing app: %v", err)
//	}
//
//	// Get an auth client from the Firebase App
//	ctx := context.Background()
//	client, err := app.Auth(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("error getting Auth client: %v", err)
//	}
//
//	// Verify the ID token
//	token, err := client.VerifyIDToken(ctx, idToken)
//	if err != nil {
//		return nil, fmt.Errorf("error verifying ID token: %v", err)
//	}
//
//	return token, nil
//}
//}
