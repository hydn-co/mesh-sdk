package permissions

import (
	"fmt"
	"strings"
)

type Permission struct {
	Scope     string `json:"scope"`     // tenant ID or "global"
	Resource  string `json:"resource"`  // resource like "user/abc123"
	Operation string `json:"operation"` // action like "read", "write", "manage"
}

// String formats the permission as "scope::resource::operation"
func (p Permission) String() string {
	return fmt.Sprintf("%s::%s::%s", p.Scope, p.Resource, p.Operation)
}

// ParsePermission parses a permission string into a Permission struct.
func ParsePermission(s string) (Permission, error) {
	parts := strings.Split(s, "::")
	if len(parts) != 3 {
		return Permission{}, fmt.Errorf("invalid permission format: %s", s)
	}
	return Permission{
		Scope:     parts[0],
		Resource:  parts[1],
		Operation: parts[2],
	}, nil
}
