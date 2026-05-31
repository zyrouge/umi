package authentication

import (
	"context"
	"net/http"

	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/utils"
)

func GetUserId(ctx context.Context) string {
	return route_data.GetUserId(ctx)
}

func GetTeamRole(ctx context.Context) repository.UmiMemberRole {
	return route_data.GetTeamRole(ctx)
}

func GetServiceId(ctx context.Context) string {
	return route_data.GetServiceId(ctx)
}

func GetServiceName(ctx context.Context) string {
	return route_data.GetServiceName(ctx)
}

func RequirePermissionMiddleware(w http.ResponseWriter, r *http.Request, minRole repository.UmiMemberRole) bool {
	role := route_data.GetTeamRole(r.Context())
	if !repository.HasMinRole(role, minRole) {
		utils.WriteHttpJsonError(w, http.StatusForbidden, constants.UmiErrorCodeForbidden)
		return false
	}
	return true
}
