package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	entity "quiz-app/internal/domain/entity"
	"strconv"
	"time"
)

// UserUseCase proxies user management to auth-service. quiz-app does not own
// user data — Postgres in auth-service is the single source of truth.
// See CLAUDE.md §E#1, §L3.
type UserUseCase struct {
	authServiceURL string
	httpClient     *http.Client
}

func NewUserUseCase(authServiceURL string) *UserUseCase {
	return &UserUseCase{
		authServiceURL: authServiceURL,
		httpClient:     &http.Client{Timeout: 5 * time.Second},
	}
}

// GetMe returns the current user's profile. Token is forwarded to auth-service.
func (uc *UserUseCase) GetMe(ctx context.Context, bearerToken string) (entity.User, error) {
	var user entity.User
	if err := uc.do(ctx, http.MethodGet, "/auth/me", bearerToken, nil, &user); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// UpdateMe applies a partial self-update.
func (uc *UserUseCase) UpdateMe(ctx context.Context, bearerToken string, req entity.UpdateMeRequest) (entity.User, error) {
	var user entity.User
	if err := uc.do(ctx, http.MethodPatch, "/auth/me", bearerToken, req, &user); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// ListUsers returns the paginated user list (admin only — enforced at auth-service).
func (uc *UserUseCase) ListUsers(ctx context.Context, bearerToken string, page, limit int) (entity.PaginatedUsers, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	path := fmt.Sprintf("/admin/users?page=%d&limit=%d", page, limit)
	var result entity.PaginatedUsers
	if err := uc.do(ctx, http.MethodGet, path, bearerToken, nil, &result); err != nil {
		return entity.PaginatedUsers{}, err
	}
	return result, nil
}

// AdminUpdateUser applies admin-only fields (permission, name) to a target user.
func (uc *UserUseCase) AdminUpdateUser(ctx context.Context, bearerToken, userID string, req entity.AdminUpdateUserRequest) (entity.User, error) {
	var user entity.User
	path := "/admin/users/" + userID
	if err := uc.do(ctx, http.MethodPatch, path, bearerToken, req, &user); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

type upstreamError struct {
	Status int
	Body   string
}

func (e *upstreamError) Error() string {
	return fmt.Sprintf("auth-service %d: %s", e.Status, e.Body)
}

// UpstreamStatus extracts the HTTP status code from a usecase error, if any.
// Returns 0 when the error did not originate from an upstream response.
func UpstreamStatus(err error) int {
	var ue *upstreamError
	if errors.As(err, &ue) {
		return ue.Status
	}
	return 0
}

func (uc *UserUseCase) do(ctx context.Context, method, path, bearerToken string, body, out any) error {
	if uc.authServiceURL == "" {
		return errors.New("auth service URL not configured")
	}
	if bearerToken == "" {
		return errors.New("missing bearer token")
	}

	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
		reader = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, uc.authServiceURL+path, reader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := uc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return &upstreamError{Status: resp.StatusCode, Body: string(raw)}
	}

	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// Tiny helper for callers that need to format limit/page from raw query — exported
// to avoid duplicating strconv in handlers.
func ParseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return def
	}
	return n
}
