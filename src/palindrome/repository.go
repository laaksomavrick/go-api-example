package palindrome

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:db,
	}
}

func (r *Repository) GetMessages() ([]Message, error) {
	messages := []Message{}

	err := r.db.Select(&messages, "SELECT * FROM messages")

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return messages, nil
}