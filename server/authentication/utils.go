package authentication

import (
	"net/http"
	"strings"

	"zyrouge.me/umi/constants"
)

func ExtractBearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}

func ExtractAccessToken(r *http.Request) string {
	cookieValue, _ := GetCookieValue(r, constants.HttpCookieAccessToken)
	if cookieValue != "" {
		return cookieValue
	}
	return ExtractBearerToken(r)
}

func SetSecureCookie(w http.ResponseWriter, name string, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	})
}

func GetCookieValue(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
