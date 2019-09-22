package api

import "time"

type Message struct {
	Id int `db:"id" json:"id"`
	Content string `db:"content" json:"content"`
	IsPalindrome bool `db:"is_palindrome" json:"isPalindrome"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// Id and Content will be sent along
// Have to call function to set IsPalindrome

func (m *Message) SetIsPalindrome() {
	// TODO: unit test
	// https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
	r := []rune(m.Content)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	reversed := string(r)
	m.IsPalindrome = reversed == m.Content
}

