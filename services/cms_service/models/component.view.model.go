package models

import "github.com/jackc/pgx/v5/pgtype"

type ComponentView struct {
	ID                  pgtype.UUID `json:"id"`
	Name                string      `json:"name"`
	Icon                string      `json:"icon"`
	Description         string      `json:"description"`
	ClassName           string      `json:"className"`
	HtmlTag             string      `json:"htmlTag"`
	ChildTagConstraints []string    `json:"childTagConstraints"`
	IsActive            bool        `json:"isActive"`
}
