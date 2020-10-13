package gotruevault

import "fmt"

type URLBuilder interface {
	SearchDocumentURL(vaultID string) string
}

// DefaultURLBuilder  ...
type DefaultURLBuilder struct{}

// SearchDocumentURL ...
func (t *DefaultURLBuilder) SearchDocumentURL(vaultID string) string {
	return fmt.Sprintf("https://api.truevault.com/v1/vaults/%s/search", vaultID)
}
