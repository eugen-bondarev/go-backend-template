package impl

type Permission struct {
	Object  string
	Action  string
	Subject string
}

func NewPermission(subject, action, object string) Permission {
	return Permission{
		Subject: subject,
		Action:  action,
		Object:  object,
	}
}

type Permissions struct {
	permissions []Permission
}

func NewPermissions() Permissions {
	return Permissions{
		permissions: make([]Permission, 0),
	}
}

func (p *Permissions) Add(subject, action, object string) {
	p.permissions = append(p.permissions, NewPermission(subject, action, object))
}

func (p *Permissions) RoleCan(subject, action, object string) bool {
	userPermission := NewPermission(subject, action, object)
	for _, permission := range p.permissions {
		if permission == userPermission {
			return true
		}
	}
	return false
}
