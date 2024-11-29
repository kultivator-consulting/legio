package v1

import (
	"context"
	"cortex_api/common/api_models"
	"cortex_api/common/api_utils"
	"cortex_api/database/db_gen"
	"cortex_api/services/auth_service/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
)

type AccountModel struct {
	Queries *db_gen.Queries
	DbCtx   context.Context
}

type AccountInterface interface {
	CreateAccount(ctx *fiber.Ctx) error
	ListAccounts(ctx *fiber.Ctx) error
	GetAccount(ctx *fiber.Ctx) error
	DeleteAccount(ctx *fiber.Ctx) error
	UpdateAccount(ctx *fiber.Ctx) error
}

func AccountController(
	queries *db_gen.Queries,
	dbCtx context.Context,
) *AccountModel {
	return &AccountModel{
		Queries: queries,
		DbCtx:   dbCtx,
	}
}

func (controller *AccountModel) CreateAccount(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	account := db_gen.CreateAccountAndReturnIdParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &account)
	if unmarshalErr != nil {
		log.Printf("CreateAccount error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	_, dbErr := controller.Queries.GetAccountByUsername(controller.DbCtx, account.Username)
	if dbErr != nil {
		if dbErr.Error() != "no rows in result set" {
			log.Printf("CreateAccount error while fetching object, error %v\n", dbErr)
			ctx.Status(http.StatusBadGateway)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusBadGateway,
				ErrorMessage: "Unable to create account",
			})
			return ctx.Send(payload)
		}
	}
	if dbErr == nil {
		log.Printf("CreateAccount error while fetching object, error account already exists\n")
		ctx.Status(http.StatusConflict)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusConflict,
			ErrorMessage: "Account already exists",
		})
		return ctx.Send(payload)
	}

	if account.UserPassword.Valid && account.UserPassword.String != "" {
		account.UserPassword = pgtype.Text{String: utils.HashPassword(account.UserPassword.String), Valid: true}
	} else {
		account.UserPassword = pgtype.Text{String: "", Valid: true}
		account.IsLocked = true
	}

	accountUuid, dbErr := controller.Queries.CreateAccountAndReturnId(controller.DbCtx, account)
	if dbErr != nil {
		log.Printf("CreateAccount error while inserting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to create account",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   accountUuid,
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

func (controller *AccountModel) ListAccounts(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	requestedPage, requestedPageSize, sortBy, sortOrder := api_utils.GetPaginationParams(ctx)
	rowCount, dbErr := controller.Queries.CountAccounts(controller.DbCtx)
	if dbErr != nil {
		log.Printf("CountAccounts error while retrieving count, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve account count",
		})
		return ctx.Send(payload)
	}

	accounts, dbErr := controller.Queries.ListAccounts(controller.DbCtx, db_gen.ListAccountsParams{
		RequestedPage:     int32(requestedPage),
		RequestedPageSize: int32(requestedPageSize),
		SortBy:            sortBy,
		SortOrder:         sortOrder,
	})
	if dbErr != nil {
		log.Printf("ListAccounts error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve accounts",
		})
		return ctx.Send(payload)
	}

	if accounts == nil {
		accounts = make([]db_gen.ListAccountsRow, 0)
	}
	results := api_utils.GetPaginatedResults(rowCount, requestedPage, requestedPageSize, accounts)
	payload, marshalErr := json.Marshal(results)
	if marshalErr != nil {
		log.Printf("ListAccounts error while marshalling object, error %v\n", marshalErr)
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

func (controller *AccountModel) GetAccount(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	accountUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("GetAccount error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	account, dbErr := controller.Queries.GetAccount(controller.DbCtx, pgtype.UUID{Bytes: accountUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Account not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetAccount error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve account",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   account,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetAccount error while marshalling object, error %v\n", marshalErr)
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

func (controller *AccountModel) DeleteAccount(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	accountUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("DeleteAccount error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	sourceAccount, dbErr := controller.Queries.GetAccount(controller.DbCtx, pgtype.UUID{Bytes: accountUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Account not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("DeleteAccount error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete account",
		})
		return ctx.Send(payload)
	}
	if sourceAccount.IsSystem {
		log.Printf("DeleteAccount unable to delete a system account\n")
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Unable to delete a system account",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeleteAccount(controller.DbCtx, pgtype.UUID{Bytes: accountUuid, Valid: true})
	if dbErr != nil {
		log.Printf("DeleteAccount error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete account",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("DeleteAccount error serialising response, error %v\n", marshalErr)
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

func (controller *AccountModel) UpdateAccount(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	accountUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		log.Printf("UpdateAccount error parsing ID, error %v\n", err.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid ID",
		})
		return ctx.Send(payload)
	}
	sourceAccount, dbErr := controller.Queries.GetAccountById(controller.DbCtx, pgtype.UUID{Bytes: accountUuid, Valid: true})
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Account not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdateAccount error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to update account",
		})
		return ctx.Send(payload)
	}
	if sourceAccount.IsSystem {
		log.Printf("UpdateAccount unable to update a system account\n")
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Unable to update a system account",
		})
		return ctx.Send(payload)
	}

	account := db_gen.UpdateAccountParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &account)
	if unmarshalErr != nil {
		log.Printf("UpdateAccount error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}
	account.ID = pgtype.UUID{Bytes: accountUuid, Valid: true}
	if account.UserPassword.Valid && account.UserPassword.String != "" {
		account.UserPassword = pgtype.Text{String: utils.HashPassword(account.UserPassword.String), Valid: true}
	} else {
		if sourceAccount.UserPassword.Valid && sourceAccount.UserPassword.String == "" {
			account.IsLocked = true
		} else {
			account.UserPassword = sourceAccount.UserPassword
		}
	}
	if !account.UserAvatarUrl.Valid {
		account.UserAvatarUrl = sourceAccount.UserAvatarUrl
	}
	if !account.UserZoneInfo.Valid {
		account.UserZoneInfo = sourceAccount.UserZoneInfo
	}
	if !account.UserLocale.Valid {
		account.UserLocale = sourceAccount.UserLocale
	}

	// Values that are not allowed to be changed when updating an account
	account.ID = sourceAccount.ID
	account.Username = sourceAccount.Username
	account.LastLogin = sourceAccount.LastLogin

	updatedAccount, dbErr := controller.Queries.UpdateAccount(controller.DbCtx, account)
	if dbErr != nil {
		log.Printf("UpdateAccount error while updating object, error %v\n", dbErr)
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
		Data:   updatedAccount,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("UpdateAccount error while marshalling object, error %v\n", marshalErr)
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
