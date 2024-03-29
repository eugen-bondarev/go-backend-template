package impl

import "go-backend-template/internal/model"

type PGTaskMapper struct {
}

func NewPGTaskMapper() model.TaskMapper[PGTask] {
	return &PGTaskMapper{}
}

func (taskMapper *PGTaskMapper) FromTask(task model.Task) PGTask {
	return PGTask{
		ID:       task.ID,
		Title:    task.Title,
		AuthorID: task.AuthorID,
		Status:   task.Status,
	}
}

func (taskMapper *PGTaskMapper) ToTask(task PGTask) model.Task {
	return model.Task{
		ID:       task.ID,
		Title:    task.Title,
		AuthorID: task.AuthorID,
		Status:   task.Status,
	}
}
