// Package helpers for storing all helper functions
package helpers

import (
	"net/http"
	"strings"
	"time"

	"github.com/odundlaw/cbt-backend/internal/jwt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a hashed password with a plaintext one.
func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func SetAuthCookies(w http.ResponseWriter, t *jwt.Tokens) {
	accessCookie := &http.Cookie{
		Name:     "access_cookie",
		Value:    t.Access,
		Expires:  t.ExpAcc,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_cookie",
		Value:    t.Refresh,
		Expires:  t.ExpRef,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode, // or Strict or None depending on frontend config
		Path:     "/",
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}

func ClearAuthCookies(w http.ResponseWriter) {
	expired := time.Now().Add(-time.Hour)

	accessCookie := &http.Cookie{
		Name:     "access_cookie",
		Value:    "",
		Expires:  expired,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_cookie",
		Value:    "",
		Expires:  expired,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
}

func BearerFromHeader(r *http.Request) string {
	h := r.Header.Get("Authorization")
	bearer, found := strings.CutPrefix(h, "Bearer ")
	if found {
		return bearer
	}

	return ""
}
