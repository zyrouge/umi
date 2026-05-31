package route_data

import (
	"net/http"

	"github.com/gorilla/mux"
)

func TeamIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)["teamId"]
		r = r.WithContext(WithTeamId(r.Context(), v))
		next.ServeHTTP(w, r)
	})
}

func ChannelIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)["channelId"]
		r = r.WithContext(WithChannelId(r.Context(), v))
		next.ServeHTTP(w, r)
	})
}

func ServiceIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)["serviceId"]
		r = r.WithContext(WithServiceId(r.Context(), v))
		next.ServeHTTP(w, r)
	})
}

func MemberUserIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)["userId"]
		r = r.WithContext(WithMemberUserId(r.Context(), v))
		next.ServeHTTP(w, r)
	})
}
