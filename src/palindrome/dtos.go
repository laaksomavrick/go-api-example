package palindrome

import "errors"

type CreateMessageDto struct {
	Content string `json:"content"`
}

func (dto *CreateMessageDto) Validate() error {
	if dto.Content == "" || len(dto.Content) > 1024 {
		return errors.New("content must be between 0 and 1024 characters")
	}

	return nil
}