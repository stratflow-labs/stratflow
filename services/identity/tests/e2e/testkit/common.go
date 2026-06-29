package testkit

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	DefaultUserPassword = "E2e-password-12345"
	DefaultTimeout      = 10 * time.Second
)

type AdminCredentials struct {
	Email    string
	Password string
}

type UserFixture struct {
	ID       string
	Login    string
	Email    string
	Password string
	Role     string
	Token    string
}

func LoadAdminCredentials() (AdminCredentials, bool) {
	email := strings.TrimSpace(os.Getenv("IDENTITY_E2E_ADMIN_EMAIL"))
	password := strings.TrimSpace(os.Getenv("IDENTITY_E2E_ADMIN_PASSWORD"))
	if email == "" || password == "" {
		return AdminCredentials{}, false
	}

	return AdminCredentials{
		Email:    email,
		Password: password,
	}, true
}

func NewUserFixture(role string) UserFixture {
	return UserFixture{
		Login:    UniqueLogin(role),
		Email:    fmt.Sprintf("%s.%s@example.test", role, UniqueSuffix(role)),
		Password: DefaultUserPassword,
		Role:     role,
	}
}

func UniqueSuffix(prefix string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", prefix, uuid.NewString()))
}

func UniqueLogin(prefix string) string {
	return strings.ToLower(fmt.Sprintf("%.10s-%s", prefix, uuid.NewString()[:12]))
}
