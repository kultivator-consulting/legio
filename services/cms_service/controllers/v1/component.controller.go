package v1

import (
	"context"
	"cortex_api/common/api_models"
	"cortex_api/common/api_utils"
	"cortex_api/database/db_gen"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
)

type ComponentField struct {
	ID           pgtype.UUID `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	DataType     string      `json:"dataType"`
	EditorType   string      `json:"editorType"`
	Validation   string      `json:"validation"`
	DefaultValue string      `json:"defaultValue"`
	IsActive     bool        `json:"isActive"`
}

type Component struct {
	ID                  pgtype.UUID      `json:"id"`
	Name                string           `json:"name"`
	Icon                string           `json:"icon"`
	Description         string           `json:"description"`
	ClassName           string           `json:"className"`
	HtmlTag             string           `json:"htmlTag"`
	ChildTagConstraints []string         `json:"childTagConstraints"`
	IsActive            bool             `json:"isActive"`
	Fields              []ComponentField `json:"fields"`
}

type CreateComponentFieldParams struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DataType     string `json:"dataType"`
	EditorType   string `json:"editorType"`
	Validation   string `json:"validation"`
	DefaultValue string `json:"defaultValue"`
	IsActive     bool   `json:"isActive"`
}

type ComponentModel struct {
	Queries *db_gen.Queries
	DbCtx   context.Context
}

type ComponentInterface interface {
	CreateComponent(ctx *fiber.Ctx) error
	ListComponents(ctx *fiber.Ctx) error
	GetComponent(ctx *fiber.Ctx) error
	DeleteComponent(ctx *fiber.Ctx) error
	UpdateComponent(ctx *fiber.Ctx) error
	AddComponentField(ctx *fiber.Ctx) error
	RemoveComponentField(ctx *fiber.Ctx) error
}

func ComponentController(
	queries *db_gen.Queries,
	dbCtx context.Context,
) *ComponentModel {
	return &ComponentModel{
		Queries: queries,
		DbCtx:   dbCtx,
	}
}

func (controller *ComponentModel) CreateComponent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	component := db_gen.CreateComponentAndReturnIdParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &component)
	if unmarshalErr != nil {
		log.Printf("CreateComponent error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	_, dbErr := controller.Queries.GetComponentByName(controller.DbCtx, component.Name)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("CreateComponent error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to create component",
			})
			return ctx.Send(payload)
		}
	}
	if dbErr == nil {
		log.Printf("CreateComponent error while fetching object, error component already exists\n")
		ctx.Status(http.StatusConflict)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusConflict,
			ErrorMessage: "Component already exists",
		})
		return ctx.Send(payload)
	}

	componentUuid, dbErr := controller.Queries.CreateComponentAndReturnId(controller.DbCtx, component)
	if dbErr != nil {
		log.Printf("CreateComponent error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to create component",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   componentUuid,
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

func (controller *ComponentModel) AddComponentField(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	componentFieldParams := CreateComponentFieldParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &componentFieldParams)
	if unmarshalErr != nil {
		log.Printf("AddComponentField error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	componentUuid, err := uuid.Parse(ctx.Params("componentId"))
	if err != nil {
		log.Printf("AddComponentField error parsing component ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid component ID",
		})
		return ctx.Send(payload)
	}

	componentFields, dbErr := controller.Queries.ListComponentFieldByComponentIdAndName(controller.DbCtx, db_gen.ListComponentFieldByComponentIdAndNameParams{
		ComponentID: pgtype.UUID{Bytes: componentUuid, Valid: true},
		Name:        componentFieldParams.Name,
	})
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("AddComponentField error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to create component field",
			})
			return ctx.Send(payload)
		}
	}
	if dbErr == nil && len(componentFields) > 0 {
		log.Printf("AddComponentField error while fetching object, error componentFieldParams already exists\n")
		ctx.Status(http.StatusConflict)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusConflict,
			ErrorMessage: "Component field already exists for this component",
		})
		return ctx.Send(payload)
	}

	componentField := db_gen.CreateComponentFieldAndReturnIdParams{
		ComponentID:  pgtype.UUID{Bytes: componentUuid, Valid: true},
		Name:         componentFieldParams.Name,
		Description:  componentFieldParams.Description,
		DataType:     componentFieldParams.DataType,
		EditorType:   componentFieldParams.EditorType,
		Validation:   componentFieldParams.Validation,
		DefaultValue: componentFieldParams.DefaultValue,
		IsActive:     componentFieldParams.IsActive,
	}

	componentFieldUuid, dbErr := controller.Queries.CreateComponentFieldAndReturnId(controller.DbCtx, componentField)
	if dbErr != nil {
		log.Printf("AddComponentField error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to add component field",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   componentFieldUuid,
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

func (controller *ComponentModel) RemoveComponentField(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	componentFieldUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("RemoveComponentField error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}

	componentUuid, err := uuid.Parse(ctx.Params("componentId"))
	if err != nil {
		log.Printf("RemoveComponentField error parsing component ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid component ID",
		})
		return ctx.Send(payload)
	}

	_, dbErr := controller.Queries.GetComponentFieldById(controller.DbCtx, pgtype.UUID{Bytes: componentFieldUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Component field not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("RemoveComponentField error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to remove component field",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeleteComponentFieldByIdAndComponentId(controller.DbCtx, db_gen.DeleteComponentFieldByIdAndComponentIdParams{
		ID:          pgtype.UUID{Bytes: componentFieldUuid, Valid: true},
		ComponentID: pgtype.UUID{Bytes: componentUuid, Valid: true},
	})
	if dbErr != nil {
		log.Printf("RemoveComponentField error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete component",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("RemoveComponentField error serialising response, error %v\n", marshalErr)
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

func (controller *ComponentModel) ListComponents(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	requestedPage, requestedPageSize, sortBy, sortOrder := api_utils.GetPaginationParams(ctx)
	rowCount, dbErr := controller.Queries.CountComponents(controller.DbCtx)
	if dbErr != nil {
		log.Printf("CountComponents error while retrieving count, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve component count",
		})
		return ctx.Send(payload)
	}

	components := make([]db_gen.Component, 0)
	if sortOrder == "desc" {
		components, dbErr = controller.Queries.ListComponentsDesc(controller.DbCtx, db_gen.ListComponentsDescParams{
			RequestedPage:     int32(requestedPage),
			RequestedPageSize: int32(requestedPageSize),
			SortBy:            sortBy,
		})
	} else {
		components, dbErr = controller.Queries.ListComponentsAsc(controller.DbCtx, db_gen.ListComponentsAscParams{
			RequestedPage:     int32(requestedPage),
			RequestedPageSize: int32(requestedPageSize),
			SortBy:            sortBy,
		})
	}
	if dbErr != nil {
		log.Printf("ListComponents error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve components",
		})
		return ctx.Send(payload)
	}

	componentViews := make([]Component, 0)
	for _, component := range components {
		fields, err := controller.Queries.ListComponentFieldByComponentId(controller.DbCtx, component.ID)
		if err != nil {
			log.Printf("ListComponents error while fetching field information, error %v\n", err)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to retrieve field information",
			})
			return ctx.Send(payload)
		}

		fieldViews := make([]ComponentField, 0)
		for _, field := range fields {
			fieldViews = append(fieldViews, ComponentField{
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

		componentViews = append(componentViews, Component{
			ID:                  component.ID,
			Name:                component.Name,
			Icon:                component.Icon,
			Description:         component.Description,
			ClassName:           component.ClassName,
			HtmlTag:             component.HtmlTag,
			ChildTagConstraints: component.ChildTagConstraints,
			IsActive:            component.IsActive,
			Fields:              fieldViews,
		})
	}

	if componentViews == nil {
		componentViews = make([]Component, 0)
	}
	results := api_utils.GetPaginatedResults(rowCount, requestedPage, requestedPageSize, componentViews)
	payload, marshalErr := json.Marshal(results)
	if marshalErr != nil {
		log.Printf("ListComponents error while marshalling object, error %v\n", marshalErr)
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

func (controller *ComponentModel) GetComponent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	componentUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetComponent error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	component, dbErr := controller.Queries.GetComponentById(controller.DbCtx, pgtype.UUID{Bytes: componentUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Component not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetComponent error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve component",
		})
		return ctx.Send(payload)
	}

	fields, err := controller.Queries.ListComponentFieldByComponentId(controller.DbCtx, component.ID)
	if err != nil {
		log.Printf("GetComponent error while fetching field information, error %v\n", err)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve field information",
		})
		return ctx.Send(payload)
	}

	fieldViews := make([]ComponentField, 0)
	for _, field := range fields {
		fieldViews = append(fieldViews, ComponentField{
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

	componentView := Component{
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

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   componentView,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetComponent error while marshalling object, error %v\n", marshalErr)
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

func (controller *ComponentModel) DeleteComponent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	componentUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("DeleteComponent error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	_, dbErr := controller.Queries.GetComponentById(controller.DbCtx, pgtype.UUID{Bytes: componentUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Component not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("DeleteComponent error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete component",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeleteComponent(controller.DbCtx, pgtype.UUID{Bytes: componentUuid, Valid: true})
	if dbErr != nil {
		log.Printf("DeleteComponent error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete component",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeleteComponentFieldByComponentId(controller.DbCtx, pgtype.UUID{Bytes: componentUuid, Valid: true})
	if dbErr != nil {
		log.Printf("DeleteComponent error while deleting component fields, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete component fields",
		})
		return ctx.Send(payload)

	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("DeleteComponent error serialising response, error %v\n", marshalErr)
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

func (controller *ComponentModel) UpdateComponent(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	componentUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("UpdateComponent error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	sourceComponent, dbErr := controller.Queries.GetComponentById(controller.DbCtx, pgtype.UUID{Bytes: componentUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Component not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdateComponent error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to update component",
		})
		return ctx.Send(payload)
	}

	component := db_gen.UpdateComponentParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &component)
	if unmarshalErr != nil {
		log.Printf("UpdateComponent error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	// Values that are not allowed to be changed when updating component
	component.ID = sourceComponent.ID

	updatedComponent, dbErr := controller.Queries.UpdateComponent(controller.DbCtx, component)
	if dbErr != nil {
		log.Printf("UpdateComponent error while updating object, error %v\n", dbErr)
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
		Data:   updatedComponent,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("UpdateComponent error while marshalling object, error %v\n", marshalErr)
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
