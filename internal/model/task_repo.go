package model

type TaskRepo interface {
	GetTasks() ([]Task, error)
	GetTaskByID(id int) (Task, error)
}
