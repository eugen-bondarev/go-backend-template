package impl

import (
	"errors"
	"go-backend-template/internal/model"

	"github.com/eugen-bondarev/go-slice-helpers/parallel"
)

type PGTaskRepo struct {
	pg         *Postgres
	taskMapper model.TaskMapper[PGTask]
}

func (taskRepo *PGTaskRepo) GetTasks() ([]model.Task, error) {
	var tasks []PGTask

	err := taskRepo.pg.GetDB().Select(&tasks, "SELECT * FROM tasks")

	if err != nil {
		return []model.Task{}, err
	}

	if len(tasks) == 0 {
		return []model.Task{}, errors.New("task not found")
	}

	return parallel.Map(tasks, taskRepo.taskMapper.ToTask), nil
}

func NewPGTaskRepo(pg *Postgres) model.TaskRepo {
	return &PGTaskRepo{
		pg:         pg,
		taskMapper: NewPGTaskMapper(),
	}
}
