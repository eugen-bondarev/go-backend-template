package model

type TaskMapper[T any] interface {
	FromTask(Task) T
	ToTask(T) Task
}
