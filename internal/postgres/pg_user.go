package postgres

type PGUser struct {
	ID           int    `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
}
