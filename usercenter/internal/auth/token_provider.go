package auth

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/square/go-jose.v2"
)

var (
	errAccessTokenNotRevocable = errors.New("not revocable")
	errInvalidSubject          = errors.New("invalid subject")
)

type TokenProvider interface {
	GrantToken(ctx context.Context, tokenInfo *TokenInfo, ttl time.Duration) (*GrantedToken, error)
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)
	RevokeToken(ctx context.Context, token string, def time.Duration) error
	RevokeAllTokens(ctx context.Context, uid uint64, deviceId string, excludes ...string) ([]string, error)

	GetPublicJWKSet(ctx context.Context, keyID string) jose.JSONWebKeySet
}

func invalidTokenError(err error) error {
	return status.Errorf(codes.Unauthenticated, "Invalid token: %s", err)
}
