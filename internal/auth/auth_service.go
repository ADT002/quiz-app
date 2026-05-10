package auth

import (
	"context"
	"errors"
	"fmt"
	entity "quiz-app/internal/domain/entity"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/option"
)

type AuthService struct {
	firebaseClient *auth.Client
	jwtSecret      []byte
}

func NewAuthService(firebaseConfigFile string, jwtSecret []byte) (*AuthService, error) {
	// jsonString := os.Getenv("FIREBASE_CONFIG")
	// opt := option.WithCredentialsJSON([]byte(jsonString))

	opt := option.WithCredentialsFile(firebaseConfigFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	firebaseClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &AuthService{firebaseClient: firebaseClient, jwtSecret: jwtSecret}, nil
}

func (s *AuthService) CreateJWT(claims entity.AuthClaims) (string, error) {
	tokenClaims := jwt.MapClaims{
		"user_id":    claims.UserID,
		"email_id":   claims.AuthorMail, // Sử dụng "email_id" để giữ nguyên cấu trúc cũ, nhưng giá trị vẫn là email
		"autho_mail": claims.AuthorMail,
		"exp":        time.Now().Add(time.Hour * time.Duration(claims.Exp)).Unix(),
		"permission": claims.Permission,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateJWT(tokenString string) (*jwt.Token, *entity.AuthClaims, error) {
	claims := &entity.AuthClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}
			return s.jwtSecret, nil
		},
	)

	if err != nil || !token.Valid {
		return nil, nil, errors.New("invalid token")
	}

	return token, claims, nil
}

func (s *AuthService) VerifyIDToken(ctx context.Context, token string) (*auth.Token, error) {
	return s.firebaseClient.VerifyIDToken(ctx, token)
}
