package paseto

import (
	"fmt"
	"paywise/internal/business/auth/token"
	tokenConfig "paywise/internal/business/auth/token"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// Paseto struct is a type which implements the TokenMaker interface
type Paseto struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func New(key string) (tokenConfig.TokenMaker, error) {
	// ensure that the length of the symmetric key is the standard of the paseto library
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key length")
	}
	return &Paseto{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(key),
	}, nil
}

func (p *Paseto) Create(username string, expiration time.Duration) (string, *token.Payload, error) {
	// create a payload for the token payload
	payload, err := tokenConfig.NewTokenPayload(username, expiration)
	if err != nil {
		return "", nil, err
	}
	token, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func (p *Paseto) Verify(token string) (*tokenConfig.Payload, error) {
	payload := tokenConfig.Payload{}
	err := p.paseto.Decrypt(token, p.symmetricKey, &payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error trying to decrypt the token | %v", err.Error())
	}

	// check if the payload is valid or not
	isValid := payload.Valid()

	if !isValid {
		return nil, fmt.Errorf("token is expired")
	}
	return &payload, nil
}
