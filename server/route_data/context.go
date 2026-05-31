package route_data

import (
	"context"

	"zyrouge.me/umi/repository"
)

type RouteDataContextKey string

const (
	ContextKeyUserIdContextKey      RouteDataContextKey = "user_id"
	ContextKeyTeamIdContextKey      RouteDataContextKey = "team_id"
	ContextKeyMemberRoleContextKey  RouteDataContextKey = "member_role"
	ContextKeyChannelIdContextKey   RouteDataContextKey = "channel_id"
	ContextKeyServiceIdContextKey   RouteDataContextKey = "service_id"
	ContextKeyServiceNameContextKey RouteDataContextKey = "service_name"
	ContextKeyMemberUserIdContextKey RouteDataContextKey = "member_user_id"
)

func WithUserId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ContextKeyUserIdContextKey, value)
}

func GetUserId(ctx context.Context) string {
	value, _ := ctx.Value(ContextKeyUserIdContextKey).(string)
	return value
}

func WithTeamId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ContextKeyTeamIdContextKey, value)
}

func GetTeamId(ctx context.Context) string {
	value, _ := ctx.Value(ContextKeyTeamIdContextKey).(string)
	return value
}

func WithMemberRole(ctx context.Context, role repository.UmiMemberRole) context.Context {
	return context.WithValue(ctx, ContextKeyMemberRoleContextKey, role)
}

func GetTeamRole(ctx context.Context) repository.UmiMemberRole {
	role, _ := ctx.Value(ContextKeyMemberRoleContextKey).(repository.UmiMemberRole)
	return role
}

func WithChannelId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ContextKeyChannelIdContextKey, value)
}

func GetChannelId(ctx context.Context) string {
	value, _ := ctx.Value(ContextKeyChannelIdContextKey).(string)
	return value
}

func WithServiceId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ContextKeyServiceIdContextKey, value)
}

func GetServiceId(ctx context.Context) string {
	value, _ := ctx.Value(ContextKeyServiceIdContextKey).(string)
	return value
}

func WithServiceName(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ContextKeyServiceNameContextKey, value)
}

func GetServiceName(ctx context.Context) string {
	value, _ := ctx.Value(ContextKeyServiceNameContextKey).(string)
	return value
}

func WithMemberUserId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ContextKeyMemberUserIdContextKey, value)
}

func GetMemberUserId(ctx context.Context) string {
	value, _ := ctx.Value(ContextKeyMemberUserIdContextKey).(string)
	return value
}
