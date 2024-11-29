package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ContentView struct {
	ID        pgtype.UUID   `json:"id"`
	DomainId  pgtype.UUID   `json:"domainId"`
	Component ComponentView `json:"component"`
	AccountId pgtype.UUID   `json:"accountId"`
	Title     string        `json:"title"`
	Slug      string        `json:"slug"`
	Data      string        `json:"data"`
	IsActive  bool          `json:"isActive"`
}
