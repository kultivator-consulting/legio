package models

import "github.com/jackc/pgx/v5/pgtype"

type PagePath struct {
	ID               pgtype.UUID        `json:"id"`
	Created          pgtype.Timestamptz `json:"created"`
	Modified         pgtype.Timestamptz `json:"modified"`
	Deleted          pgtype.Timestamptz `json:"deleted"`
	DomainID         pgtype.UUID        `json:"domainId"`
	AccountID        pgtype.UUID        `json:"accountId"`
	ParentPagePathID pgtype.UUID        `json:"parentPagePathId"`
	Title            string             `json:"title"`
	Slug             string             `json:"slug"`
	IsActive         bool               `json:"isActive"`
	Folders          []PagePath         `json:"folders"`
	Pages            []Page             `json:"pages"`
	Templates        []PageTemplate     `json:"templates"`
	Extensions       []Extension        `json:"extensions"`
}

type PagePathLink struct {
	ID         pgtype.UUID `json:"id"`
	PagePathID pgtype.UUID `json:"pagePathId"`
	Title      string      `json:"title"`
	Link       string      `json:"link"`
}
