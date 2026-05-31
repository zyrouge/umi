package route_data

import (
	"context"

	"zyrouge.me/umi/repository"
)

type RouteDataContextKey string

const (
	RouteDataContextKeyUserId      RouteDataContextKey = "user_id"
	RouteDataContextKeyTeamId      RouteDataContextKey = "team_id"
	RouteDataContextKeyMemberRole  RouteDataContextKey = "member_role"
	RouteDataContextKeyChannelId   RouteDataContextKey = "channel_id"
	RouteDataContextKeyServiceId   RouteDataContextKey = "service_id"
	RouteDataContextKeyServiceName RouteDataContextKey = "service_name"
	RouteDataContextKeyMemberUserId RouteDataContextKey = "member_user_id"
)

func WithUserId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, RouteDataContextKeyUserId, value)
}

func GetUserId(ctx context.Context) string {
	value, _ := ctx.Value(RouteDataContextKeyUserId).(string)
	return value
}

func WithTeamId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, RouteDataContextKeyTeamId, value)
}

func GetTeamId(ctx context.Context) string {
	value, _ := ctx.Value(RouteDataContextKeyTeamId).(string)
	return value
}

func WithMemberRole(ctx context.Context, role repository.UmiMemberRole) context.Context {
	return context.WithValue(ctx, RouteDataContextKeyMemberRole, role)
}

func GetTeamRole(ctx context.Context) repository.UmiMemberRole {
	role, _ := ctx.Value(RouteDataContextKeyMemberRole).(repository.UmiMemberRole)
	return role
}

func WithChannelId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, RouteDataContextKeyChannelId, value)
}

func GetChannelId(ctx context.Context) string {
	value, _ := ctx.Value(RouteDataContextKeyChannelId).(string)
	return value
}

func WithServiceId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, RouteDataContextKeyServiceId, value)
}

func GetServiceId(ctx context.Context) string {
	value, _ := ctx.Value(RouteDataContextKeyServiceId).(string)
	return value
}

func WithServiceName(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, RouteDataContextKeyServiceName, value)
}

func GetServiceName(ctx context.Context) string {
	value, _ := ctx.Value(RouteDataContextKeyServiceName).(string)
	return value
}

func WithMemberUserId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, RouteDataContextKeyMemberUserId, value)
}

func GetMemberUserId(ctx context.Context) string {
	value, _ := ctx.Value(RouteDataContextKeyMemberUserId).(string)
	return value
}
