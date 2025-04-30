package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	GithubAppID      = 123456
	GithubPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
<YOUR_PRIVATE_KEY_CONTENT_HERE>
-----END RSA PRIVATE KEY-----`
)

func generateGithubAppJWT(appID int, privateKey string) (string, error) {
	// Parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %w", err)
	}

	// Prepare the claims for the JWT
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),                       // Issued at: current time
		"exp": time.Now().Add(10 * time.Minute).Unix(), // Expiry: 10 minutes from iat
		"iss": appID,                                   // GitHub App ID
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token using the private key
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return signedToken, nil
}
