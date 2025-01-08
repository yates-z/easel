package rbac

import "sync"

// Enforcer manages policies and performs permission checks.
type Enforcer struct {
	policies []Policy
	adapter  Adapter
	mu       sync.RWMutex
}

// NewEnforcer creates a new Enforcer with the given adapter.
func NewEnforcer(adapter Adapter) (*Enforcer, error) {
	policies, err := adapter.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return &Enforcer{
		policies: policies,
		adapter:  adapter,
	}, nil
}

// AddPolicy adds a new policy to the enforcer.
func (e *Enforcer) AddPolicy(role, resource, action string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	// avoid adding duplicate policies.
	for _, p := range e.policies {
		if p.Role == role && p.Resource == resource && p.Action == action {
			return
		}
	}
	e.policies = append(e.policies, Policy{Role: role, Resource: resource, Action: action})
}

// RemovePolicy removes a policy from the enforcer.
func (e *Enforcer) RemovePolicy(role, resource, action string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.policies = filterPolicies(e.policies, func(p Policy) bool {
		return !(p.Role == role && p.Resource == resource && p.Action == action)
	})
}

func (e *Enforcer) Save() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.adapter.SavePolicy(e.policies)
}

// Enforce checks if a role has permission to perform an action on a resource.
func (e *Enforcer) Enforce(role, resource, action string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, p := range e.policies {
		if p.Role == role && p.Resource == resource && p.Action == action {
			return true
		}
	}
	return false
}

// filterPolicies filters a slice of policies based on a condition.
func filterPolicies(policies []Policy, condition func(Policy) bool) []Policy {
	var result []Policy
	for _, p := range policies {
		if condition(p) {
			result = append(result, p)
		}
	}
	return result
}
