package v1

import (
	"bytes"
	"context"
	"cortex_api/common"
	"cortex_api/common/api_models"
	"cortex_api/common/api_utils"
	"cortex_api/database/db_gen"
	"cortex_api/services/auth_service/config"
	"cortex_api/services/auth_service/models"
	"cortex_api/services/auth_service/utils"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"text/template"
	"time"
)

const MaxSecureDelay = 5
const WebAccessExpiryMinutes = 15
const WebRefreshExpiryMinutes = 60

type AuthModel struct {
	Queries *db_gen.Queries
	DbCtx   context.Context
}

type AuthInterface interface {
	AuthLogin(ctx *fiber.Ctx) error
	AuthForgotPassword(ctx *fiber.Ctx) error
	AuthLogout(ctx *fiber.Ctx) error
	AuthRefresh(ctx *fiber.Ctx) error
	AuthValidate(ctx *fiber.Ctx) error
}

func AuthController(
	queries *db_gen.Queries,
	dbCtx context.Context,
) *AuthModel {
	return &AuthModel{
		Queries: queries,
		DbCtx:   dbCtx,
	}
}

func secureDelay() {
	secureDelay := rand.Intn(MaxSecureDelay)
	time.Sleep(time.Second * time.Duration(secureDelay))
}

func createSession(
	controller *AuthModel,
	accountId pgtype.UUID,
	clientId pgtype.UUID,
	clientAgent string,
	clientIp string,
	clientBundleId string,
	accessToken string,
	refreshToken string,
	accessExpiresIn time.Duration,
	refreshExpiresIn time.Duration,
	isDeviceApp bool) (db_gen.Session, error) {

	session := db_gen.CreateSessionParams{
		AccountID:          accountId,
		ClientID:           clientId,
		ClientAgent:        clientAgent,
		ClientIp:           clientIp,
		ClientBundleID:     pgtype.Text{String: clientBundleId, Valid: true},
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpiry:  pgtype.Timestamptz{Time: time.Now().Add(accessExpiresIn), Valid: true},
		RefreshTokenExpiry: pgtype.Timestamptz{Time: time.Now().Add(refreshExpiresIn), Valid: true},
		IsDeviceApp:        isDeviceApp,
	}

	return controller.Queries.CreateSession(controller.DbCtx, session)
}

func setCookie(
	cookieName string,
	cookieValue string,
	cookieExpires time.Duration,
	cookieDomain string,
	cookiePath string,
	cookieSecure bool,
	cookieHttpOnly bool) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = cookieName
	cookie.Value = cookieValue
	cookie.HTTPOnly = cookieHttpOnly
	cookie.Expires = time.Now().Add(cookieExpires)
	cookie.Domain = cookieDomain
	cookie.Secure = cookieSecure
	cookie.Path = cookiePath
	return cookie
}

func invalidateCookie(
	cookieName string,
	cookieDomain string,
	cookiePath string,
	cookieSecure bool,
	cookieHttpOnly bool) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = cookieName
	cookie.Value = "deleted"
	cookie.HTTPOnly = cookieHttpOnly
	cookie.Expires = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	cookie.Domain = cookieDomain
	cookie.Secure = cookieSecure
	cookie.Path = cookiePath
	return cookie
}

func getSafeClientInstanceId(clientId string) (pgtype.UUID, error) {
	clientUuid, err := uuid.Parse(clientId)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("could not parse device ID: %w", err)
	}

	return pgtype.UUID{Bytes: clientUuid, Valid: true}, nil
}

