package models

// All returns every model that participates in the schema. The core service
// owns migrations so the schema has a single source of truth.
func All() []any {
	return []any{
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
		&BlacklistedToken{},
	}
}
