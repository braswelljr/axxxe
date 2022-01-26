package helper

import (
  "os"
  "time"

  "github.com/braswelljr/goax/model"
  "github.com/golang-jwt/jwt/v4"
)

var (
	SecretKey = os.Getenv("SECRET_KEY")
)

type SignedParams struct {
	User model.TokenizedUserParams
	jwt.RegisteredClaims
}

func GetAllTokens(user model.TokenizedUserParams) (string, string, error) {
  //
	if SecretKey != "" {
		SecretKey = "xxyyzzaa"
	}

	// params
	signedParams := &SignedParams{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Local().Add(time.Hour * time.Duration(168)),
			},
		},
	}

	// refresh token
	refreshClaims := &SignedParams{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Local().Add(time.Hour * time.Duration(168)),
			},
		},
	}

	// create token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, signedParams).SignedString([]byte(SecretKey))

	// create refresh token
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SecretKey))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

// ValidateToken validates a token
func ValidateToken(token string) (*SignedParams, error) {
	if SecretKey != "" {
		SecretKey = "xxyyzzaa"
	}

	// parse token
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&SignedParams{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := tokenClaims.Claims.(*SignedParams)
	if !ok {
		return nil, err
	}

	// Ensure token is valid not expired
	if claims.VerifyExpiresAt(time.Now().Local(), true) == false {
		return nil, err
	}

	return claims, nil
}