func buildSession(
	ctx *fiber.Ctx,
	controller *AuthModel,
	account db_gen.Account,
	clientId string,
	isDeviceApp bool,
) error {
	clientUuid, err := getSafeClientInstanceId(clientId)
	if err != nil {
		log.Printf("AuthLogin error while generating client instance id, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	accessExpiresIn, err := time.ParseDuration(config.AppConfig().GetAllConfig().AccessTokenExpiresIn)
	if err != nil {
		log.Printf("AuthLogin error while parsing access token expiration, error %v\n", err)
		accessExpiresIn = time.Second * WebAccessExpiryMinutes
	}

	accessToken, err := utils.CreateToken(true, accessExpiresIn, clientUuid, account, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		log.Printf("AuthLogin error while creating access token, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	refreshExpiresIn, err := time.ParseDuration(config.AppConfig().GetAllConfig().RefreshTokenExpiresIn)
	if err != nil {
		log.Printf("AuthLogin error while parsing refresh token expiration, error %v\n", err)
		refreshExpiresIn = time.Second * WebRefreshExpiryMinutes
	}

	refreshToken, err := utils.CreateToken(false, refreshExpiresIn, clientUuid, account, os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"))
	if err != nil {
		log.Printf("AuthLogin error while creating refresh token, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	if !isDeviceApp {
		ctx.Cookie(setCookie(common.RefreshTokenCookie, refreshToken, refreshExpiresIn, config.AppConfig().GetAllConfig().CookieDomain, "/", false, true))
	}

	clientAgent := ctx.Get("User-Agent")
	clientIp := ctx.IP()
	clientBundleId := ctx.Get("X-Bundle-Id")

	previousSessions, err := controller.Queries.GetSessionsByAccountIdClientId(controller.DbCtx, db_gen.GetSessionsByAccountIdClientIdParams{
		AccountID: account.ID,
		ClientID:  clientUuid,
	})
	if err != nil {
		log.Printf("AuthLogin error while finding previous session, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	var sessionInstance db_gen.Session
	if len(previousSessions) != 0 {
		log.Printf("AuthLogin two or more sessions found for this account, deleting previous sessions.\n")
		if err := controller.Queries.DeleteSessionsByAccountIdClientId(controller.DbCtx, db_gen.DeleteSessionsByAccountIdClientIdParams{
			AccountID: account.ID,
			ClientID:  clientUuid,
		}); err != nil {
			log.Printf("AuthLogin error while deleting session, error %v\n", err)
			ctx.Status(http.StatusUnauthorized)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "Invalid login credentials, check your username and password and try again",
			})
			secureDelay()
			return ctx.Send(payload)
		}
	}

	sessionInstance, err = createSession(
		controller,
		account.ID,
		clientUuid,
		clientAgent,
		clientIp,
		clientBundleId,
		accessToken,
		refreshToken,
		accessExpiresIn,
		refreshExpiresIn,
		isDeviceApp)
	if err != nil {
		log.Printf("AuthLogin error while creating session, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	mobileRefreshToken := ""
	if isDeviceApp {
		mobileRefreshToken = refreshToken
	}
	response := models.LocalLoginResponse{
		ID:           utils.GetIdAsString(sessionInstance.ID),
		AccountID:    utils.GetIdAsString(account.ID),
		AccessToken:  accessToken,
		RefreshToken: mobileRefreshToken,
	}
	success := api_models.Success{
		Status: api_utils.SuccessResponse,
		Data:   response,
	}
	payload, err := json.Marshal(success)
	if err != nil {
		log.Printf("AuthLogin error while marshalling session response, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	ctx.Status(http.StatusOK)
	return ctx.Send(payload)
}

func (controller *AuthModel) AuthLogin(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	log.Printf("AuthLogin incoming login request from %v\n", ctx.IP())

	var credentials models.LocalLoginRequest
	if unmarshalErr := json.Unmarshal(ctx.Body(), &credentials); unmarshalErr != nil {
		log.Printf("AuthLogin error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	account, err := controller.Queries.GetAccountByUsername(controller.DbCtx, credentials.Username)
	if err != nil {
		log.Printf("AuthLogin error while finding object, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	if err := utils.ComparePassword(account.UserPassword.String, credentials.Password); err != nil {
		log.Printf("AuthLogin error while comparing password, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	if err := controller.Queries.UpdateAccountLastLogin(controller.DbCtx, account.ID); err != nil {
		log.Printf("AuthLogin error while updating last login, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	isDeviceApp, err := strconv.ParseBool(ctx.Get(common.IsDeviceApp))
	if err != nil {
		isDeviceApp = false
	}

	return buildSession(ctx, controller, account, credentials.ClientId, isDeviceApp)
}

func (controller *AuthModel) AuthForgotPassword(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	log.Printf("AuthForgotPassword incoming forgot password request from %v\n", ctx.IP())

	var forgotPassword models.LocalForgotPasswordRequest
	if unmarshalErr := json.Unmarshal(ctx.Body(), &forgotPassword); unmarshalErr != nil {
		log.Printf("AuthForgotPassword error while unmarshalling object, error %v\n", unmarshalErr.Error())
		secureDelay()
		ctx.Status(http.StatusOK)
		return ctx.Send(nil)
	}

	account, err := controller.Queries.GetAccountByEmailAddress(controller.DbCtx, forgotPassword.Email)
	if err != nil {
		log.Printf("AuthForgotPassword error while finding object, error %v\n", err)
		secureDelay()
		ctx.Status(http.StatusOK)
		return ctx.Send(nil)
	}

	resetPasswordToken, err := common.GenerateRandomStringURLSafe(60)
	if err != nil {
		log.Printf("AuthForgotPassword error while generating reset token, error %v\n", err)
		secureDelay()
		ctx.Status(http.StatusOK)
		return ctx.Send(nil)
	}

	passwordResetDuration, err := time.ParseDuration(os.Getenv("FORGOT_PASSWORD_EXPIRY"))
	passwordResetExpires := time.Now().Add(passwordResetDuration)

	_, err = controller.Queries.UpdateAccountResetPasswordToken(controller.DbCtx, db_gen.UpdateAccountResetPasswordTokenParams{
		ID:                       account.ID,
		ResetPasswordToken:       pgtype.Text{String: resetPasswordToken, Valid: true},
		ResetPasswordTokenExpiry: pgtype.Timestamptz{Time: passwordResetExpires, Valid: true},
	})
	if err != nil {
		log.Printf("AuthForgotPassword error while updating reset token, error %v\n", err)
		secureDelay()
		ctx.Status(http.StatusOK)
		return ctx.Send(nil)
	}

	from := os.Getenv("SMTP_FROM")
	subject := os.Getenv("FORGOT_PASSWORD_SUBJECT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	to := []string{
		account.UserEmailAddress,
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	templateFile := os.Getenv("FORGOT_PASSWORD_TEMPLATE")
	log.Printf("AuthForgotPassword using template file %v\n", templateFile)
	htmlTemplate, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Printf("AuthForgotPassword error while parsing template file %v, error %v\n", htmlTemplate, err)
		secureDelay()
		ctx.Status(http.StatusOK)
		return ctx.Send(nil)
	}
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	date := time.Now().Format(time.RFC822Z)
	body.Write([]byte(fmt.Sprintf("Date: %s\nFrom: Legio Admin<%s>\nSubject: %s\nTo: %s %s<%s>\n%s\n\n",
		date,
		from,
		subject,
		account.FirstName.String,
		account.LastName.String,
		to,
		mimeHeaders,
	)))

	err = htmlTemplate.Execute(&body, struct {
		FirstName string
		LastName  string
		ResetLink string
	}{
		FirstName: account.FirstName.String,
		LastName:  account.LastName.String,
		ResetLink: os.Getenv("FORGOT_PASSWORD_URL") + resetPasswordToken,
	})
	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		log.Printf("AuthForgotPassword error while sending email, error %v\n", err)
		secureDelay()
		ctx.Status(http.StatusOK)
		return ctx.Send(nil)
	}

	log.Printf("AuthForgotPassword successfully sent email\n")

	ctx.Status(http.StatusOK)
	return ctx.Send(nil)
}

func (controller *AuthModel) AuthResetPassword(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	log.Printf("AuthResetPassword incoming forgot password request from %v\n", ctx.IP())

	var resetPassword models.LocalResetPasswordRequest
	if unmarshalErr := json.Unmarshal(ctx.Body(), &resetPassword); unmarshalErr != nil {
		log.Printf("AuthResetPassword error while unmarshalling object, error %v\n", unmarshalErr.Error())
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid request, please try again later",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	if resetPassword.Password != resetPassword.ConfirmPassword {
		log.Printf("AuthResetPassword password and confirm password do not match\n")
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Passwords do not match, please try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	_, err := controller.Queries.UpdateAccountResetPasswordByToken(controller.DbCtx, db_gen.UpdateAccountResetPasswordByTokenParams{
		ResetPasswordToken: pgtype.Text{String: resetPassword.Token, Valid: true},
		UserPassword:       pgtype.Text{String: utils.HashPassword(resetPassword.Password), Valid: true},
	})
	if err != nil {
		log.Printf("AuthResetPassword error while finding object, error %v\n", err)
		ctx.Status(http.StatusBadRequest)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusBadRequest,
			ErrorMessage: "Invalid request, please try again later",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	ctx.Cookie(invalidateCookie(common.RefreshTokenCookie, config.AppConfig().GetAllConfig().CookieDomain, "/", false, true))

	ctx.Status(http.StatusOK)
	return ctx.Send(nil)
}

func (controller *AuthModel) AuthLogout(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	if err := controller.Queries.DeleteSessionsByAccountIdClientId(controller.DbCtx, db_gen.DeleteSessionsByAccountIdClientIdParams{
		AccountID: ctx.Locals(common.ContextAccount).(pgtype.UUID),
		ClientID:  ctx.Locals(common.ContextSession).(pgtype.UUID),
	}); err != nil {
		log.Printf("AuthLogin error while deleting session, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid login credentials, check your username and password and try again",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	ctx.Cookie(invalidateCookie(common.RefreshTokenCookie, config.AppConfig().GetAllConfig().CookieDomain, "/", false, true))

	ctx.Status(http.StatusOK)
	return ctx.Send(nil)
}

func (controller *AuthModel) AuthRefresh(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	checkDeviceAppSession := false
	token := ctx.Cookies(common.RefreshTokenCookie, "")
	if token == "" {
		checkDeviceAppSession = true
		token = common.ExtractToken(ctx)
		if token == "" {
			log.Printf("AuthRefresh invalid refresh token\n")
			ctx.Status(http.StatusUnauthorized)
			payload, _ := json.Marshal(api_models.Failed{
				Status:       api_utils.FailedResponse,
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "Invalid session, please try logging in again later",
			})
			secureDelay()
			return ctx.Send(payload)
		}
	}

	publicKey := os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	_, clientId, err := common.ValidateToken(token, publicKey)
	if err != nil {
		log.Printf("AuthRefresh error validating refresh token, %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid session, please try logging in again later",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	sessions, err := controller.Queries.GetSessionsByRefreshToken(controller.DbCtx, token)
	if err != nil {
		log.Printf("AuthRefresh error while validating existing session, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid session, please try logging in again later",
		})
		secureDelay()
		return ctx.Send(payload)
	}
	if len(sessions) != 1 {
		log.Printf("AuthRefresh no matching session for this token\n")
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid session, please try logging in again later",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	session := sessions[0]
	if session.IsDeviceApp != checkDeviceAppSession {
		log.Printf("AuthRefresh misconfigured session token, check isDeviceApp flag matches initial login request\n")
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorCode:    http.StatusUnauthorized,
			ErrorMessage: "Invalid session, please try logging in again later",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	account, err := controller.Queries.GetUnlockedAccountById(controller.DbCtx, session.AccountID)
	if err != nil {
		log.Printf("AuthRefresh error while validating existing session, error %v\n", err)
		ctx.Status(http.StatusUnauthorized)
		payload, _ := json.Marshal(api_models.Failed{
			Status:       api_utils.FailedResponse,
			ErrorMessage: "Invalid session, please try logging in again later",
		})
		secureDelay()
		return ctx.Send(payload)
	}

	return buildSession(ctx, controller, account, common.FormatIdAsString(clientId.(pgtype.UUID)), session.IsDeviceApp)
}

func (controller *AuthModel) AuthValidate(ctx *fiber.Ctx) error {
	api_utils.SetRESTHeaders(ctx)

	log.Printf("AuthValidate incoming login request from %v\n", ctx.IP())

	ctx.Status(http.StatusOK)
	return ctx.Send(nil)
}
