package v1

import (
	"context"
	"cortex_api/common/api_models"
	"cortex_api/common/api_utils"
	"cortex_api/database/db_gen"
	"cortex_api/services/cms_service/models"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type ContentModel struct {
	Queries *db_gen.Queries
	DbCtx   context.Context
}

type ContentInterface interface {
	CreateContent(ctx *fiber.Ctx) error
	ListContents(ctx *fiber.Ctx) error
	GetContent(ctx *fiber.Ctx) error
	DeleteContent(ctx *fiber.Ctx) error
	UpdateContent(ctx *fiber.Ctx) error
	AddChildContent(ctx *fiber.Ctx) error
	UpdateChildContent(ctx *fiber.Ctx) error
	RemoveChildContent(ctx *fiber.Ctx) error
}

func ContentController(
	queries *db_gen.Queries,
	dbCtx context.Context,
) *ContentModel {
	return &ContentModel{
		Queries: queries,
		DbCtx:   dbCtx,
	}
}

func (controller *ContentModel) ClonePageContentIteratively(sourceContentUuid pgtype.UUID, ctx *fiber.Ctx, newTitle *string, newSlug *string) (*db_gen.Content, error) {
	// get the source content
	sourceContent, dbErr := controller.Queries.GetContentById(controller.DbCtx, sourceContentUuid)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			return nil, nil
		}
		log.Printf("CloneContent error while fetching object, error %v\n", dbErr)
		return nil, errors.New("content not found")
	}

	if newTitle != nil {
		sourceContent.Title = *newTitle
	}
	if newSlug != nil {
		sourceContent.Slug = *newSlug
	}

	// create a new content object with the same data as the source content
	newContent := db_gen.CreateContentParams{
		DomainID:    sourceContent.DomainID,
		AccountID:   sourceContent.AccountID,
		Title:       sourceContent.Title,
		Slug:        sourceContent.Slug,
		Data:        sourceContent.Data,
		IsActive:    sourceContent.IsActive,
		ComponentID: sourceContent.ComponentID,
	}

	// insert the new content object
	destContent, dbErr := controller.Queries.CreateContent(controller.DbCtx, newContent)
	if dbErr != nil {
		log.Printf("CloneContent error while inserting object, error %v\n", dbErr)
		return nil, errors.New("unable to clone content")
	}

	// get any child content blocks of the source content
	children, err := controller.Queries.ListContentCollectionByParentId(controller.DbCtx, sourceContent.ID)
	if err != nil {
		log.Printf("CloneContent error while fetching child information, error %v\n", err)
		return nil, errors.New("unable to retrieve child content information")
	}

	// iterate through children and clone them
	clonedChildren := make([]db_gen.Content, 0)
	for _, child := range children {
		clonedChild, dbErr := controller.ClonePageContentIteratively(child.ContentID, ctx, nil, nil)
		if dbErr != nil {
			log.Printf("CloneContent error while cloning child, error %v\n", dbErr)
			return nil, errors.New("unable to clone content")
		}

		if clonedChild != nil {
			clonedChildren = append(clonedChildren, *clonedChild)

			// add the cloned child to the new content
			contentCollectionEntry := db_gen.CreateContentCollectionAndReturnIdParams{
				ParentID:  destContent.ID,
				ContentID: pgtype.UUID{Bytes: clonedChild.ID.Bytes, Valid: true},
				Ordering:  child.Ordering,
				IsActive:  child.IsActive,
			}
			_, dbErr = controller.Queries.CreateContentCollectionAndReturnId(controller.DbCtx, contentCollectionEntry)
			if dbErr != nil {
				log.Printf("CloneContent error while inserting object, error %v\n", dbErr)
				return nil, errors.New("unable to clone child content")
			}
		}
	}

	return &destContent, nil
}

func (controller *ContentModel) ClonePageContent(sourceContentUuid pgtype.UUID, ctx *fiber.Ctx, title string, slug string) (db_gen.Content, error) {
	content, err := controller.ClonePageContentIteratively(sourceContentUuid, ctx, &title, &slug)
	if err != nil {
		return db_gen.Content{}, err
	}
	return *content, nil
}

