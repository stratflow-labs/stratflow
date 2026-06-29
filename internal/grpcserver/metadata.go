package grpcserver

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	HeaderAuthorization = "authorization"
	HeaderRequestID     = "x-request-id"
	HeaderCorrelationID = "x-correlation-id"
)

var (
	ErrAuthorizationHeaderNotFound = errors.New("grpc authorization header not found")
	ErrAuthorizationSchemaInvalid  = errors.New("grpc authorization schema is invalid")
	ErrAuthorizationTokenEmpty     = errors.New("grpc authorization token is empty")
)

func BearerTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrAuthorizationHeaderNotFound
	}

	values := md.Get(HeaderAuthorization)
	if len(values) == 0 {
		return "", ErrAuthorizationHeaderNotFound
	}

	header := strings.TrimSpace(values[0])
	if !strings.HasPrefix(header, "Bearer ") {
		return "", ErrAuthorizationSchemaInvalid
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if token == "" {
		return "", ErrAuthorizationTokenEmpty
	}

	return token, nil
}
