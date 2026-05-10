package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	entity "quiz-app/internal/domain/entity"
)

type AuthUseCase struct {
	authServiceURL string
	httpClient     *http.Client
}

func NewAuthUseCase(authServiceURL string) *AuthUseCase {
	return &AuthUseCase{
		authServiceURL: authServiceURL,
		httpClient:     &http.Client{Timeout: 5 * time.Second},
	}
}

type verifyResponse struct {
	UserID     string   `json:"user_id"`
	EmailID    string   `json:"email_id"`
	Email      string   `json:"email"`
	Permission []string `json:"permission"`
}

func (uc *AuthUseCase) Authenticate(ctx context.Context, tokenString string) (entity.AuthClaims, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uc.authServiceURL+"/auth/verify", nil)
	if err != nil {
		return entity.AuthClaims{}, fmt.Errorf("failed to create verify request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err := uc.httpClient.Do(req)
	if err != nil {
		return entity.AuthClaims{}, fmt.Errorf("auth service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entity.AuthClaims{}, errors.New("unauthorized")
	}

	var result verifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return entity.AuthClaims{}, fmt.Errorf("failed to decode verify response: %w", err)
	}

	return entity.AuthClaims{
		UserID:     result.UserID,
		EmailID:    result.EmailID,
		AuthorMail: result.Email,
		Permission: result.Permission,
	}, nil
}
