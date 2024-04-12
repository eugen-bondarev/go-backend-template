package permissions

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

type Subject struct {
	subject string
}

func NewSubject(subject string) Subject {
	return Subject{subject: subject}
}

func (s *Subject) Can(action, object string) [3]string {
	return [3]string{s.subject, action, object}
}

func NewPolicies() Policies {
	return Policies{
		policies: make([]Policy, 0),
	}
}

func (p *Policies) Add(policy [3]string) {
	p.policies = append(p.policies, NewPolicy(policy[0], policy[1], policy[2]))
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
