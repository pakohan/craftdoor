package lib

import "github.com/google/uuid"

// State contains information about the current card in front of the reader
type State struct {
	ID              uuid.UUID `json:"id"`
	IsCardAvailable bool      `json:"is_card_available"`
	CardData        []string  `json:"card_data"`
}
