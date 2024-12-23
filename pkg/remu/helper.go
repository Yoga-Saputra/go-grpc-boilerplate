package remu

import "strings"

// Helper function return redis key with namespaces
func (r *Remu) getFullKey(key string) string {
	return strings.TrimSpace(r.nameSpace) + ":" + strings.TrimSpace(key)
}