func (controller *ContentModel) CloneContent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	content := db_gen.CreateContentParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &content)
	if unmarshalErr != nil {
		log.Printf("CreateContent error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	// get the source content ID
	sourceContentUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("CloneContent error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}

	// clone the content iteratively
	destContent, err := controller.ClonePageContent(pgtype.UUID{Bytes: sourceContentUuid, Valid: true}, ctx, content.Title, content.Slug)
	if err != nil {
		log.Printf("CloneContent error while cloning content, error %v\n", err)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to clone content",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   destContent,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("CloneContent error while marshalling object, error %v\n", marshalErr)
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

func (controller *ContentModel) CreateContent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	contentParams := db_gen.CreateContentParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &contentParams)
	if unmarshalErr != nil {
		log.Printf("CreateContent error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}
	log.Printf("CreateContent contentParams %v\n", contentParams.Slug)

	_, dbErr := controller.Queries.GetContentByDomainIdAndSlug(controller.DbCtx, db_gen.GetContentByDomainIdAndSlugParams{
		DomainID: contentParams.DomainID,
		Slug:     contentParams.Slug,
	})
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("CreateContent error while fetching object, error content already exists\n")
			ctx.Status(http.StatusConflict)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusConflict,
				ErrorMessage: "Content already exists",
			})
			return ctx.Send(payload)
		}
	}

	newContent, dbErr := controller.Queries.CreateContent(controller.DbCtx, contentParams)
	if dbErr != nil {
		log.Printf("CreateContent error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to create contentParams",
		})
		return ctx.Send(payload)
	}

	component, err := controller.Queries.GetComponentById(controller.DbCtx, newContent.ComponentID)
	if err != nil {
		log.Printf("CreateContent error while fetching component, error %v\n", err)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve component information",
		})
		return ctx.Send(payload)
	}

	componentView := models.ComponentView{
		ID:                  component.ID,
		Name:                component.Name,
		Icon:                component.Icon,
		Description:         component.Description,
		ClassName:           component.ClassName,
		HtmlTag:             component.HtmlTag,
		ChildTagConstraints: component.ChildTagConstraints,
		IsActive:            component.IsActive,
	}

	contentView := models.ContentView{
		ID:        newContent.ID,
		DomainId:  newContent.DomainID,
		Component: componentView,
		AccountId: newContent.AccountID,
		Title:     newContent.Title,
		Slug:      newContent.Slug,
		Data:      newContent.Data,
		IsActive:  newContent.IsActive,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   contentView,
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

func (controller *ContentModel) AddChildContent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	contentCollection := models.CreateContentCollectionParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &contentCollection)
	if unmarshalErr != nil {
		log.Printf("AddChildContent error while unmarshalling object, error %v\n", unmarshalErr.Error())
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
		log.Printf("AddChildContent error parsing parent content ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid parent content block ID",
		})
		return ctx.Send(payload)
	}

	_, dbErr := controller.Queries.GetContentCollectionByParentIdAndContentId(controller.DbCtx, db_gen.GetContentCollectionByParentIdAndContentIdParams{
		ParentID:  pgtype.UUID{Bytes: parentUuid, Valid: true},
		ContentID: contentCollection.ContentID,
	})
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("AddChildContent error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to add child content block",
			})
			return ctx.Send(payload)
		}
	}
	if dbErr == nil {
		log.Printf("AddChildContent error while fetching object, error content child already exists for this parent\n")
		ctx.Status(http.StatusConflict)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusConflict,
			ErrorMessage: "This child content block already exists as part of this parent block",
		})
		return ctx.Send(payload)
	}

	dbContentCollection := db_gen.CreateContentCollectionAndReturnIdParams{
		ParentID:  pgtype.UUID{Bytes: parentUuid, Valid: true},
		ContentID: contentCollection.ContentID,
		Ordering:  contentCollection.Ordering,
		IsActive:  contentCollection.IsActive,
	}

	contentCollectionUuid, dbErr := controller.Queries.CreateContentCollectionAndReturnId(controller.DbCtx, dbContentCollection)
	if dbErr != nil {
		log.Printf("AddChildContent error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to add child content block",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   contentCollectionUuid,
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

func (controller *ContentModel) UpdateChildContent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	updateContentCollection := models.UpdateContentCollectionParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &updateContentCollection)
	if unmarshalErr != nil {
		log.Printf("UpdateChildContent error while unmarshalling object, error %v\n", unmarshalErr.Error())
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
		log.Printf("UpdateChildContent error parsing parent content ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid parent content block ID",
		})
		return ctx.Send(payload)
	}

	contentUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("UpdateChildContent error parsing content ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid content block ID",
		})
		return ctx.Send(payload)
	}

	originalContentCollection, dbErr := controller.Queries.GetContentCollectionByParentIdAndContentId(controller.DbCtx, db_gen.GetContentCollectionByParentIdAndContentIdParams{
		ParentID:  pgtype.UUID{Bytes: parentUuid, Valid: true},
		ContentID: pgtype.UUID{Bytes: contentUuid, Valid: true},
	})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			log.Printf("UpdateChildContent error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to update child content block",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdateChildContent error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusNotFound)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Child content block not found",
		})
		return ctx.Send(payload)
	}

	updatedContentCollection := db_gen.UpdateContentCollectionByIdParams{
		ID:        originalContentCollection.ID,
		ParentID:  originalContentCollection.ParentID,
		ContentID: originalContentCollection.ContentID,
		Ordering:  updateContentCollection.Ordering,
		IsActive:  updateContentCollection.IsActive,
	}

	newContentCollection, dbErr := controller.Queries.UpdateContentCollectionById(controller.DbCtx, updatedContentCollection)
	if dbErr != nil {
		log.Printf("UpdateChildContent error while updating object, error %v\n", dbErr)
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
		Data:   newContentCollection,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("UpdateChildContent error while marshalling object, error %v\n", marshalErr)
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

func (controller *ContentModel) RemoveChildContent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	contentCollectionUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("RemoveChildContent error parsing parent content ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid child content block ID",
		})
		return ctx.Send(payload)
	}

	contentCollection, dbErr := controller.Queries.GetContentCollectionById(controller.DbCtx, pgtype.UUID{Bytes: contentCollectionUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Child content block not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("RemoveChildContent error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove child content block",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeleteContentCollection(controller.DbCtx, contentCollection.ID)
	if dbErr != nil {
		log.Printf("RemoveChildContent error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove child content block",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("RemoveChildContent error serialising response, error %v\n", marshalErr)
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

func (controller *ContentModel) GetContentById(ctx *fiber.Ctx, contentId pgtype.UUID) (*models.Content, error) {
	// TODO - add max depth to prevent infinite recursion
	content, dbErr := controller.Queries.GetContentById(controller.DbCtx, contentId)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Content not found",
			})
			return nil, ctx.Send(payload)
		}
		log.Printf("GetContentById error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve content",
		})
		return nil, ctx.Send(payload)
	}

	component, dbErr := controller.Queries.GetComponentById(controller.DbCtx, content.ComponentID)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Component not found",
			})
			return nil, ctx.Send(payload)
		}
		log.Printf("GetContentById error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve component",
		})
		return nil, ctx.Send(payload)
	}

	fields, err := controller.Queries.ListComponentFieldByComponentId(controller.DbCtx, component.ID)
	if err != nil {
		log.Printf("GetContentById error while fetching field information, error %v\n", err)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve field information",
		})
		return nil, ctx.Send(payload)
	}

	fieldViews := make([]models.ContentComponentField, 0)
	for _, field := range fields {
		fieldViews = append(fieldViews, models.ContentComponentField{
			ID:           field.ID,
			Name:         field.Name,
			Description:  field.Description,
			DataType:     field.DataType,
			EditorType:   field.EditorType,
			Validation:   field.Validation,
			DefaultValue: field.DefaultValue,
			IsActive:     field.IsActive,
		})
	}

	contentComponent := models.ContentComponent{
		ID:                  component.ID,
		Name:                component.Name,
		Icon:                component.Icon,
		Description:         component.Description,
		ClassName:           component.ClassName,
		HtmlTag:             component.HtmlTag,
		ChildTagConstraints: component.ChildTagConstraints,
		IsActive:            component.IsActive,
		Fields:              fieldViews,
	}

	children, err := controller.Queries.ListContentCollectionByParentId(controller.DbCtx, content.ID)
	if err != nil {
		log.Printf("GetContentById error while fetching child information, error %v\n", err)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve child content information",
		})
		return nil, ctx.Send(payload)
	}

	childContentList := make([]*models.Content, 0)
	for _, child := range children {
		childContent, err := controller.GetContentById(ctx, child.ContentID)
		if err != nil {
			log.Printf("GetContentById error while fetching child content, error %v\n", err)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve child content",
			})
			return nil, ctx.Send(payload)
		}
		if childContent != nil {
			childContentList = append(childContentList, childContent)
		}
	}

	return &models.Content{
		ID:        content.ID,
		DomainID:  content.DomainID,
		AccountID: content.AccountID,
		Title:     content.Title,
		Slug:      content.Slug,
		Data:      content.Data,
		IsActive:  content.IsActive,
		Component: contentComponent,
		Children:  childContentList,
	}, nil
}

