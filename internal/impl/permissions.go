package impl

type Policy struct {
	Object  string
	Action  string
	Subject string
}

func NewPolicy(subject, action, object string) Policy {
	return Policy{
		Subject: subject,
		Action:  action,
		Object:  object,
	}
}

type Policies struct {
	policies []Policy
}

func NewPolicies() Policies {
	return Policies{
		policies: make([]Policy, 0),
	}
}

func (p *Policies) Add(subject, action, object string) {
	p.policies = append(p.policies, NewPolicy(subject, action, object))
}

func (p *Policies) RoleCan(subject, action, object string) bool {
	userPolicy := NewPolicy(subject, action, object)
	for _, policy := range p.policies {
		if policy == userPolicy {
			return true
		}
	}
	return false
}
