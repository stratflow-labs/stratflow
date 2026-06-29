package auth

import tx "github.com/stratflow-labs/stratflow/internal/foundation/tx"

type Service struct {
	user        CredentialFinder
	password    PasswordVerifier
	token       TokenService
	accessToken AccessTokenRevoker
	txManager   tx.Manager
}

func NewService(
	user CredentialFinder,
	password PasswordVerifier,
	token TokenService,
	accessToken AccessTokenRevoker,
	txManager tx.Manager,
) *Service {
	return &Service{
		user:        user,
		password:    password,
		token:       token,
		accessToken: accessToken,
		txManager:   txManager,
	}
}