func (controller *ContentModel) ListContents(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	requestedPage, requestedPageSize, sortBy, sortOrder := api_utils.GetPaginationParams(ctx)
	rowCount, dbErr := controller.Queries.CountContents(controller.DbCtx)
	if dbErr != nil {
		log.Printf("CountContents error while retrieving count, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve content count",
		})
		return ctx.Send(payload)
	}

	contents := make([]db_gen.Content, 0)
	if sortOrder == "desc" {
		contents, dbErr = controller.Queries.ListContentsDesc(controller.DbCtx, db_gen.ListContentsDescParams{
			RequestedPage:     int32(requestedPage),
			RequestedPageSize: int32(requestedPageSize),
			SortBy:            sortBy,
		})
	} else {
		contents, dbErr = controller.Queries.ListContentsAsc(controller.DbCtx, db_gen.ListContentsAscParams{
			RequestedPage:     int32(requestedPage),
			RequestedPageSize: int32(requestedPageSize),
			SortBy:            sortBy,
		})
	}
	if dbErr != nil {
		log.Printf("ListContents error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve contents",
		})
		return ctx.Send(payload)
	}

	if contents == nil {
		contents = make([]db_gen.Content, 0)
	}

	contentsView := make([]*models.Content, 0)
	for _, content := range contents {
		contentView, err := controller.GetContentById(ctx, content.ID)
		if err != nil {
			return err
		}

		contentsView = append(contentsView, contentView)
	}

	results := api_utils.GetPaginatedResults(rowCount, requestedPage, requestedPageSize, contentsView)
	payload, marshalErr := json.Marshal(results)
	if marshalErr != nil {
		log.Printf("ListContents error while marshalling object, error %v\n", marshalErr)
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

func (controller *ContentModel) GetContent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	contentUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetContent error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}

	contentView, err := controller.GetContentById(ctx, pgtype.UUID{Bytes: contentUuid, Valid: true})
	if err != nil {
		return err
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   contentView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetContent error while marshalling object, error %v\n", marshalErr)
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

func (controller *ContentModel) DeleteContent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	contentUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("DeleteContent error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	_, dbErr := controller.Queries.GetContentById(controller.DbCtx, pgtype.UUID{Bytes: contentUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Content not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("DeleteContent error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete content",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeleteContent(controller.DbCtx, pgtype.UUID{Bytes: contentUuid, Valid: true})
	if dbErr != nil {
		log.Printf("DeleteContent error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete content",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("DeleteContent error serialising response, error %v\n", marshalErr)
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

func (controller *ContentModel) UpdateContent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	contentUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("UpdateContent error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	sourceContent, dbErr := controller.Queries.GetContentById(controller.DbCtx, pgtype.UUID{Bytes: contentUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Content not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdateContent error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to update content",
		})
		return ctx.Send(payload)
	}

	content := db_gen.UpdateContentParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &content)
	if unmarshalErr != nil {
		log.Printf("UpdateContent error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	// Values that are not allowed to be changed when updating an content
	content.ID = sourceContent.ID

	updatedContent, dbErr := controller.Queries.UpdateContent(controller.DbCtx, content)
	if dbErr != nil {
		log.Printf("UpdateContent error while updating object, error %v\n", dbErr)
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
		Data:   updatedContent,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("UpdateContent error while marshalling object, error %v\n", marshalErr)
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

func (controller *ContentModel) GetPageContentByFilename(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	domainName := ctx.Get("X-Origin-Domain")
	if domainName == "" {
		log.Printf("GetPageContentByFilename error missing domain header\n")
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Missing domain header",
		})
		return ctx.Send(payload)
	}

	domain, dbErr := controller.Queries.GetDomainByDomainName(controller.DbCtx, domainName)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Domain not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetPageContentByFilename error while fetching domain, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve domain",
		})
		return ctx.Send(payload)
	}

	filePath, filename := filepath.Split(ctx.Params("*"))
	if filename == "" {
		filename = "index"
	}

	fullPath := strings.Split(filePath, "/")

	// TODO - implement path walking here
	// walk through filePath from split above and get the parent folder

	parentFolder := db_gen.PagePath{}
	if (len(fullPath) == 1 && fullPath[0] == "") || (len(fullPath) == 1 && fullPath[0] == "/") {
		parentFolder, dbErr = controller.Queries.GetPagePathByDomainIdRootPathAndSlug(controller.DbCtx, db_gen.GetPagePathByDomainIdRootPathAndSlugParams{
			DomainID: domain.ID,
			Slug:     "index",
		})
		if dbErr != nil {
			if dbErr.Error() != "no rows in result set" {
				log.Printf("GetPageContentByFilename error while fetching object, error %v\n", dbErr)
				ctx.Status(http.StatusBadGateway)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to create page path",
				})
				return ctx.Send(payload)
			}
			log.Printf("GetPageContentByFilename error while fetching parent folder, error %v\n", dbErr)
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Folder not found",
			})
			return ctx.Send(payload)
		}
	} else {
		parentFolder, dbErr = controller.Queries.GetPagePathByDomainIdAndSlug(controller.DbCtx, db_gen.GetPagePathByDomainIdAndSlugParams{
			DomainID: domain.ID,
			Slug:     fullPath[len(fullPath)-2],
		})
		if dbErr != nil {
			if dbErr.Error() != "no rows in result set" {
				log.Printf("GetPageContentByFilename error while fetching object, error %v\n", dbErr)
				ctx.Status(http.StatusBadGateway)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusBadGateway,
					ErrorMessage: "Unable to create page path",
				})
				return ctx.Send(payload)
			}
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Folder not found",
			})
			return ctx.Send(payload)
		}
	}

	page, dbErr := controller.Queries.GetPageByDomainIdPagePathIdAndSlug(controller.DbCtx, db_gen.GetPageByDomainIdPagePathIdAndSlugParams{
		DomainID:   domain.ID,
		PagePathID: parentFolder.ID,
		Slug:       filename,
	})
	if dbErr != nil {
		log.Printf("GetPageContentByFilename error while fetching page, error %v\n", dbErr)
		ctx.Status(http.StatusNotFound)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Page not found",
		})
		return ctx.Send(payload)
	}

	if !page.ContentID.Valid {
		ctx.Status(http.StatusNotFound)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Content not found",
		})
		return ctx.Send(payload)
	}

	breadcrumbs := make([]models.BreadCrumb, 0)
	if parentFolder.ID.Valid {
		breadcrumbs = append(breadcrumbs, models.BreadCrumb{
			ID:       parentFolder.ID,
			ParentID: parentFolder.ParentPagePathID,
			Title:    parentFolder.Title,
			Slug:     parentFolder.Slug,
		})

		currentPagePath := parentFolder
		for currentPagePath.ParentPagePathID.Valid {
			currentPagePath, dbErr = controller.Queries.GetPagePathById(controller.DbCtx, currentPagePath.ParentPagePathID)
			if dbErr != nil {
				log.Printf("GetPageContentByFilename error while fetching parent folder, error %v\n", dbErr)
				ctx.Status(http.StatusNotFound)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusNotFound,
					ErrorMessage: "Folder not found",
				})
				return ctx.Send(payload)
			}

			breadcrumbs = append(breadcrumbs, models.BreadCrumb{
				ID:       currentPagePath.ID,
				ParentID: currentPagePath.ParentPagePathID,
				Title:    currentPagePath.Title,
				Slug:     currentPagePath.Slug,
			})
		}

		slices.Reverse(breadcrumbs)
	}

	renderedContent, err := controller.GetContentById(ctx, page.ContentID)
	if err != nil {
		log.Printf("GetPageContentByFilename error while getting page, error %v\n", err)
		ctx.Status(http.StatusNotFound)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Page not found",
		})
		return ctx.Send(payload)
	}

	pageContent := models.PageContent{
		Page:        page,
		Content:     renderedContent,
		BreadCrumbs: breadcrumbs,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   pageContent,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetPageContentByFilename error serialising response, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	ctx.Set("Content-Type", "application/json")
	return ctx.Send(payload)
}

