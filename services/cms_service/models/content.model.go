package models

import (
	"cortex_api/database/db_gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type ContentComponentField struct {
	ID           pgtype.UUID `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	DataType     string      `json:"dataType"`
	EditorType   string      `json:"editorType"`
	Validation   string      `json:"validation"`
	DefaultValue string      `json:"defaultValue"`
	IsActive     bool        `json:"isActive"`
}

type ContentComponent struct {
	ID                  pgtype.UUID             `json:"id"`
	Name                string                  `json:"name"`
	Icon                string                  `json:"icon"`
	Description         string                  `json:"description"`
	ClassName           string                  `json:"className"`
	HtmlTag             string                  `json:"htmlTag"`
	ChildTagConstraints []string                `json:"childTagConstraints"`
	IsActive            bool                    `json:"isActive"`
	Fields              []ContentComponentField `json:"fields"`
}

type ContentCollection struct {
	ID        pgtype.UUID `json:"id"`
	ParentID  pgtype.UUID `json:"parentId"`
	ContentID pgtype.UUID `json:"contentId"`
	Ordering  int32       `json:"ordering"`
	IsActive  bool        `json:"isActive"`
}

type Content struct {
	ID        pgtype.UUID      `json:"id"`
	DomainID  pgtype.UUID      `json:"domainId"`
	AccountID pgtype.UUID      `json:"accountId"`
	Title     string           `json:"title"`
	Slug      string           `json:"slug"`
	Data      string           `json:"data"`
	IsActive  bool             `json:"isActive"`
	Component ContentComponent `json:"component"`
	Children  []*Content       `json:"children"`
}

type BreadCrumb struct {
	ID       pgtype.UUID `json:"id"`
	ParentID pgtype.UUID `json:"parentId"`
	Title    string      `json:"title"`
	Slug     string      `json:"slug"`
}

type PageContent struct {
	Page        db_gen.Page  `json:"page"`
	Content     *Content     `json:"content"`
	BreadCrumbs []BreadCrumb `json:"breadcrumbs"`
}
