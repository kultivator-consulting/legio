package v1

import (
	"context"
	"cortex_api/common/api_models"
	"cortex_api/common/api_utils"
	"cortex_api/database/db_gen"
	"cortex_api/services/cms_service/models"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"strings"
)

type PagePathModel struct {
	Queries *db_gen.Queries
	DbCtx   context.Context
}

type PagePathInterface interface {
	CreatePagePath(ctx *fiber.Ctx) error
	ListPagePaths(ctx *fiber.Ctx) error
	GetPagePath(ctx *fiber.Ctx) error
	DeletePagePath(ctx *fiber.Ctx) error
	UpdatePagePath(ctx *fiber.Ctx) error
	AddPage(ctx *fiber.Ctx) error
	UpdatePage(ctx *fiber.Ctx) error
	RemovePage(ctx *fiber.Ctx) error
	GetPage(ctx *fiber.Ctx) error
	ValidatePagePath(ctx *fiber.Ctx) error
	CreateDraftPage(ctx *fiber.Ctx) error
}

func PagePathController(
	queries *db_gen.Queries,
	dbCtx context.Context,
) *PagePathModel {
	return &PagePathModel{
		Queries: queries,
		DbCtx:   dbCtx,
	}
}

func (controller *PagePathModel) GetPageLink(pageID pgtype.UUID) (string, error) {
	page, dbErr := controller.Queries.GetPageById(controller.DbCtx, pageID)
	if dbErr != nil {
		log.Printf("GetPageLink error while fetching object, error %v\n", dbErr)
		return "", dbErr
	}

	parentPagePath, dbErr := controller.Queries.GetPagePathById(controller.DbCtx, page.PagePathID)
	if dbErr != nil {
		log.Printf("GetPageLink error while fetching object, error %v\n", dbErr)
		return "", dbErr
	}

	var path []string

	pagePath := parentPagePath
	for pagePath.ID.Valid {
		if pagePath.Slug != "index" {
			path = append(path, pagePath.Slug)
		}
		pagePath, dbErr = controller.Queries.GetPagePathById(controller.DbCtx, pagePath.ParentPagePathID)
		if dbErr != nil {
			if dbErr.Error() == "no rows in result set" {
				break
			}
			log.Printf("GetPageLink error while fetching object, error %v\n", dbErr)
			return "", dbErr
		}
	}

	for i := 0; i < len(path)/2; i++ {
		opp := len(path) - i - 1
		path[i], path[opp] = path[opp], path[i]
	}
	path = append(path, page.Slug)

	return fmt.Sprintf("/%s", strings.Join(path, "/")), nil
}

func (controller *PagePathModel) ListChildLinks(linkPath []string, domainId pgtype.UUID, parentId *pgtype.UUID) ([]models.PagePathLink, error) {
	dbErr := error(nil)

	pagePaths := make([]db_gen.PagePath, 0)
	if parentId == nil {
		pagePaths, dbErr = controller.Queries.ListLinksAtRoot(controller.DbCtx, domainId)
	} else {
		pagePaths, dbErr = controller.Queries.ListLinksByParentId(controller.DbCtx, db_gen.ListLinksByParentIdParams{
			DomainID:         domainId,
			ParentPagePathID: *parentId,
		})
	}
	if dbErr != nil {
		log.Printf("ListLinks error while retrieving object, error %v\n", dbErr)
		return nil, dbErr
	}

	links := make([]models.PagePathLink, 0)
	for _, folderPath := range pagePaths {
		pathLink := folderPath.Slug
		if folderPath.Slug == "index" {
			pathLink = ""
		}

		currentPath := append(linkPath, pathLink)

		calculatedLink := strings.Join(currentPath, "/")
		if len(currentPath) == 1 {
			calculatedLink = "/"
		}

		links = append(links, models.PagePathLink{
			ID:    folderPath.ID,
			Title: folderPath.Title,
			Link:  calculatedLink,
		})

		childLinks, err := controller.ListChildLinks(currentPath, domainId, &folderPath.ID)
		if err != nil {
			return nil, err
		}

		links = append(links, childLinks...)

		for _, childLink := range childLinks {
			childPages, dbErr := controller.Queries.ListPagesByPagePathId(controller.DbCtx, childLink.ID)
			if dbErr != nil {
				log.Printf("ListLinks error while retrieving object, error %v\n", dbErr)
				return nil, dbErr
			}

			for _, page := range childPages {
				pathPageLink := page.Slug
				if page.Slug == "index" {
					pathPageLink = ""
				}

				pageLink := models.PagePathLink{
					ID:         page.ID,
					PagePathID: page.PagePathID,
					Title:      childLink.Title + " | " + page.Title,
					Link:       childLink.Link + "/" + pathPageLink,
				}
				links = append(links, pageLink)
			}
		}
	}

	return links, nil
}

