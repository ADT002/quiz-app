package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	EmailAddress string    `json:"author_mail"`
	EmailID      string    `json:"email_id"`
	Password     string    `json:"-"`
	Permission   []string  `json:"permission"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthClaims struct {
	UserID     string   `json:"user_id"`
	EmailID    string   `json:"email_id"`
	AuthorMail string   `json:"author_mail"`
	Exp        int64    `json:"exp"`
	Permission []string `json:"permission"`
	jwt.RegisteredClaims
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Password  *string `json:"password,omitempty"`
}

// NewUser constructs a new User entity and performs validation.
func NewUser(user_id, firstName, lastName, email, password string, permission []string) (User, error) {
	if strings.TrimSpace(firstName) == "" || strings.TrimSpace(lastName) == "" {
		return User{}, errors.New("first and last name cannot be empty")
	}
	if !isValidEmail(email) {
		return User{}, errors.New("invalid email address")
	}
	if len(password) < 6 {
		return User{}, errors.New("password must be at least 6 characters long")
	}

	return User{
		ID:           primitive.NewObjectID().Hex(),
		FirstName:    firstName,
		LastName:     lastName,
		EmailAddress: strings.TrimSpace(email),
		Password:     hashPassword(password), // Hash the password
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(), // Set initial updated_at time
		Permission:   permission,
	}, nil
}

// UpdatePassword allows updating the user's password with hashing.
func (u *User) UpdatePassword(newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	u.Password = hashPassword(newPassword)
	u.UpdatedAt = time.Now() // Update the updated_at field
	return nil
}

// Auxiliary functions for validation and hashing

// Placeholder for hashing function
func hashPassword(password string) string {
	// Implement the actual password hashing here
	return password // Replace with a proper hash function
}

// Simple email validation
func isValidEmail(email string) bool {
	// Implement email validation logic (can use regex or external libraries)
	return strings.Contains(email, "@")
}
func permissionTeacher() []string {
	return []string{
		// Question
		"CREATE_QUESTION",
		"UPDATE_QUESTION",
		"DELETE_QUESTION",
		"VIEW_QUESTION",

		// Test
		"CREATE_TEST",
		"UPDATE_TEST",
		"DELETE_TEST",
		"VIEW_TEST",
		"ASSIGN_TEST_TO_CLASS",

		// Class
		"CREATE_CLASS",
		"UPDATE_CLASS",
		"DELETE_CLASS",
		"VIEW_CLASS",

		// User submit & result
		"VIEW_USER_SUBMIT",
		"VIEW_TEST_RESULT",
		"EXPORT_TEST_RESULT",

		// Topic & Level
		"CREATE_TOPIC",
		"UPDATE_TOPIC",
		"DELETE_TOPIC",
		"VIEW_TOPIC",

		"CREATE_LEVEL",
		"UPDATE_LEVEL",
		"DELETE_LEVEL",
		"VIEW_LEVEL",

		// Dashboard
		"VIEW_DASHBOARD",
	}
}
