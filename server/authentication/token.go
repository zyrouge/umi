package authentication

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

type UmiAccessTokenClaims struct {
	UserId     string  `json:"user_id"`
	TeamId     *string `json:"team_id"`
	MemberRole *string `json:"member_role"`
	jwt.RegisteredClaims
}

const AccessTokenTTL = 15 * time.Minute
const RefreshTokenTTL = 24 * time.Hour

func GenerateAccessToken(userId string, teamId *string, memberRole *repository.UmiMemberRole, secret []byte) (string, error) {
	var memberRoleStr *string
	if memberRole != nil {
		str := string(*memberRole)
		memberRoleStr = &str
	}
	claims := UmiAccessTokenClaims{
		UserId:     userId,
		TeamId:     teamId,
		MemberRole: memberRoleStr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signed, nil
}

func ValidateAccessToken(token string, secret []byte) (*UmiAccessTokenClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &UmiAccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := jwtToken.Claims.(*UmiAccessTokenClaims)
	if !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func GenerateRefreshToken(userId string) (string, error) {
	id, err := utils.GenerateUUIDv7()
	if err != nil {
		return "", err
	}
	tokenBytes, err := utils.GenerateRandomBytes(32)
	if err != nil {
		return "", err
	}
	secret := utils.BytesToBase64(tokenBytes)
	token := fmt.Sprintf("%s.%s", id, secret)
	hashBytes := sha256.Sum256([]byte(secret))
	hash := utils.BytesToHex(hashBytes[:])
	now := time.Now()
	tokenData := repository.UmiRefreshToken{
		Id:        id,
		UserId:    userId,
		TokenHash: hash,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(RefreshTokenTTL).Unix(),
	}
	if err := repository.InsertRefreshToken(&tokenData); err != nil {
		utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to insert refresh token")
		return "", err
	}
	return token, nil
}

func ValidateRefreshToken(token string) (*repository.UmiRefreshToken, error) {
	parts := strings.SplitN(token, ".", 2)
	id := parts[0]
	hash := parts[1]
	tokenData, err := repository.GetRefreshTokenById(id)
	if err != nil {
		return nil, err
	}
	if tokenData == nil || tokenData.TokenHash != hash || !tokenData.IsValid() {
		return nil, fmt.Errorf("invalid refresh token")
	}
	return tokenData, nil
}
