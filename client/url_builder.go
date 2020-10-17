package client

import "fmt"

//go:generate mockery --name URLBuilder
type URLBuilder interface {
	SearchDocumentURL(vaultID string) string
}

// defaultURLBuilder  ...
type defaultURLBuilder struct{}

// SearchDocumentURL ...
func (t *defaultURLBuilder) SearchDocumentURL(vaultID string) string {
	return fmt.Sprintf("https://api.truevault.com/v1/vaults/%s/search", vaultID)
}
