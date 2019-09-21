package palindrome

import "errors"

// Since the dto is so simple between create and updates, we can use the same one
// More complicated shapes would require separate createFooDto and updateFooDto structs
type UpsertMessageDto struct {
	Content string `json:"content"`
}

func (dto *UpsertMessageDto) Validate() error {
	if dto.Content == "" || len(dto.Content) > 1024 {
		return errors.New("content must be between 0 and 1024 characters")
	}

	return nil
}