package roles

const (
	AdminRole   = "admin"
	UserRole    = "user"
)

var validRoles = map[string]bool{
	AdminRole:   true,
	UserRole:      true,
}

func IsValidRole(role string) bool {
	return validRoles[role]
}

func IsAdminRole(role string) bool {
  return role == AdminRole
}

func IsUserRole(role string) bool {
	return role == UserRole
}
