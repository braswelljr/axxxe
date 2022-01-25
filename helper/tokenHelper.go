package helper

import (
	"github.com/braswelljr/goax/model"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

var (
	SECRET_KEY = os.Getenv("SECRET_KEY")
)

type SignedParams struct {
	user model.TokenizedUserParams
	jwt.RegisteredClaims
}

func GetAllTokens(user model.TokenizedUserParams) (string, string, error) {
	if SECRET_KEY != "" {
		SECRET_KEY = "xxyyzzaa"
	}

	// params
	signedParams := &SignedParams{
		user: user,
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
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, signedParams).SignedString([]byte(SECRET_KEY))

	// create refresh token
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}
