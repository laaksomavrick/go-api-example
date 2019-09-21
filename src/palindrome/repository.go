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

func (r *Repository) CreateMessage(content string) (Message, error) {
	var id int

	message := Message{
		Content: content,
	}
	message.SetIsPalindrome()

	err := r.db.QueryRow(
			"INSERT INTO messages (content, is_palindrome) VALUES ($1, $2) RETURNING id",
			message.Content,
			message.IsPalindrome,
		).Scan(&id)

	if err != nil {
		return Message{}, err
	}

	row := r.db.QueryRowx( "SELECT * FROM messages WHERE id = $1", id)
	err = row.StructScan(&message)

	if err != nil {
		return Message{}, err
	}

	return message, nil
}