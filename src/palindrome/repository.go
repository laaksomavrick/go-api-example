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
		log.Print(err)
		return Message{}, err
	}

	row := r.db.QueryRowx( "SELECT * FROM messages WHERE id = $1", id)
	err = row.StructScan(&message)

	if err != nil {
		log.Print(err)
		return Message{}, err
	}

	return message, nil
}

func (r *Repository) GetMessage(id int) (Message, error) {
	var message Message

	row := r.db.QueryRowx( "SELECT * FROM messages WHERE id = $1", id)
	err := row.StructScan(&message)

	if err != nil {
		// No log here, as "no rows in result set" is most likely cause of err
		// *Could* be issues with the sql statement or db connection issues, but we'd likely see
		// other things blowing up in that case as well
		return Message{}, err
	}

	return message, nil
}

func (r *Repository) UpdateMessage(id int, content string) (Message, error) {
	message := Message{
		Content: content,
	}
	message.SetIsPalindrome()

	err := r.db.QueryRow(
		`UPDATE messages
				SET (content, is_palindrome) = ($1, $2) 
				WHERE id = $3
				RETURNING id`,
		message.Content,
		message.IsPalindrome,
		id,
	).Scan(&id)

	if err != nil {
		log.Print(err)
		return Message{}, err
	}

	row := r.db.QueryRowx( "SELECT * FROM messages WHERE id = $1", id)
	err = row.StructScan(&message)

	if err != nil {
		log.Print(err)
		return Message{}, err
	}

	return message, nil
}

func (r *Repository) DeleteMessage(id int) error {
	_, err := r.db.Exec(
		"DELETE FROM messages WHERE id = $1",
		id,
	)

	return err
}