func (controller *PagePathModel) CreateDraftPage(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pageUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("UpdatePage error parsing sourcePage ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid sourcePage ID",
		})
		return ctx.Send(payload)
	}

	sourcePage, dbErr := controller.Queries.GetPageById(controller.DbCtx, pgtype.UUID{Bytes: pageUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdatePage error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to update sourcePage",
		})
		return ctx.Send(payload)
	}

	clonedContent, dbErr := ContentController(controller.Queries, controller.DbCtx).ClonePageContent(sourcePage.ContentID, ctx, sourcePage.Title, sourcePage.Slug)
	if dbErr != nil {
		log.Printf("CreateDraftPage error while cloning content, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to clone content",
		})
		return ctx.Send(payload)
	}

	dbPage := db_gen.CreatePageParams{
		DomainID:       sourcePage.DomainID,
		AccountID:      sourcePage.AccountID,
		ContentID:      clonedContent.ID,
		Title:          sourcePage.Title,
		Slug:           sourcePage.Slug,
		SeoTitle:       sourcePage.SeoTitle,
		SeoDescription: sourcePage.SeoDescription,
		SeoKeywords:    sourcePage.SeoKeywords,
		DraftPageID:    pgtype.UUID{Bytes: pageUuid, Valid: true},
		PageTemplateID: sourcePage.PageTemplateID,
		IsActive:       sourcePage.IsActive,
	}

	newPage, dbErr := controller.Queries.CreatePage(controller.DbCtx, dbPage)
	if dbErr != nil {
		log.Printf("AddPage error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to add sourcePage",
		})
		return ctx.Send(payload)
	}

	pageView := models.Page{
		ID:             newPage.ID,
		Created:        newPage.Created,
		Modified:       newPage.Modified,
		Deleted:        newPage.Deleted,
		DomainID:       newPage.DomainID,
		AccountID:      newPage.AccountID,
		ContentID:      newPage.ContentID,
		PagePathID:     newPage.PagePathID,
		Title:          newPage.Title,
		Slug:           newPage.Slug,
		SeoTitle:       newPage.SeoTitle,
		SeoDescription: newPage.SeoDescription,
		SeoKeywords:    newPage.SeoKeywords,
		DraftPageID:    newPage.DraftPageID,
		PageTemplateID: newPage.PageTemplateID,
		PublishAt:      newPage.PublishAt,
		UnpublishAt:    newPage.UnpublishAt,
		Version:        newPage.Version,
		IsActive:       newPage.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   pageView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: marshalErr.Error(),
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) CreatePagePath(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pagePath := models.CreatePagePathParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &pagePath)
	if unmarshalErr != nil {
		log.Printf("CreatePagePath error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	createPageParams := db_gen.CreatePagePathAndReturnIdParams{}
	dbErr := error(nil)
	if (ctx.Params("parentId") != "null") && (ctx.Params("parentId") != "") {
		parentUuid, err := uuid.Parse(ctx.Params("parentId"))
		if err != nil {
			log.Printf("CreatePagePath error parsing parent page path ID, error %v\n", err.Error())
			ctx.Status(http.StatusBadRequest)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "Invalid parent folder ID",
			})
			return ctx.Send(payload)
		}

		_, dbErr = controller.Queries.GetPagePathByDomainIdParentPagePathIdAndSlug(controller.DbCtx, db_gen.GetPagePathByDomainIdParentPagePathIdAndSlugParams{
			DomainID:         pagePath.DomainID,
			Slug:             pagePath.Slug,
			ParentPagePathID: pgtype.UUID{Bytes: parentUuid, Valid: true},
		})

		createPageParams = db_gen.CreatePagePathAndReturnIdParams{
			DomainID:         pagePath.DomainID,
			AccountID:        pagePath.AccountID,
			ParentPagePathID: pgtype.UUID{Bytes: parentUuid, Valid: true},
			Title:            pagePath.Title,
			Slug:             pagePath.Slug,
			IsActive:         pagePath.IsActive,
		}
	} else {
		_, dbErr = controller.Queries.GetPagePathByDomainIdRootPathAndSlug(controller.DbCtx, db_gen.GetPagePathByDomainIdRootPathAndSlugParams{
			DomainID: pagePath.DomainID,
			Slug:     pagePath.Slug,
		})

		createPageParams = db_gen.CreatePagePathAndReturnIdParams{
			DomainID:  pagePath.DomainID,
			AccountID: pagePath.AccountID,
			Title:     pagePath.Title,
			Slug:      pagePath.Slug,
			IsActive:  pagePath.IsActive,
		}
	}
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("CreatePagePath error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to create folder",
			})
			return ctx.Send(payload)
		}
	}
	if dbErr == nil {
		log.Printf("CreatePagePath error while fetching object, error pagePath already exists\n")
		ctx.Status(http.StatusConflict)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusConflict,
			ErrorMessage: "Folder already exists",
		})
		return ctx.Send(payload)
	}

	pagePathUuid, dbErr := controller.Queries.CreatePagePathAndReturnId(controller.DbCtx, createPageParams)
	if dbErr != nil {
		log.Printf("CreatePagePath error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to create folder",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   pagePathUuid,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: marshalErr.Error(),
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) AddPage(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	page := models.CreatePageParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &page)
	if unmarshalErr != nil {
		log.Printf("AddPage error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	parentUuid, err := uuid.Parse(ctx.Params("parentId"))
	if err != nil {
		log.Printf("AddPage error parsing parent page path ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid parent page path ID",
		})
		return ctx.Send(payload)
	}

	_, dbErr := controller.Queries.GetPageByPagePathIdAndSlug(controller.DbCtx, db_gen.GetPageByPagePathIdAndSlugParams{
		PagePathID: pgtype.UUID{Bytes: parentUuid, Valid: true},
		Slug:       page.Slug,
	})
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("AddPage error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to add content block",
			})
			return ctx.Send(payload)
		}
	}
	if !page.DraftPageID.Valid { // If draft page ID is not provided, it means this is a new page and should not overlap slug with existing pages
		if dbErr == nil {
			log.Printf("AddPage error while fetching object, error page already exists for this parent\n")
			ctx.Status(http.StatusConflict)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusConflict,
				ErrorMessage: "A page with this name already exists in this folder",
			})
			return ctx.Send(payload)
		}
	}

	dbPage := db_gen.CreatePageParams{
		PagePathID:     pgtype.UUID{Bytes: parentUuid, Valid: true},
		DomainID:       page.DomainID,
		AccountID:      page.AccountID,
		ContentID:      page.ContentID,
		Title:          page.Title,
		Slug:           page.Slug,
		SeoTitle:       page.SeoTitle,
		SeoDescription: page.SeoDescription,
		SeoKeywords:    page.SeoKeywords,
		DraftPageID:    page.DraftPageID,
		PageTemplateID: page.PageTemplateID,
		Version:        1,
		IsActive:       page.IsActive,
	}

	pageResult, dbErr := controller.Queries.CreatePage(controller.DbCtx, dbPage)
	if dbErr != nil {
		log.Printf("AddPage error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to add page",
		})
		return ctx.Send(payload)
	}

	pageView := models.Page{
		pageResult.ID,
		pageResult.Created,
		pageResult.Modified,
		pageResult.Deleted,
		pageResult.DomainID,
		pageResult.AccountID,
		pageResult.ContentID,
		pageResult.PagePathID,
		pageResult.Title,
		pageResult.Slug,
		pageResult.SeoTitle,
		pageResult.SeoDescription,
		pageResult.SeoKeywords,
		pageResult.DraftPageID,
		pageResult.PageTemplateID,
		pageResult.PublishAt,
		pageResult.UnpublishAt,
		pageResult.Version,
		pageResult.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   pageView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: marshalErr.Error(),
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) UpdatePage(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pageUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("UpdatePage error parsing page ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid page ID",
		})
		return ctx.Send(payload)
	}

	sourcePage, dbErr := controller.Queries.GetPageById(controller.DbCtx, pgtype.UUID{Bytes: pageUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdatePage error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to update page",
		})
		return ctx.Send(payload)
	}

	page := db_gen.UpdatePageByIdParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &page)
	if unmarshalErr != nil {
		log.Printf("UpdatePage error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	pendingPage, dbErr := controller.Queries.GetPageByPagePathIdAndSlug(controller.DbCtx, db_gen.GetPageByPagePathIdAndSlugParams{
		PagePathID: page.PagePathID,
		Slug:       page.Slug,
	})
	if dbErr == nil {
		if pendingPage.ID != sourcePage.ID {
			log.Printf("UpdatePage error page already exists for this path\n")
			ctx.Status(http.StatusConflict)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusConflict,
				ErrorMessage: "A page with this name already exists in this folder",
			})
			return ctx.Send(payload)
		}
	} else {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("UpdatePage error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to update page",
			})
			return ctx.Send(payload)
		}
	}

	// Values that are not allowed to be changed when updating content collection children
	page.ID = sourcePage.ID
	page.Version = sourcePage.Version + 1

	updatedPage, dbErr := controller.Queries.UpdatePageById(controller.DbCtx, page)
	if dbErr != nil {
		log.Printf("UpdatePage error while updating object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to perform update",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   updatedPage,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("UpdatePage error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) RemovePage(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pageUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("RemovePage error parsing page ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid page ID",
		})
		return ctx.Send(payload)
	}

	page, dbErr := controller.Queries.GetPageById(controller.DbCtx, pgtype.UUID{Bytes: pageUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Child page not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("RemovePage error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove page",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeletePageById(controller.DbCtx, page.ID)
	if dbErr != nil {
		log.Printf("RemovePage error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove page",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("RemovePage error serialising response, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) GetPage(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pageUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetPage error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	page, dbErr := controller.Queries.GetPageById(controller.DbCtx, pgtype.UUID{Bytes: pageUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetPage error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve page",
		})
		return ctx.Send(payload)
	}

	pageView := models.Page{
		page.ID,
		page.Created,
		page.Modified,
		page.Deleted,
		page.DomainID,
		page.AccountID,
		page.ContentID,
		page.PagePathID,
		page.Title,
		page.Slug,
		page.SeoTitle,
		page.SeoDescription,
		page.SeoKeywords,
		page.DraftPageID,
		page.PageTemplateID,
		page.PublishAt,
		page.UnpublishAt,
		page.Version,
		page.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   pageView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetPage error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) ValidatePagePath(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	found := true

	pagePathUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("ValidatePagePath error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	pagePath, dbErr := controller.Queries.GetPagePathById(controller.DbCtx, pgtype.UUID{Bytes: pagePathUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			found = false
		} else {
			log.Printf("ValidatePagePath error while retrieving object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve page path",
			})
			return ctx.Send(payload)
		}
	} else {
		_, dbErr = controller.Queries.GetPageByPagePathIdAndSlug(controller.DbCtx, db_gen.GetPageByPagePathIdAndSlugParams{
			PagePathID: pagePath.ID,
			Slug:       ctx.Params("slug"),
		})
		if dbErr != nil {
			if dbErr.Error() == "no rows in result set" {
				found = false
			} else {
				log.Printf("ValidatePagePath error while retrieving object, error %v\n", dbErr)
				ctx.Status(http.StatusBadGateway)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to retrieve page",
				})
				return ctx.Send(payload)
			}
		}
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   found,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("ValidatePagePath error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) AddPageTemplate(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	page := models.CreatePageTemplateParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &page)
	if unmarshalErr != nil {
		log.Printf("AddPageTemplate error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	parentFolderUuid := pgtype.UUID{}
	parentFolder := &parentFolderUuid
	if (ctx.Params("parentId") != "null") && (ctx.Params("parentId") != "") {
		parentUuid, err := uuid.Parse(ctx.Params("parentId"))
		if err != nil {
			log.Printf("AddPageTemplate error parsing parent page path ID, error %v\n", err.Error())
			ctx.Status(http.StatusBadRequest)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "Invalid parent page path ID",
			})
			return ctx.Send(payload)
		}

		parentFolderUuid = pgtype.UUID{Bytes: parentUuid, Valid: true}

		_, dbErr := controller.Queries.GetPageTemplateByPagePathIdAndSlug(controller.DbCtx, db_gen.GetPageTemplateByPagePathIdAndSlugParams{
			ParentPagePathID: pgtype.UUID{Bytes: parentUuid, Valid: true},
			Slug:             page.Slug,
		})
		if dbErr != nil {
			if dbErr.Error() != "no rows in result set" {
				log.Printf("AddPageTemplate error while fetching object, error %v\n", dbErr)
				ctx.Status(http.StatusBadGateway)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to add content block",
				})
				return ctx.Send(payload)
			}
		}
		if dbErr == nil {
			log.Printf("AddPageTemplate error while fetching object, error page already exists for this parent\n")
			ctx.Status(http.StatusConflict)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusConflict,
				ErrorMessage: "A page with this name already exists in this folder",
			})
			return ctx.Send(payload)
		}
	} else {
		_, dbErr := controller.Queries.GetPageTemplateAtRootAndBySlug(controller.DbCtx, page.Slug)
		if dbErr != nil {
			if dbErr.Error() != "no rows in result set" {
				log.Printf("AddPageTemplate error while fetching object, error %v\n", dbErr)
				ctx.Status(http.StatusBadGateway)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to add content block",
				})
				return ctx.Send(payload)
			}
		}
		if dbErr == nil {
			log.Printf("AddPageTemplate error while fetching object, error page already exists for this parent\n")
			ctx.Status(http.StatusConflict)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusConflict,
				ErrorMessage: "A page with this name already exists in this folder",
			})
			return ctx.Send(payload)
		}
	}

	dbPageTemplate := db_gen.CreatePageTemplateAndReturnIdParams{
		ParentPagePathID: *parentFolder,
		DomainID:         page.DomainID,
		AccountID:        page.AccountID,
		ContentID:        page.ContentID,
		Title:            page.Title,
		Slug:             page.Slug,
		Description:      page.Description,
		IsActive:         page.IsActive,
	}

	pageTemplateUuid, dbErr := controller.Queries.CreatePageTemplateAndReturnId(controller.DbCtx, dbPageTemplate)
	if dbErr != nil {
		log.Printf("AddPageTemplate error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to add page template",
		})
		return ctx.Send(payload)
	}

	pageTemplateResult, dbErr := controller.Queries.GetPageTemplateById(controller.DbCtx, pageTemplateUuid)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page template not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("AddPageTemplate error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve page template",
		})
		return ctx.Send(payload)
	}

	pageTemplateView := models.PageTemplate{
		pageTemplateResult.ID,
		pageTemplateResult.Created,
		pageTemplateResult.Modified,
		pageTemplateResult.Deleted,
		pageTemplateResult.DomainID,
		pageTemplateResult.AccountID,
		pageTemplateResult.ContentID,
		pageTemplateResult.ParentPagePathID,
		pageTemplateResult.Title,
		pageTemplateResult.Slug,
		pageTemplateResult.Description,
		pageTemplateResult.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   pageTemplateView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: marshalErr.Error(),
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) UpdatePageTemplate(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pageTemplateUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("UpdatePageTemplate error parsing page ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid page template ID",
		})
		return ctx.Send(payload)
	}

	sourcePageTemplate, dbErr := controller.Queries.GetPageTemplateById(controller.DbCtx, pgtype.UUID{Bytes: pageTemplateUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page template not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdatePageTemplate error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to update page template",
		})
		return ctx.Send(payload)
	}

	pageTemplate := db_gen.UpdatePageTemplateByIdParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &pageTemplate)
	if unmarshalErr != nil {
		log.Printf("UpdatePageTemplate error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	pendingPageTemplate, dbErr := controller.Queries.GetPageTemplateByPagePathIdAndSlug(controller.DbCtx, db_gen.GetPageTemplateByPagePathIdAndSlugParams{
		ParentPagePathID: pageTemplate.ParentPagePathID,
		Slug:             pageTemplate.Slug,
	})
	if dbErr == nil {
		if pendingPageTemplate.ID != sourcePageTemplate.ID {
			log.Printf("UpdatePageTemplate error page already exists for this path\n")
			ctx.Status(http.StatusConflict)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusConflict,
				ErrorMessage: "A page template with this name already exists in this folder",
			})
			return ctx.Send(payload)
		}
	} else {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("UpdatePageTemplate error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to update page",
			})
			return ctx.Send(payload)
		}
	}

	// Values that are not allowed to be changed when updating content collection children
	pageTemplate.ID = sourcePageTemplate.ID

	updatedPage, dbErr := controller.Queries.UpdatePageTemplateById(controller.DbCtx, pageTemplate)
	if dbErr != nil {
		log.Printf("UpdatePageTemplate error while updating object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to perform update",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   updatedPage,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("UpdatePageTemplate error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) RemovePageTemplate(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pageTemplateUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("RemovePageTemplate error parsing pageTemplate ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid pageTemplate ID",
		})
		return ctx.Send(payload)
	}

	pageTemplate, dbErr := controller.Queries.GetPageTemplateById(controller.DbCtx, pgtype.UUID{Bytes: pageTemplateUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Child page template not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("RemovePageTemplate error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove page template",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeleteContent(controller.DbCtx, pageTemplate.ContentID)
	if dbErr != nil {
		log.Printf("RemovePageTemplate error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove page template",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DetachTemplateFromPageByTemplateId(controller.DbCtx, pageTemplate.ID)
	if dbErr != nil {
		log.Printf("RemovePageTemplate error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove page template",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeletePageTemplateById(controller.DbCtx, pageTemplate.ID)
	if dbErr != nil {
		log.Printf("RemovePageTemplate error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove page template",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("RemovePageTemplate error serialising response, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) GetPageTemplate(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pageTemplateUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetPageTemplate error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	page, dbErr := controller.Queries.GetPageTemplateById(controller.DbCtx, pgtype.UUID{Bytes: pageTemplateUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page template not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetPageTemplate error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve page template",
		})
		return ctx.Send(payload)
	}

	pageTemplateView := models.PageTemplate{
		page.ID,
		page.Created,
		page.Modified,
		page.Deleted,
		page.DomainID,
		page.AccountID,
		page.ContentID,
		page.ParentPagePathID,
		page.Title,
		page.Slug,
		page.Description,
		page.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   pageTemplateView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetPageTemplate error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) ListPageTemplates(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	requestedPage, requestedPageSize, sortBy, sortOrder := api_utils.GetPaginationParams(ctx)
	rowCount, dbErr := controller.Queries.CountPageTemplates(controller.DbCtx)
	if dbErr != nil {
		log.Printf("ListPageTemplates error while retrieving count, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve page template count",
		})
		return ctx.Send(payload)
	}

	query, err := api_utils.ParseQueryString(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid query",
		})
		return ctx.Send(payload)
	}

	domainId, err := uuid.Parse(query.Get("domain"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid query, invalid domain ID",
		})
		return ctx.Send(payload)
	}

	pageTemplates := make([]db_gen.PageTemplate, 0)
	if (query.Get("path") != "null") && (query.Get("path") != "") {
		pagePathPathId, err := uuid.Parse(query.Get("path"))
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "Invalid query, invalid page path ID",
			})
			return ctx.Send(payload)
		}

		if sortOrder == "desc" {
			pageTemplates, dbErr = controller.Queries.ListPageTemplatesByDomainIdAndPagePathIdDesc(controller.DbCtx, db_gen.ListPageTemplatesByDomainIdAndPagePathIdDescParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				ParentPagePathID:  pgtype.UUID{Bytes: pagePathPathId, Valid: true},
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		} else {
			pageTemplates, dbErr = controller.Queries.ListPageTemplatesByDomainIdAndPagePathIdAsc(controller.DbCtx, db_gen.ListPageTemplatesByDomainIdAndPagePathIdAscParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				ParentPagePathID:  pgtype.UUID{Bytes: pagePathPathId, Valid: true},
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		}
	} else {
		if sortOrder == "desc" {
			pageTemplates, dbErr = controller.Queries.ListPageTemplatesByDomainIdAndAtRootDesc(controller.DbCtx, db_gen.ListPageTemplatesByDomainIdAndAtRootDescParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		} else {
			pageTemplates, dbErr = controller.Queries.ListPageTemplatesByDomainIdAndAtRootAsc(controller.DbCtx, db_gen.ListPageTemplatesByDomainIdAndAtRootAscParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		}
	}
	if dbErr != nil {
		log.Printf("ListPageTemplates error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve page templates",
		})
		return ctx.Send(payload)
	}

	pageTemplateViews := make([]models.PageTemplate, 0)
	for _, pageTemplate := range pageTemplates {
		pageTemplateView := models.PageTemplate{
			pageTemplate.ID,
			pageTemplate.Created,
			pageTemplate.Modified,
			pageTemplate.Deleted,
			pageTemplate.DomainID,
			pageTemplate.AccountID,
			pageTemplate.ContentID,
			pageTemplate.ParentPagePathID,
			pageTemplate.Title,
			pageTemplate.Slug,
			pageTemplate.Description,
			pageTemplate.IsActive,
		}

		pageTemplateViews = append(pageTemplateViews, pageTemplateView)
	}

	if pageTemplateViews == nil {
		pageTemplateViews = make([]models.PageTemplate, 0)
	}
	results := api_utils.GetPaginatedResults(rowCount, requestedPage, requestedPageSize, pageTemplateViews)
	payload, marshalErr := json.Marshal(results)
	if marshalErr != nil {
		log.Printf("ListPageTemplates error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) ListPagePaths(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	requestedPage, requestedPageSize, sortBy, sortOrder := api_utils.GetPaginationParams(ctx)
	rowCount, dbErr := controller.Queries.CountPagePaths(controller.DbCtx)
	if dbErr != nil {
		log.Printf("CountPagePaths error while retrieving count, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve page path count",
		})
		return ctx.Send(payload)
	}

	query, err := api_utils.ParseQueryString(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid query",
		})
		return ctx.Send(payload)
	}

	domainId, err := uuid.Parse(query.Get("domain"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid query, invalid domain ID",
		})
		return ctx.Send(payload)
	}

	pagePaths := make([]db_gen.PagePath, 0)
	if (query.Get("path") != "null") && (query.Get("path") != "") {
		pagePathPathId, err := uuid.Parse(query.Get("path"))
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "Invalid query, invalid page path ID",
			})
			return ctx.Send(payload)
		}

		if sortOrder == "desc" {
			pagePaths, dbErr = controller.Queries.ListPagePathsByDomainIdAndParentPagePathIdDesc(controller.DbCtx, db_gen.ListPagePathsByDomainIdAndParentPagePathIdDescParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				ParentPagePathID:  pgtype.UUID{Bytes: pagePathPathId, Valid: true},
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		} else {
			pagePaths, dbErr = controller.Queries.ListPagePathsByDomainIdAndParentPagePathIdAsc(controller.DbCtx, db_gen.ListPagePathsByDomainIdAndParentPagePathIdAscParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				ParentPagePathID:  pgtype.UUID{Bytes: pagePathPathId, Valid: true},
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		}
	} else {
		if sortOrder == "desc" {
			pagePaths, dbErr = controller.Queries.ListPagePathsByDomainIdAtRootDesc(controller.DbCtx, db_gen.ListPagePathsByDomainIdAtRootDescParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		} else {
			pagePaths, dbErr = controller.Queries.ListPagePathsByDomainIdAtRootAsc(controller.DbCtx, db_gen.ListPagePathsByDomainIdAtRootAscParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		}
	}
	if dbErr != nil {
		log.Printf("ListPagePaths error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve folders",
		})
		return ctx.Send(payload)
	}

	pages := make([]models.PagePath, 0)
	for _, pagePath := range pagePaths {

		pagePathFolders := make([]db_gen.PagePath, 0)
		if sortOrder == "desc" {
			pagePathFolders, dbErr = controller.Queries.ListPagePathsByDomainIdAndParentPagePathIdDesc(controller.DbCtx, db_gen.ListPagePathsByDomainIdAndParentPagePathIdDescParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				ParentPagePathID:  pagePath.ID,
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		} else {
			pagePathFolders, dbErr = controller.Queries.ListPagePathsByDomainIdAndParentPagePathIdAsc(controller.DbCtx, db_gen.ListPagePathsByDomainIdAndParentPagePathIdAscParams{
				DomainID:          pgtype.UUID{Bytes: domainId, Valid: true},
				ParentPagePathID:  pagePath.ID,
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		}
		if dbErr != nil {
			if dbErr.Error() != "no rows in result set" {
				log.Printf("ListPagePaths error while retrieving object, error %v\n", dbErr)
				ctx.Status(http.StatusBadGateway)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to retrieve child page paths",
				})
				return ctx.Send(payload)
			}
		}

		folders := make([]models.PagePath, 0)
		for _, pagePathFolder := range pagePathFolders {
			folder := models.PagePath{
				ID:               pagePathFolder.ID,
				Created:          pagePathFolder.Created,
				Modified:         pagePathFolder.Modified,
				Deleted:          pagePathFolder.Deleted,
				DomainID:         pagePathFolder.DomainID,
				AccountID:        pagePathFolder.AccountID,
				ParentPagePathID: pagePathFolder.ParentPagePathID,
				Title:            pagePathFolder.Title,
				Slug:             pagePathFolder.Slug,
				IsActive:         pagePathFolder.IsActive,
			}
			folders = append(folders, folder)
		}

		pagePathTemplates, dbErr := controller.Queries.ListAllPageTemplatesByPagePathIdAsc(controller.DbCtx, pagePath.ID)
		if dbErr != nil {
			if dbErr.Error() != "no rows in result set" {
				log.Printf("ListPagePaths error while retrieving object, error %v\n", dbErr)
				ctx.Status(http.StatusBadGateway)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to retrieve child page templates",
				})
				return ctx.Send(payload)
			}
		}

		templates := make([]models.PageTemplate, 0)
		for _, pagePathTemplate := range pagePathTemplates {
			template := models.PageTemplate{
				ID:               pagePathTemplate.ID,
				Created:          pagePathTemplate.Created,
				Modified:         pagePathTemplate.Modified,
				Deleted:          pagePathTemplate.Deleted,
				DomainID:         pagePathTemplate.DomainID,
				AccountID:        pagePathTemplate.AccountID,
				ContentID:        pagePathTemplate.ContentID,
				ParentPagePathID: pagePathTemplate.ParentPagePathID,
				Title:            pagePathTemplate.Title,
				Slug:             pagePathTemplate.Slug,
				IsActive:         pagePathTemplate.IsActive,
			}
			templates = append(templates, template)
		}

		pagePathContent := make([]db_gen.Page, 0)
		if sortOrder == "desc" {
			pagePathContent, dbErr = controller.Queries.ListPagesByPagePathIdDesc(controller.DbCtx, db_gen.ListPagesByPagePathIdDescParams{
				PagePathID:        pagePath.ID,
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		} else {
			pagePathContent, dbErr = controller.Queries.ListPagesByPagePathIdAsc(controller.DbCtx, db_gen.ListPagesByPagePathIdAscParams{
				PagePathID:        pagePath.ID,
				RequestedPage:     int32(requestedPage),
				RequestedPageSize: int32(requestedPageSize),
				SortBy:            sortBy,
			})
		}
		if dbErr != nil {
			if dbErr.Error() != "no rows in result set" {
				log.Printf("ListPagePaths error while retrieving object, error %v\n", dbErr)
				ctx.Status(http.StatusBadGateway)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to retrieve pages",
				})
				return ctx.Send(payload)
			}
		}

		pageContents := make([]models.Page, 0)
		for _, pageContent := range pagePathContent {
			pageContent := models.Page{
				ID:             pageContent.ID,
				Created:        pageContent.Created,
				Modified:       pageContent.Modified,
				Deleted:        pageContent.Deleted,
				DomainID:       pageContent.DomainID,
				AccountID:      pageContent.AccountID,
				ContentID:      pageContent.ContentID,
				PagePathID:     pageContent.PagePathID,
				Title:          pageContent.Title,
				Slug:           pageContent.Slug,
				SeoTitle:       pageContent.SeoTitle,
				SeoDescription: pageContent.SeoDescription,
				SeoKeywords:    pageContent.SeoKeywords,
				PublishAt:      pageContent.PublishAt,
				UnpublishAt:    pageContent.UnpublishAt,
				Version:        pageContent.Version,
				IsActive:       pageContent.IsActive,
			}
			pageContents = append(pageContents, pageContent)
		}

		page := models.PagePath{
			ID:               pagePath.ID,
			Created:          pagePath.Created,
			Modified:         pagePath.Modified,
			Deleted:          pagePath.Deleted,
			DomainID:         pagePath.DomainID,
			AccountID:        pagePath.AccountID,
			ParentPagePathID: pagePath.ParentPagePathID,
			Title:            pagePath.Title,
			Slug:             pagePath.Slug,
			IsActive:         pagePath.IsActive,
			Folders:          folders,
			Pages:            pageContents,
			Templates:        templates,
		}

		pages = append(pages, page)
	}

	if pages == nil {
		pages = make([]models.PagePath, 0)
	}
	results := api_utils.GetPaginatedResults(rowCount, requestedPage, requestedPageSize, pages)
	payload, marshalErr := json.Marshal(results)
	if marshalErr != nil {
		log.Printf("ListPagePaths error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) GetPagePath(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pagePathUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetPagePath error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	pagePath, dbErr := controller.Queries.GetPagePathById(controller.DbCtx, pgtype.UUID{Bytes: pagePathUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page path not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetPagePath error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve page path",
		})
		return ctx.Send(payload)
	}

	pagePathFolders, dbErr := controller.Queries.ListAllPagePathsByParentPagePathIdAsc(controller.DbCtx, pagePath.ID)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("ListPagePaths error while retrieving object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve child page paths",
			})
			return ctx.Send(payload)
		}
	}

	folders := make([]models.PagePath, 0)
	for _, pagePathFolder := range pagePathFolders {
		folder := models.PagePath{
			ID:               pagePathFolder.ID,
			Created:          pagePathFolder.Created,
			Modified:         pagePathFolder.Modified,
			Deleted:          pagePathFolder.Deleted,
			DomainID:         pagePathFolder.DomainID,
			AccountID:        pagePathFolder.AccountID,
			ParentPagePathID: pagePathFolder.ParentPagePathID,
			Title:            pagePathFolder.Title,
			Slug:             pagePathFolder.Slug,
			IsActive:         pagePathFolder.IsActive,
		}
		folders = append(folders, folder)
	}

	pagePathContent, dbErr := controller.Queries.ListAllPagesByPagePathIdAsc(controller.DbCtx, pagePath.ID)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("GetPagePath error while retrieving object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve pages",
			})
			return ctx.Send(payload)
		}
	}

	pageContents := make([]models.Page, 0)
	for _, pageContent := range pagePathContent {
		pageContent := models.Page{
			ID:             pageContent.ID,
			Created:        pageContent.Created,
			Modified:       pageContent.Modified,
			Deleted:        pageContent.Deleted,
			DomainID:       pageContent.DomainID,
			AccountID:      pageContent.AccountID,
			ContentID:      pageContent.ContentID,
			PagePathID:     pageContent.PagePathID,
			Title:          pageContent.Title,
			Slug:           pageContent.Slug,
			SeoTitle:       pageContent.SeoTitle,
			SeoDescription: pageContent.SeoDescription,
			SeoKeywords:    pageContent.SeoKeywords,
			PublishAt:      pageContent.PublishAt,
			UnpublishAt:    pageContent.UnpublishAt,
			Version:        pageContent.Version,
			IsActive:       pageContent.IsActive,
		}
		pageContents = append(pageContents, pageContent)
	}

	pagePathTemplates, dbErr := controller.Queries.ListAllPageTemplatesByPagePathIdAsc(controller.DbCtx, pagePath.ID)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("GetPagePath error while retrieving object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve child page templates",
			})
			return ctx.Send(payload)
		}
	}

	templates := make([]models.PageTemplate, 0)
	for _, pagePathTemplate := range pagePathTemplates {
		template := models.PageTemplate{
			ID:               pagePathTemplate.ID,
			Created:          pagePathTemplate.Created,
			Modified:         pagePathTemplate.Modified,
			Deleted:          pagePathTemplate.Deleted,
			DomainID:         pagePathTemplate.DomainID,
			AccountID:        pagePathTemplate.AccountID,
			ContentID:        pagePathTemplate.ContentID,
			ParentPagePathID: pagePathTemplate.ParentPagePathID,
			Title:            pagePathTemplate.Title,
			Slug:             pagePathTemplate.Slug,
			IsActive:         pagePathTemplate.IsActive,
		}
		templates = append(templates, template)
	}

	pagePathExtensions, dbErr := controller.Queries.ListPagePathExtensionsByPagePathId(controller.DbCtx, pagePath.ID)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("GetPagePath error while retrieving object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve child page extensions",
			})
			return ctx.Send(payload)
		}
	}

	extensions := make([]models.Extension, 0)
	for _, extension := range pagePathExtensions {
		dbExtension, dbErr := controller.Queries.GetExtensionById(controller.DbCtx, extension.ExtensionID)
		if dbErr != nil {
			if dbErr.Error() != "no rows in result set" {
				log.Printf("GetPagePath error while retrieving object, error %v\n", dbErr)
				ctx.Status(http.StatusBadGateway)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to retrieve child page extensions",
				})
				return ctx.Send(payload)
			}

			continue
		}

		extension := models.Extension{
			ID:       dbExtension.ID,
			Name:     dbExtension.Name,
			Slug:     dbExtension.Slug,
			Icon:     dbExtension.Icon,
			Data:     dbExtension.Data,
			IsActive: dbExtension.IsActive,
		}
		extensions = append(extensions, extension)
	}

	page := models.PagePath{
		ID:               pagePath.ID,
		Created:          pagePath.Created,
		Modified:         pagePath.Modified,
		Deleted:          pagePath.Deleted,
		DomainID:         pagePath.DomainID,
		AccountID:        pagePath.AccountID,
		ParentPagePathID: pagePath.ParentPagePathID,
		Title:            pagePath.Title,
		Slug:             pagePath.Slug,
		IsActive:         pagePath.IsActive,
		Folders:          folders,
		Pages:            pageContents,
		Templates:        templates,
		Extensions:       extensions,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   page,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetPagePath error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) DeletePagePath(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pagePathUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("DeletePagePath error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	_, dbErr := controller.Queries.GetPagePathById(controller.DbCtx, pgtype.UUID{Bytes: pagePathUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page path not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("DeletePagePath error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete page path",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeletePagePath(controller.DbCtx, pgtype.UUID{Bytes: pagePathUuid, Valid: true})
	if dbErr != nil {
		log.Printf("DeletePagePath error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete page path",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("DeletePagePath error serialising response, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) UpdatePagePath(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pagePathUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("UpdatePagePath error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	sourcePagePath, dbErr := controller.Queries.GetPagePathById(controller.DbCtx, pgtype.UUID{Bytes: pagePathUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page path not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdatePagePath error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to update page path",
		})
		return ctx.Send(payload)
	}

	pagePath := db_gen.UpdatePagePathParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &pagePath)
	if unmarshalErr != nil {
		log.Printf("UpdatePagePath error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	// Values that are not allowed to be changed when updating an pagePath
	pagePath.ID = sourcePagePath.ID

	updatedPagePath, dbErr := controller.Queries.UpdatePagePath(controller.DbCtx, pagePath)
	if dbErr != nil {
		log.Printf("UpdatePagePath error while updating object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to perform update",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   updatedPagePath,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("UpdatePagePath error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) ListLinks(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	query, err := api_utils.ParseQueryString(ctx)
	if err != nil {
		log.Printf("ListLinks error while parsing query, error %v\n", err)
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid query",
		})
		return ctx.Send(payload)
	}

	domainId, err := uuid.Parse(query.Get("domain"))
	if err != nil {
		log.Printf("ListLinks error while parsing query, error %v\n", err)
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid query, invalid domain ID",
		})
		return ctx.Send(payload)
	}

	pageLinks, err := controller.ListChildLinks([]string{}, pgtype.UUID{Bytes: domainId, Valid: true}, nil)
	if err != nil {
		log.Printf("ListLinks error while retrieving object, error %v\n", err)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve links",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.ResultsResponse,
		Data:   pageLinks,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("ListLinks error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) AddBlog(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	blog := models.BlogParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &blog)
	if unmarshalErr != nil {
		log.Printf("AddBlog error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	_, dbErr := controller.Queries.GetBlogByPageId(controller.DbCtx, blog.PageID)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("AddBlog error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusConflict)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusConflict,
				ErrorMessage: "Unable to add blog blog",
			})
			return ctx.Send(payload)
		}
	}

	dbBlog := db_gen.CreateBlogParams{
		DomainID:    blog.DomainID,
		AccountID:   blog.AccountID,
		PageID:      blog.PageID,
		Title:       blog.Title,
		Description: blog.Description,
		Image:       blog.Image,
		ImageInfo:   blog.ImageInfo,
		Keywords:    blog.Keywords,
		IsActive:    blog.IsActive,
	}

	blogResult, dbErr := controller.Queries.CreateBlog(controller.DbCtx, dbBlog)
	if dbErr != nil {
		log.Printf("AddBlog error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to add blog",
		})
		return ctx.Send(payload)
	}

	pageLink, err := controller.GetPageLink(blog.PageID)
	if err != nil {
		log.Printf("AddBlog error while fetching object, error %v\n", err)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: err.Error(),
		})
		return ctx.Send(payload)
	}

	blogView := models.BlogView{
		ID:          blogResult.ID,
		DomainID:    blogResult.DomainID,
		AccountID:   blogResult.AccountID,
		PageID:      blogResult.PageID,
		Link:        pageLink,
		Title:       blogResult.Title,
		Description: blogResult.Description,
		Image:       blogResult.Image,
		ImageInfo:   blogResult.ImageInfo,
		Keywords:    blogResult.Keywords,
		IsActive:    blogResult.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   blogView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: marshalErr.Error(),
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) UpdateBlog(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	blogUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("UpdateBlog error parsing blog ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid blog ID",
		})
		return ctx.Send(payload)
	}

	sourceBlog, dbErr := controller.Queries.GetBlogById(controller.DbCtx, pgtype.UUID{Bytes: blogUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Blog not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdateBlog error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to update blog",
		})
		return ctx.Send(payload)
	}

	blog := db_gen.UpdateBlogByIdParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &blog)
	if unmarshalErr != nil {
		log.Printf("UpdateBlog error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	pageLink, err := controller.GetPageLink(sourceBlog.PageID)
	if err != nil {
		log.Printf("UpdateBlog error while fetching object, error %v\n", err)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: err.Error(),
		})
		return ctx.Send(payload)
	}

	// Values that are not allowed to be changed when updating a blog
	blog.ID = sourceBlog.ID
	blog.PageID = sourceBlog.PageID
	blog.DomainID = sourceBlog.DomainID
	blog.AccountID = sourceBlog.AccountID

	updatedBlog, dbErr := controller.Queries.UpdateBlogById(controller.DbCtx, blog)
	if dbErr != nil {
		log.Printf("UpdateBlog error while updating object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to perform update",
		})
		return ctx.Send(payload)
	}

	blogView := models.BlogView{
		updatedBlog.ID,
		updatedBlog.DomainID,
		updatedBlog.AccountID,
		updatedBlog.PageID,
		pageLink,
		updatedBlog.Title,
		updatedBlog.Description,
		updatedBlog.Image,
		updatedBlog.ImageInfo,
		updatedBlog.Keywords,
		updatedBlog.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   blogView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("UpdateBlog error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) RemoveBlog(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	blogUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("RemoveBlog error parsing blog ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid blog ID",
		})
		return ctx.Send(payload)
	}

	blog, dbErr := controller.Queries.GetBlogById(controller.DbCtx, pgtype.UUID{Bytes: blogUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Blog not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("RemoveBlog error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove blog",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeleteBlogById(controller.DbCtx, blog.ID)
	if dbErr != nil {
		log.Printf("RemoveBlog error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove blog",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("RemoveBlog error serialising response, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) GetBlog(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	blogUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetBlog error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}

	blog, dbErr := controller.Queries.GetBlogById(controller.DbCtx, pgtype.UUID{Bytes: blogUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetPage error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve page",
		})
		return ctx.Send(payload)
	}

	pageLink, err := controller.GetPageLink(blog.PageID)
	if err != nil {
		log.Printf("UpdateBlog error while fetching object, error %v\n", err)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: err.Error(),
		})
		return ctx.Send(payload)
	}

	blogView := models.BlogView{
		ID:          blog.ID,
		DomainID:    blog.DomainID,
		AccountID:   blog.AccountID,
		PageID:      blog.PageID,
		Link:        pageLink,
		Title:       blog.Title,
		Description: blog.Description,
		Image:       blog.Image,
		ImageInfo:   blog.ImageInfo,
		Keywords:    blog.Keywords,
		IsActive:    blog.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   blogView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetBlog error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) GetBlogByPage(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	pageUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetBlog error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}

	blog, dbErr := controller.Queries.GetBlogByPageId(controller.DbCtx, pgtype.UUID{Bytes: pageUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Page not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetPage error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve page",
		})
		return ctx.Send(payload)
	}

	pageLink, err := controller.GetPageLink(blog.PageID)
	if err != nil {
		log.Printf("UpdateBlog error while fetching object, error %v\n", err)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: err.Error(),
		})
		return ctx.Send(payload)
	}

	blogView := models.BlogView{
		ID:          blog.ID,
		DomainID:    blog.DomainID,
		AccountID:   blog.AccountID,
		PageID:      blog.PageID,
		Link:        pageLink,
		Title:       blog.Title,
		Description: blog.Description,
		Image:       blog.Image,
		ImageInfo:   blog.ImageInfo,
		Keywords:    blog.Keywords,
		IsActive:    blog.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   blogView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetBlog error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *PagePathModel) ListBlogs(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	blogs, dbErr := controller.Queries.ListAllBlogs(controller.DbCtx)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			blogs = make([]db_gen.Blog, 0)
		} else {
			log.Printf("ListBlogs error while retrieving object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve blogs",
			})
			return ctx.Send(payload)
		}
	}

	blogViews := make([]models.BlogItemView, 0)
	for _, blog := range blogs {
		pageLink, err := controller.GetPageLink(blog.PageID)
		if err != nil {
			log.Printf("ListBlogs error while fetching object, error %v\n", err)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: err.Error(),
			})
			return ctx.Send(payload)
		}

		blogViews = append(blogViews, models.BlogItemView{
			ID:          blog.ID,
			Created:     blog.Created,
			Link:        pageLink,
			Title:       blog.Title,
			Description: blog.Description,
			Image:       blog.Image,
			ImageInfo:   blog.ImageInfo,
			Keywords:    blog.Keywords,
		})
	}

	success := api_models.Success{
		Status: api_utils.ResultsResponse,
		Data:   blogViews,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetBlog error while marshalling object, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}
