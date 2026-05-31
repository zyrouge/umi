package authentication

import (
	"net/http"

	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/utils"
)

func MemberPermissionMiddleware(next http.Handler, permission repository.UmiMemberPermission) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := route_data.GetTeamRole(r.Context())
		if !role.HasPermission(permission) {
			utils.WriteHttpJsonError(w, http.StatusForbidden, constants.ErrorCodeForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
