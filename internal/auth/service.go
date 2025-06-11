package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Service handles generation and refreshing of JWT tokens.
type Service struct {
	Secret []byte
	// AccessTokenTTL defines expiration for access tokens
	AccessTokenTTL time.Duration
	// RefreshTokenTTL defines expiration for refresh tokens
	RefreshTokenTTL time.Duration
}

func NewService(secret []byte) *Service {
	return &Service{
		Secret:          secret,
		AccessTokenTTL:  30 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
}

// GenerateToken returns signed access and refresh JWT tokens for given subject.
func (s *Service) GenerateToken(sub string, roles []string) (string, string, error) {
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   sub,
		"roles": roles,
		"exp":   now.Add(s.AccessTokenTTL).Unix(),
	})
	tokenStr, err := t.SignedString(s.Secret)
	if err != nil {
		return "", "", err
	}
	r := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   sub,
		"roles": roles,
		"exp":   now.Add(s.RefreshTokenTTL).Unix(),
	})
	refreshStr, err := r.SignedString(s.Secret)
	if err != nil {
		return "", "", err
	}
	return tokenStr, refreshStr, nil
}

// Refresh validates refresh token and issues new tokens.
func (s *Service) Refresh(refreshToken string) (string, string, error) {
	tok, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return s.Secret, nil
	})
	if err != nil || !tok.Valid {
		return "", "", err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", jwt.ErrTokenMalformed
	}
	sub, _ := claims["sub"].(string)
	rolesIface, _ := claims["roles"].([]interface{})
	roles := make([]string, 0, len(rolesIface))
	for _, r := range rolesIface {
		if s, ok := r.(string); ok {
			roles = append(roles, s)
		}
	}
	return s.GenerateToken(sub, roles)
}
