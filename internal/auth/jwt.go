package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrEmptyJWTSecret = errors.New("jwt secret is empty")
	ErrInvalidToken   = errors.New("invalid jwt token")
	ErrExpiredToken   = errors.New("jwt token expired")
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

type Claims struct {
	UserID    int64
	ExpiresAt time.Time
	IssuedAt  time.Time
}

func NewJWTManager(secret string, ttl time.Duration) (*JWTManager, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, ErrEmptyJWTSecret
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	return &JWTManager{
		secret: []byte(secret),
		ttl:    ttl,
	}, nil
}

func (m *JWTManager) Generate(userID int64) (string, error) {
	now := time.Now().UTC()
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	payload := map[string]any{
		"sub": strconv.FormatInt(userID, 10),
		"iat": now.Unix(),
		"exp": now.Add(m.ttl).Unix(),
	}

	headerPart, err := encodeJSON(header)
	if err != nil {
		return "", err
	}
	payloadPart, err := encodeJSON(payload)
	if err != nil {
		return "", err
	}

	unsigned := headerPart + "." + payloadPart
	signature := m.sign(unsigned)

	return unsigned + "." + signature, nil
}

func (m *JWTManager) Parse(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSignature := m.sign(unsigned)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSignature)) {
		return nil, ErrInvalidToken
	}

	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, ErrInvalidToken
	}
	if header.Algorithm != "HS256" || header.Type != "JWT" {
		return nil, ErrInvalidToken
	}

	var payload struct {
		Subject   string `json:"sub"`
		ExpiresAt int64  `json:"exp"`
		IssuedAt  int64  `json:"iat"`
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, ErrInvalidToken
	}

	userID, err := strconv.ParseInt(payload.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return nil, ErrInvalidToken
	}

	expiresAt := time.Unix(payload.ExpiresAt, 0).UTC()
	if !time.Now().UTC().Before(expiresAt) {
		return nil, ErrExpiredToken
	}

	return &Claims{
		UserID:    userID,
		ExpiresAt: expiresAt,
		IssuedAt:  time.Unix(payload.IssuedAt, 0).UTC(),
	}, nil
}

func (m *JWTManager) TTL() time.Duration {
	return m.ttl
}

func (m *JWTManager) sign(unsigned string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func encodeJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal jwt part: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}
