package teams

import (
	"net/http"

	"github.com/gorilla/mux"
	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/services"
)

func AttachRoutes(router *mux.Router) error {
	r := router.NewRoute().Subrouter()
	r.Use(authentication.AuthMiddleware)

	r.Path("").Methods(http.MethodPost).HandlerFunc(CreateTeamRoute)
	r.Path("").Methods(http.MethodGet).HandlerFunc(ListTeamsRoute)

	teamRouter := r.PathPrefix("/{teamId}").Subrouter()
	teamRouter.Use(route_data.TeamIdMiddleware)
	teamRouter.Use(authentication.TeamMiddleware)

	teamRouter.Path("").Methods(http.MethodGet).HandlerFunc(GetTeamRoute)
	teamRouter.Path("").Methods(http.MethodPut).HandlerFunc(UpdateTeamRoute)
	teamRouter.Path("").Methods(http.MethodDelete).HandlerFunc(DeleteTeamRoute)

	teamRouter.Path("/members").Methods(http.MethodGet).HandlerFunc(ListMembersRoute)
	teamRouter.Path("/members").Methods(http.MethodPost).HandlerFunc(AddMemberRoute)

	memberRouter := teamRouter.PathPrefix("/members/{userId}").Subrouter()
	memberRouter.Use(route_data.MemberUserIdMiddleware)
	memberRouter.Path("").Methods(http.MethodPut).HandlerFunc(UpdateMemberRoleRoute)
	memberRouter.Path("").Methods(http.MethodDelete).HandlerFunc(RemoveMemberRoute)

	teamRouter.Path("/channels").Methods(http.MethodPost).HandlerFunc(CreateChannelRoute)
	teamRouter.Path("/channels").Methods(http.MethodGet).HandlerFunc(ListChannelsRoute)

	channelRouter := teamRouter.PathPrefix("/channels/{channelId}").Subrouter()
	channelRouter.Use(route_data.ChannelIdMiddleware)
	channelRouter.Path("").Methods(http.MethodGet).HandlerFunc(GetChannelRoute)
	channelRouter.Path("").Methods(http.MethodPut).HandlerFunc(UpdateChannelRoute)
	channelRouter.Path("").Methods(http.MethodDelete).HandlerFunc(DeleteChannelRoute)

	services.AttachRoutes(teamRouter)

	return nil
}
