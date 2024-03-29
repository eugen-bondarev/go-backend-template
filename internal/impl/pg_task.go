package impl

type PGTask struct {
	ID       int    `db:"id"`
	Title    string `db:"title"`
	AuthorID int    `db:"author_id"`
	Status   int    `db:"status"`
}
