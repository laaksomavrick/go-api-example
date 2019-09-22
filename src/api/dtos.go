package api

import "errors"

// UpsertMessageDto defines the shape of create and update requests against a message.
// Since the dto is so simple between create and updates, we can use the same one.
// More complicated shapes would require separate createFooDto and updateFooDto structs.
type UpsertMessageDto struct {
	Content string `json:"content"`
}

// Validate performs business logic checks against the DTO, confirming that the data conforms to our
// business rules
func (dto *UpsertMessageDto) Validate() error {
	if dto.Content == "" || len(dto.Content) > 1024 {
		return errors.New("content must be between 1 and 1024 characters")
	}

	return nil
}