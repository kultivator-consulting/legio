package models

import "github.com/jackc/pgx/v5/pgtype"

type CreatePageParams struct {
	DomainID       pgtype.UUID        `json:"domainId"`
	AccountID      pgtype.UUID        `json:"accountId"`
	ContentID      pgtype.UUID        `json:"contentId"`
	Title          string             `json:"title"`
	Slug           string             `json:"slug"`
	SeoTitle       string             `json:"seoTitle"`
	SeoDescription string             `json:"seoDescription"`
	SeoKeywords    string             `json:"seoKeywords"`
	DraftPageID    pgtype.UUID        `json:"draftPageId"`
	PageTemplateID pgtype.UUID        `json:"pageTemplateId"`
	PublishAt      pgtype.Timestamptz `json:"publishAt"`
	UnpublishAt    pgtype.Timestamptz `json:"unpublishAt"`
	Version        int32              `json:"version"`
	IsActive       bool               `json:"isActive"`
}

type CreatePagePathParams struct {
	DomainID  pgtype.UUID `json:"domainId"`
	AccountID pgtype.UUID `json:"accountId"`
	Title     string      `json:"title"`
	Slug      string      `json:"slug"`
	IsActive  bool        `json:"isActive"`
}

type CreatePageTemplateParams struct {
	DomainID    pgtype.UUID `json:"domainId"`
	AccountID   pgtype.UUID `json:"accountId"`
	ContentID   pgtype.UUID `json:"contentId"`
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	Description string      `json:"description"`
	IsActive    bool        `json:"isActive"`
}
