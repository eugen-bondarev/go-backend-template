package dto

import "go-backend-template/internal/model"

type TaskResolvers struct {
	GetAuthor func() User
}

type Task struct {
	ID        int32
	Title     string
	Status    int32
	resolvers TaskResolvers
}

func TaskFromModel(m model.Task, resolvers TaskResolvers) Task {
	return Task{
		ID:        int32(m.ID),
		Title:     m.Title,
		Status:    int32(m.Status),
		resolvers: resolvers,
	}
}

func (t Task) Author() User {
	return t.resolvers.GetAuthor()
}
