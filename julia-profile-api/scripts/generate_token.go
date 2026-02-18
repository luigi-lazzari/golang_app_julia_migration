package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/generate_token.go <userID> [secret]")
		os.Exit(1)
	}

	userID := os.Args[1]
	secret := "your-default-secret-for-local-testing"
	if len(os.Args) > 2 {
		secret = os.Args[2]
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
		"iss": "bff-julia",
		"aud": "julia-app",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated JWT for UserID: %s\n", userID)
	fmt.Printf("Token: %s\n\n", tokenString)
	fmt.Printf("Curl example:\n")
	fmt.Printf("curl -H \"Authorization: Bearer %s\" http://localhost:8090/api/v1/user/profile\n", tokenString)
}
