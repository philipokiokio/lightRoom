package schemas

import "github.com/google/uuid"

type TagPayload struct {
	Title string `json:"title"`
}

type TagProfile struct {
	TagPayload
	Id uuid.UUID `json:"id"`
}

type PortfolioPayload struct {
	Title           string      `json:"title"`
	Description     *string     `json:"description"`
	Price           int         `json:"price"`
	Tags            []uuid.UUID `json:"tags"`
	PaywalledImages []string    `json:"paywalled_images"`
	Images          []string    `json:"images"`
}
