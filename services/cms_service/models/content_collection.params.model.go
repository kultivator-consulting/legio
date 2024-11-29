package models

import "github.com/jackc/pgx/v5/pgtype"

type CreateContentCollectionParams struct {
	ContentID pgtype.UUID `json:"contentId"`
	Ordering  int32       `json:"ordering"`
	IsActive  bool        `json:"isActive"`
}

type UpdateContentCollectionParams struct {
	Ordering int32 `json:"ordering"`
	IsActive bool  `json:"isActive"`
}
