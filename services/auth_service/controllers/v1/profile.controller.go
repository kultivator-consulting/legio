package v1

import (
	"context"
	"cortex_api/common"
	"cortex_api/common/api_models"
	"cortex_api/common/api_utils"
	"cortex_api/database/db_gen"
	"cortex_api/services/auth_service/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
)

type ProfileModel struct {
	Queries *db_gen.Queries
	DbCtx   context.Context
}

type ProfileInterface interface {
	GetProfile(ctx *fiber.Ctx) error
	UpdateProfile(ctx *fiber.Ctx) error
	DeleteProfile(ctx *fiber.Ctx) error
}

func ProfileController(
	queries *db_gen.Queries,
	dbCtx context.Context,
) *ProfileModel {
	return &ProfileModel{
		Queries: queries,
		DbCtx:   dbCtx,
	}
}

func (controller *ProfileModel) GetProfile(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	account, dbErr := controller.Queries.GetAccount(controller.DbCtx, ctx.Locals(common.ContextAccount).(pgtype.UUID))
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Profile not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("GetProfile error while retrieving object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to retrieve profile",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   account,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("GetProfile error while marshalling object, error %v\n", marshalErr)
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

func (controller *ProfileModel) DeleteProfile(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	sourceAccount, dbErr := controller.Queries.GetAccount(controller.DbCtx, ctx.Locals(common.ContextAccount).(pgtype.UUID))
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Profile not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("DeleteProfile error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete profile",
		})
		return ctx.Send(payload)
	}
	if sourceAccount.IsSystem {
		log.Printf("DeleteProfile unable to delete a system account\n")
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Unable to delete a system profile",
		})
		return ctx.Send(payload)
	}

	dbErr = controller.Queries.DeleteAccount(controller.DbCtx, sourceAccount.ID)
	if dbErr != nil {
		log.Printf("DeleteProfile error while deleting object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to delete profile",
		})
		return ctx.Send(payload)
	}

	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   nil,
	}
	payload, marshalErr := json.Marshal(success)
	if marshalErr != nil {
		log.Printf("DeleteProfile error serialising response, error %v\n", marshalErr)
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

func (controller *ProfileModel) UpdateProfile(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	sourceAccount, dbErr := controller.Queries.GetAccountById(controller.DbCtx, ctx.Locals(common.ContextAccount).(pgtype.UUID))
	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			ctx.Status(http.StatusNotFound)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusNotFound,
				ErrorMessage: "Profile not found",
			})
			return ctx.Send(payload)
		}
		log.Printf("UpdateProfile error while fetching object, error %v\n", dbErr)
		ctx.Status(http.StatusBadGateway)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadGateway,
			ErrorMessage: "Unable to update profile",
		})
		return ctx.Send(payload)
	}

	account := db_gen.UpdateAccountParams{}
	unmarshalErr := json.Unmarshal(ctx.Body(), &account)
	if unmarshalErr != nil {
		log.Printf("UpdateProfile error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid body",
		})
		return ctx.Send(payload)
	}

	account.ID = sourceAccount.ID
	if account.UserPassword.Valid && account.UserPassword.String != "" {
		account.UserPassword = pgtype.Text{String: utils.HashPassword(account.UserPassword.String), Valid: true}
	} else {
		account.UserPassword = sourceAccount.UserPassword
	}

	if !sourceAccount.IsSystem {
		// Additional values that are not allowed to be changed when updating a profile
		account.Username = sourceAccount.Username
		account.AccessLevel = sourceAccount.AccessLevel
		account.CreditBalance = sourceAccount.CreditBalance
		account.IsLocked = sourceAccount.IsLocked
		account.LastLogin = sourceAccount.LastLogin

		updatedAccount, dbErr := controller.Queries.UpdateAccount(controller.DbCtx, account)
		if dbErr != nil {
			log.Printf("UpdateProfile error while updating object, error %v\n", dbErr)
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
			log.Printf("UpdateProfile error while marshalling object, error %v\n", marshalErr)
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

	updatedAccount, dbErr := controller.Queries.UpdateSystemProfile(controller.DbCtx, db_gen.UpdateSystemProfileParams{ID: account.ID, UserPassword: account.UserPassword})
	if dbErr != nil {
		log.Printf("UpdateProfile error while updating object, error %v\n", dbErr)
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
		log.Printf("UpdateProfile error while marshalling object, error %v\n", marshalErr)
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
