package teams

import (
	"encoding/json"
	"net/http"
	"time"

	"zyrouge.me/umi/application"
	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/utils"
)

func ListMembersRoute(w http.ResponseWriter, r *http.Request) {
	teamId := route_data.GetTeamId(r.Context())
	members, err := repository.ListMembersByTeamId(teamId)
	if err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Msg("failed to list members")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, members)
}

type AddMemberRequest struct {
	UserId string                   `json:"user_id" validate:"required"`
	Role   repository.UmiMemberRole `json:"role" validate:"required,oneof=admin member"`
}

func AddMemberRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.MemberRoleAdmin) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	config, err := application.GetConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to get config")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	user, err := repository.GetUserById(req.UserId, config.Secret.UserEncryptionKeyBytes)
	if err != nil || user == nil {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.ErrorCodeNotFound)
		return
	}
	existing, err := repository.GetMember(req.UserId, teamId)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", req.UserId).Str("teamId", teamId).Msg("failed to get existing member")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	if existing != nil {
		utils.WriteHttpJsonError(w, http.StatusConflict, constants.ErrorCodeConflict)
		return
	}
	now := time.Now().Unix()
	member := repository.UmiMember{
		UserId: req.UserId, TeamId: teamId, Role: req.Role, CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.InsertMember(&member); err != nil {
		utils.Logger.Error().Err(err).Str("userId", req.UserId).Str("teamId", teamId).Msg("failed to insert member")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusCreated, &member)
}

type UpdateMemberRoleRequest struct {
	Role repository.UmiMemberRole `json:"role" validate:"required,oneof=admin member"`
}

func UpdateMemberRoleRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.MemberRoleOwner) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	targetUserId := route_data.GetMemberUserId(r.Context())
	var req UpdateMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := repository.UpdateMemberRole(targetUserId, teamId, req.Role); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Str("targetUserId", targetUserId).Msg("failed to update member role")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	member, _ := repository.GetMember(targetUserId, teamId)
	utils.WriteHttpJsonResponse(w, http.StatusOK, member)
}

func RemoveMemberRoute(w http.ResponseWriter, r *http.Request) {
	teamId := route_data.GetTeamId(r.Context())
	targetUserId := route_data.GetMemberUserId(r.Context())
	callerUserId := authentication.GetUserId(r.Context())
	callerRole := authentication.GetTeamRole(r.Context())
	if callerUserId != targetUserId && !repository.HasMinRole(callerRole, repository.MemberRoleAdmin) {
		utils.WriteHttpJsonError(w, http.StatusForbidden, constants.ErrorCodeForbidden)
		return
	}
	target, err := repository.GetMember(targetUserId, teamId)
	if err != nil || target == nil {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.ErrorCodeNotFound)
		return
	}
	if target.Role == repository.MemberRoleOwner && callerUserId != targetUserId {
		utils.WriteHttpJsonError(w, http.StatusForbidden, constants.ErrorCodeForbidden)
		return
	}
	if err := repository.DeleteMemberByUserIdAndTeamId(targetUserId, teamId); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Str("targetUserId", targetUserId).Msg("failed to delete member")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, nil)
}
