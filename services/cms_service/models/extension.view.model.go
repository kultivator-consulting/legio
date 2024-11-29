package models

import "github.com/jackc/pgx/v5/pgtype"

type Extension struct {
	ID       pgtype.UUID `json:"id"`
	Name     string      `json:"name"`
	Slug     string      `json:"slug"`
	Icon     string      `json:"icon"`
	Data     string      `json:"data"`
	IsActive bool        `json:"isActive"`
}
