package github

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateGithubAppJWT(appId int, privateKey string) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %w", err)
	}

	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),                       // Issued at: current time
		"exp": time.Now().Add(10 * time.Minute).Unix(), // Expiry: 10 minutes from iat
		"iss": appId,                                   // GitHub App ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return signedToken, nil
}
