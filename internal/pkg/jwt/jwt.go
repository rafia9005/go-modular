package jwt

import (
	"fmt"
	"sync"
	"time"

	gojwt "github.com/golang-jwt/jwt"
)

// JWT interface defines methods for handling JWT operations.
type JWT interface {
	GenerateToken(data map[string]interface{}) (string, error)
	ValidateToken(token string) (bool, error)
	ParseToken(tokenString string) (map[string]interface{}, error)
	Logout(tokenString string, blacklist *Blacklist) error
}

// JWTImpl implements the JWT interface.
type JWTImpl struct {
	SignatureKey string
	Expiration   int
}

// Blacklist manages invalidated tokens.
type Blacklist struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

// NewJWTImpl creates a new instance of JWTImpl.
func NewJWTImpl(signatureKey string, expiration int) JWT {
	return &JWTImpl{SignatureKey: signatureKey, Expiration: expiration}
}

// GenerateToken generates a JWT with the given data.
func (j *JWTImpl) GenerateToken(data map[string]interface{}) (string, error) {
	var mySigningKey = []byte(j.SignatureKey)
	token := gojwt.New(gojwt.SigningMethodHS256)
	claims := token.Claims.(gojwt.MapClaims)

	for key, value := range data {
		claims[key] = value
	}

	expirationDuration := time.Duration(j.Expiration) * 24 * time.Hour
	expirationTime := time.Now().Add(expirationDuration).Unix()
	claims["exp"] = expirationTime

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken checks if a JWT is valid and not expired.
func (j *JWTImpl) ValidateToken(tokenString string) (bool, error) {
	token, err := gojwt.Parse(tokenString, func(token *gojwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SignatureKey), nil
	})

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(gojwt.MapClaims)
	if !ok || !token.Valid {
		return false, nil
	}

	expirationTime := claims["exp"].(float64)
	if time.Now().Unix() > int64(expirationTime) {
		return false, nil
	}

	return true, nil
}

// ParseToken extracts claims from a JWT.
func (j *JWTImpl) ParseToken(tokenString string) (map[string]interface{}, error) {
	token, err := gojwt.Parse(tokenString, func(token *gojwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SignatureKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(gojwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// NewBlacklist creates a new instance of Blacklist.
func NewBlacklist() *Blacklist {
	return &Blacklist{
		tokens: make(map[string]time.Time),
	}
}

// Add adds a token to the blacklist with an expiration time.
func (b *Blacklist) Add(token string, expiration time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens[token] = expiration
}

// IsBlacklisted checks if a token is blacklisted.
func (b *Blacklist) IsBlacklisted(token string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	expiration, exists := b.tokens[token]
	if !exists {
		return false
	}

	// Remove expired tokens from the blacklist.
	if time.Now().After(expiration) {
		delete(b.tokens, token)
		return false
	}

	return true
}

// Logout adds a token to the blacklist, effectively logging it out.
func (j *JWTImpl) Logout(tokenString string, blacklist *Blacklist) error {
	blacklist.Add(tokenString, time.Now().Add(1*time.Hour))
	return nil
}