func (controller *ContentModel) GetPageContentById(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	domainName := ctx.Get("X-Origin-Domain")
	if domainName == "" {
		log.Printf("GetPageContentById error missing domain header\n")
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Missing domain header",
		})
		return ctx.Send(payload)
	}

	domain, dbErr := controller.Queries.GetDomainByDomainName(controller.DbCtx, domainName)
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Domain not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetPageContentById error while fetching domain, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve domain",
		})
		return ctx.Send(payload)
	}

	pageUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetPageContentById error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}

	page, dbErr := controller.Queries.GetPageByDomainIdAndId(controller.DbCtx, db_gen.GetPageByDomainIdAndIdParams{
		DomainID: domain.ID,
		ID:       pgtype.UUID{Bytes: pageUuid, Valid: true},
	})
	if dbErr != nil {
		log.Printf("GetPageContentById error while fetching page, error %v\n", dbErr)
		ctx.Status(http.StatusNotFound)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Page not found",
		})
		return ctx.Send(payload)
	}

	if !page.ContentID.Valid {
		ctx.Status(http.StatusNotFound)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Content not found",
		})
		return ctx.Send(payload)
	}

	if !page.DraftPageID.Valid {
		ctx.Status(http.StatusNotFound)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Content not found",
		})
		return ctx.Send(payload)
	}

	breadcrumbs := make([]models.BreadCrumb, 0)
	if page.PagePathID.Valid {
		pagePathParent, dbErr := controller.Queries.GetPagePathById(controller.DbCtx, page.PagePathID)
		if dbErr != nil {
			log.Printf("GetPageContentById error while fetching parent folder, error %v\n", dbErr)
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Folder not found",
			})
			return ctx.Send(payload)
		}

		breadcrumbs = append(breadcrumbs, models.BreadCrumb{
			ID:       pagePathParent.ID,
			ParentID: pagePathParent.ParentPagePathID,
			Title:    pagePathParent.Title,
			Slug:     pagePathParent.Slug,
		})

		currentPagePath := pagePathParent
		for currentPagePath.ParentPagePathID.Valid {
			currentPagePath, dbErr = controller.Queries.GetPagePathById(controller.DbCtx, currentPagePath.ParentPagePathID)
			if dbErr != nil {
				log.Printf("GetPageContentById error while fetching parent folder, error %v\n", dbErr)
				ctx.Status(http.StatusNotFound)
				payload, _ := json.Marshal(api_models.Failed{
					Status:       api_utils.FailedResponse,
					ErrorCode:    http.StatusNotFound,
					ErrorMessage: "Folder not found",
				})
				return ctx.Send(payload)
			}

			breadcrumbs = append(breadcrumbs, models.BreadCrumb{
				ID:       currentPagePath.ID,
				ParentID: currentPagePath.ParentPagePathID,
				Title:    currentPagePath.Title,
				Slug:     currentPagePath.Slug,
			})
		}

		slices.Reverse(breadcrumbs)
	}

	renderedContent, err := controller.GetContentById(ctx, page.ContentID)
	if err != nil {
		log.Printf("GetPageContentById error while getting page, error %v\n", err)
		ctx.Status(http.StatusNotFound)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusNotFound,
			ErrorMessage: "Page not found",
		})
		return ctx.Send(payload)
	}

	pageContent := models.PageContent{
		Page:        page,
		Content:     renderedContent,
		BreadCrumbs: breadcrumbs,
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   pageContent,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetPageContentById error serialising response, error %v\n", marshalErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to serialise response",
		})
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	ctx.Set("Content-Type", "application/json")
	return ctx.Send(payload)
}
