package rbac

// Policy defines a role-based access control policy.
type Policy struct {
	Role     string
	Resource string
	Action   string
}
