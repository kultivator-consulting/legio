package models

import "github.com/jackc/pgx/v5/pgtype"

type Page struct {
	ID             pgtype.UUID        `json:"id"`
	Created        pgtype.Timestamptz `json:"created"`
	Modified       pgtype.Timestamptz `json:"modified"`
	Deleted        pgtype.Timestamptz `json:"deleted"`
	DomainID       pgtype.UUID        `json:"domainId"`
	AccountID      pgtype.UUID        `json:"accountId"`
	ContentID      pgtype.UUID        `json:"contentId"`
	PagePathID     pgtype.UUID        `json:"pagePathId"`
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

type PageTemplate struct {
	ID               pgtype.UUID        `json:"id"`
	Created          pgtype.Timestamptz `json:"created"`
	Modified         pgtype.Timestamptz `json:"modified"`
	Deleted          pgtype.Timestamptz `json:"deleted"`
	DomainID         pgtype.UUID        `json:"domainId"`
	AccountID        pgtype.UUID        `json:"accountId"`
	ContentID        pgtype.UUID        `json:"contentId"`
	ParentPagePathID pgtype.UUID        `json:"parentPagePathId"`
	Title            string             `json:"title"`
	Slug             string             `json:"slug"`
	Description      string             `json:"description"`
	IsActive         bool               `json:"isActive"`
}
