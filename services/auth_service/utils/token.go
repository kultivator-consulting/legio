package utils

import (
	"cortex_api/common"
	"cortex_api/database/db_gen"
	"cortex_api/services/auth_service/config"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

const RandomStringLen = 32

func GetIdAsString(id pgtype.UUID) string {
	return base64.URLEncoding.EncodeToString([]byte(common.FormatIdAsString(id)))
}

func CreateToken(full bool, ttl time.Duration, clientInstanceID pgtype.UUID, account db_gen.Account, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("unable to parse key: %w", err)
	}

	now := time.Now().UTC()

	randomString, err := common.GenerateRandomStringURLSafe(RandomStringLen)

	claims := make(jwt.MapClaims)

	claims["client_id"] = GetIdAsString(clientInstanceID)
	claims["sub"] = GetIdAsString(account.ID)

	if full {
		claims["first_name"] = account.FirstName
		claims["last_name"] = account.LastName
		claims["access_level"] = account.AccessLevel
		claims["credit_balance"] = account.CreditBalance
		claims["username"] = account.Username
		claims["user_email_address"] = account.UserEmailAddress
		claims["user_avatar_url"] = account.UserAvatarUrl
		claims["user_zone_info"] = account.UserZoneInfo
		claims["user_locale"] = account.UserLocale
		claims["last_login"] = account.LastLogin
	}

	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["jti"] = randomString
	claims["aud"] = config.AppConfig().GetAllConfig().BundleId

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS512, claims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return token, nil
}
