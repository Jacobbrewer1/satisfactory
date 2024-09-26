package repositories

import (
	"context"

	"github.com/Jacobbrewer1/satisfactory/pkg/vault"
	"github.com/spf13/viper"
)

type ConnectionOption func(c *databaseConnector)

func WithVaultClient(client vault.Client) ConnectionOption {
	return func(c *databaseConnector) {
		c.client = client
	}
}

func WithViper(v *viper.Viper) ConnectionOption {
	return func(c *databaseConnector) {
		c.vip = v
	}
}

func WithCurrentSecrets(secrets *vault.Secrets) ConnectionOption {
	return func(c *databaseConnector) {
		c.currentSecrets = secrets
	}
}

func WithContext(ctx context.Context) ConnectionOption {
	return func(c *databaseConnector) {
		c.ctx = ctx
	}
}
