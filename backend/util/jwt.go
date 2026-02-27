package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess      = "access"
	TokenTypeRefresh     = "refresh"
	TokenTypeTOTPPending = "totp_pending"

	minSecretLength = 64
)

type TOTPSessionClaims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateTOTPSessionToken(cfg *JWTConfig, userID string) (string, error) {
	if err := ValidateSecret(cfg.Secret); err != nil {
		return "", err
	}

	jti, err := generateJTI()
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP session JTI: %w", err)
	}

	now := time.Now()
	claims := TOTPSessionClaims{
		UserID:    userID,
		TokenType: TokenTypeTOTPPending,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{cfg.Issuer},
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign TOTP session token: %w", err)
	}
	return signed, nil
}

func ValidateTOTPSessionToken(tokenStr, secret, issuer string) (*TOTPSessionClaims, error) {
	if err := ValidateSecret(secret); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &TOTPSessionClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(secret), nil
	},
		jwt.WithIssuer(issuer),
		jwt.WithAudience(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("TOTP session token validation failed: %w", err)
	}

	claims, ok := token.Claims.(*TOTPSessionClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid TOTP session token claims")
	}

	if claims.TokenType != TokenTypeTOTPPending {
		return nil, errors.New("token is not a TOTP session token")
	}

	if claims.Subject == "" {
		return nil, errors.New("TOTP session token subject is empty")
	}

	return claims, nil
}

type AccessTokenClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type JWTConfig struct {
	Secret              string
	Issuer              string
	AccessExpiryMinutes int
	RefreshExpiryDays   int
}

func ValidateSecret(secret string) error {
	if len(secret) < minSecretLength {
		return fmt.Errorf("JWT secret must be at least %d characters for HS512, got %d", minSecretLength, len(secret))
	}
	return nil
}

func GenerateSecureToken(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateTokenPair(cfg *JWTConfig, userID, email, role string) (*TokenPair, error) {
	if err := ValidateSecret(cfg.Secret); err != nil {
		return nil, err
	}

	signingKey := []byte(cfg.Secret)
	now := time.Now()

	accessJTI, err := generateJTI()
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token JTI: %w", err)
	}

	accessClaims := AccessTokenClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{cfg.Issuer},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.AccessExpiryMinutes) * time.Minute)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        accessJTI,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	accessTokenStr, err := accessToken.SignedString(signingKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshJTI, err := generateJTI()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token JTI: %w", err)
	}

	refreshClaims := RefreshTokenClaims{
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{cfg.Issuer},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.RefreshExpiryDays) * 24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        refreshJTI,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(signingKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    cfg.AccessExpiryMinutes * 60,
	}, nil
}

func ValidateAccessToken(tokenStr, secret, issuer string) (*AccessTokenClaims, error) {
	if err := ValidateSecret(secret); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(secret), nil
	},
		jwt.WithIssuer(issuer),
		jwt.WithAudience(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if claims.TokenType != TokenTypeAccess {
		return nil, errors.New("token is not an access token")
	}

	if claims.Subject == "" {
		return nil, errors.New("token subject is empty")
	}

	if claims.ID == "" {
		return nil, errors.New("token JTI is empty")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenStr, secret, issuer string) (*RefreshTokenClaims, error) {
	if err := ValidateSecret(secret); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &RefreshTokenClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(secret), nil
	},
		jwt.WithIssuer(issuer),
		jwt.WithAudience(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if claims.TokenType != TokenTypeRefresh {
		return nil, errors.New("token is not a refresh token")
	}

	if claims.Subject == "" {
		return nil, errors.New("token subject is empty")
	}

	if claims.ID == "" {
		return nil, errors.New("token JTI is empty")
	}

	return claims, nil
}
