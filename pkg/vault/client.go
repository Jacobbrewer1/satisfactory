package vault

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	vault "github.com/hashicorp/vault/api"
)

var (
	ErrSecretNotFound = vault.ErrSecretNotFound
)

type ClientHandler interface {
	Client() *vault.Client
}

type Client interface {
	ClientHandler

	// GetKvSecretV2 returns a map of secrets for the given path.
	GetKvSecretV2(ctx context.Context, name string) (*vault.KVSecret, error)

	// GetSecret returns a map of secrets for the given path.
	GetSecret(ctx context.Context, path string) (*vault.Secret, error)

	// TransitEncrypt encrypts the given data.
	TransitEncrypt(ctx context.Context, data string) (*vault.Secret, error)

	// TransitDecrypt decrypts the given data.
	TransitDecrypt(ctx context.Context, data string) (string, error)
}

type (
	RenewalFunc func() (*vault.Secret, error)
	loginFunc   func(v *vault.Client) (*vault.Secret, error)
)

type client struct {
	ctx context.Context

	transitPathEncrypt string
	transitPathDecrypt string

	kvv2Mount string

	auth loginFunc

	// Below are set on initialization
	v         *vault.Client
	authCreds *vault.Secret
}

func NewClient(opts ...ClientOption) (Client, error) {
	c := new(client)

	for _, opt := range opts {
		opt(c)
	}

	if c.ctx == nil {
		c.ctx = context.Background()
	} else if c.v == nil {
		return nil, errors.New("vault client is nil")
	} else if c.auth == nil {
		return nil, errors.New("auth method is nil")
	}

	authCreds, err := c.auth(c.v)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate with Vault: %w", err)
	}

	c.authCreds = authCreds

	go c.renewAuthInfo()

	return c, nil
}

func (c *client) renewAuthInfo() {
	err := RenewLease(c.ctx, c, "auth", c.authCreds, func() (*vault.Secret, error) {
		authInfo, err := c.auth(c.v)
		if err != nil {
			return nil, fmt.Errorf("unable to renew auth info: %w", err)
		}

		c.authCreds = authInfo

		return authInfo, nil
	})
	if err != nil {
		slog.Error("unable to renew auth info", slog.String(logging.KeyError, err.Error()))
		os.Exit(1)
	}
}

func (c *client) Client() *vault.Client {
	return c.v
}

func (c *client) GetKvSecretV2(ctx context.Context, name string) (*vault.KVSecret, error) {
	secret, err := c.v.KVv2(c.kvv2Mount).Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}
	return secret, nil
}

func (c *client) GetSecret(ctx context.Context, path string) (*vault.Secret, error) {
	secret, err := c.v.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("unable to read secrets: %w", err)
	} else if secret == nil {
		return nil, ErrSecretNotFound
	}
	return secret, nil
}

func (c *client) TransitEncrypt(ctx context.Context, data string) (*vault.Secret, error) {
	plaintext := base64.StdEncoding.EncodeToString([]byte(data))

	// Encrypt the data using the transit engine
	encryptData, err := c.v.Logical().WriteWithContext(ctx, c.transitPathEncrypt, map[string]any{
		"plaintext": plaintext,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt data: %w", err)
	}

	return encryptData, nil
}

func (c *client) TransitDecrypt(ctx context.Context, data string) (string, error) {
	// Decrypt the data using the transit engine
	decryptData, err := c.v.Logical().WriteWithContext(ctx, c.transitPathDecrypt, map[string]any{
		"ciphertext": data,
	})
	if err != nil {
		return "", fmt.Errorf("unable to decrypt data: %w", err)
	}

	// Decode the base64 encoded data
	decodedData, err := base64.StdEncoding.DecodeString(decryptData.Data["plaintext"].(string))
	if err != nil {
		return "", fmt.Errorf("unable to decode data: %w", err)
	}

	return string(decodedData), nil
}
