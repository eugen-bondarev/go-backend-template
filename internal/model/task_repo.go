package model

type TaskRepo interface {
	GetTasks() ([]Task, error)
}
