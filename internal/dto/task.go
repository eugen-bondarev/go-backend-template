package dto

type Task struct {
	ID        int32
	Title     string
	Status    int32
	getAuthor func() User
}

func NewTask(id int32, title string, status int32, getAuthor func() User) Task {
	return Task{
		ID:        id,
		Title:     title,
		Status:    status,
		getAuthor: getAuthor,
	}
}

func (t Task) Author() User {
	return t.getAuthor()
}
