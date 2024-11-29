package models

import "github.com/jackc/pgx/v5/pgtype"

type BlogParams struct {
	ID          pgtype.UUID `json:"id"`
	DomainID    pgtype.UUID `json:"domainId"`
	AccountID   pgtype.UUID `json:"accountId"`
	PageID      pgtype.UUID `json:"pageId"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	ImageInfo   string      `json:"imageInfo"`
	Keywords    []string    `json:"keywords"`
	IsActive    bool        `json:"isActive"`
}
