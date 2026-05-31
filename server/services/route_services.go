package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/utils"
)

type CreateServiceRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=128"`
	Token string `json:"token" validate:"required,min=16"`
}

func CreateServiceRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.MemberRoleAdmin) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	var req CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	hash := sha256.Sum256([]byte(req.Token))
	tokenHash := fmt.Sprintf("%x", hash)
	now := time.Now().Unix()
	serviceId, err := utils.GenerateUUIDv7()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to generate service id")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	service := repository.UmiService{
		Id:        serviceId,
		Name:      req.Name,
		TeamId:    teamId,
		TokenHash: tokenHash,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repository.CreateService(&service); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Msg("failed to create service")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusCreated, &service)
}

func ListServicesRoute(w http.ResponseWriter, r *http.Request) {
	teamId := route_data.GetTeamId(r.Context())
	services, err := repository.ListServicesByTeamId(teamId)
	if err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Msg("failed to list services")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, services)
}

func GetServiceRoute(w http.ResponseWriter, r *http.Request) {
	teamId := route_data.GetTeamId(r.Context())
	serviceId := route_data.GetServiceId(r.Context())
	service, err := repository.GetServiceById(serviceId)
	if err != nil || service == nil || service.TeamId != teamId {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.ErrorCodeNotFound)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, service)
}

type UpdateServiceRequest struct {
	Name string `json:"name" validate:"required,min=1,max=128"`
}

func UpdateServiceRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.MemberRoleAdmin) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	serviceId := route_data.GetServiceId(r.Context())
	service, err := repository.GetServiceById(serviceId)
	if err != nil || service == nil || service.TeamId != teamId {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.ErrorCodeNotFound)
		return
	}
	var req UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := repository.UpdateServiceName(serviceId, req.Name); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Str("serviceId", serviceId).Msg("failed to update service name")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	service.Name = req.Name
	utils.WriteHttpJsonResponse(w, http.StatusOK, service)
}

type RotateServiceTokenRequest struct {
	Token string `json:"token" validate:"required,min=16"`
}

func RotateServiceTokenRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.MemberRoleAdmin) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	serviceId := route_data.GetServiceId(r.Context())
	service, err := repository.GetServiceById(serviceId)
	if err != nil || service == nil || service.TeamId != teamId {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.ErrorCodeNotFound)
		return
	}
	var req RotateServiceTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	hash := sha256.Sum256([]byte(req.Token))
	tokenHash := fmt.Sprintf("%x", hash)
	if err := repository.UpdateServiceTokenHash(serviceId, tokenHash); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Str("serviceId", serviceId).Msg("failed to rotate service token")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, nil)
}

func DeleteServiceRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.MemberRoleAdmin) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	serviceId := route_data.GetServiceId(r.Context())
	service, err := repository.GetServiceById(serviceId)
	if err != nil || service == nil || service.TeamId != teamId {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.ErrorCodeNotFound)
		return
	}
	if err := repository.DeleteService(serviceId); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Str("serviceId", serviceId).Msg("failed to delete service")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, nil)
}
