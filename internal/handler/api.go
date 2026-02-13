package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	v1 "github.com/soltiHQ/control-plane/api/v1"
	"github.com/soltiHQ/control-plane/domain/kind"
	"github.com/soltiHQ/control-plane/internal/service/access"
	"github.com/soltiHQ/control-plane/internal/service/session"
	"github.com/soltiHQ/control-plane/internal/service/user"
	"github.com/soltiHQ/control-plane/internal/storage"
	"github.com/soltiHQ/control-plane/internal/storage/inmemory"
	"github.com/soltiHQ/control-plane/internal/transport/http/apimap"
	"github.com/soltiHQ/control-plane/internal/transport/http/responder"
	"github.com/soltiHQ/control-plane/internal/transport/http/response"
	"github.com/soltiHQ/control-plane/internal/transport/http/route"
	contentUser "github.com/soltiHQ/control-plane/ui/templates/content/user"
)

// API handlers.
type API struct {
	logger     zerolog.Logger
	accessSVC  *access.Service
	userSVC    *user.Service
	sessionSVC *session.Service
}

// NewAPI creates a new API handler.
func NewAPI(
	logger zerolog.Logger,
	accessSVC *access.Service,
	userSVC *user.Service,
	sessionSVC *session.Service,
) *API {
	if accessSVC == nil {
		panic("handler.API: accessSVC is nil")
	}
	if userSVC == nil {
		panic("handler.API: userSVC is nil")
	}
	if sessionSVC == nil {
		panic("handler.API: sessionSVC is nil")
	}
	return &API{
		logger:     logger.With().Str("handler", "api").Logger(),
		accessSVC:  accessSVC,
		userSVC:    userSVC,
		sessionSVC: sessionSVC,
	}
}

// Routes registers API routes.
func (a *API) Routes(mux *http.ServeMux, auth route.BaseMW, perm route.PermMW, common ...route.BaseMW) {
	route.HandleFunc(mux, "/api/v1/users", a.UsersList, append(common, auth, perm(kind.UsersGet))...)
	route.HandleFunc(mux, "/api/v1/users/detail/", a.UsersDetail, append(common, auth, perm(kind.UsersGet))...)
	route.HandleFunc(mux, "/api/v1/users/session/", a.UsersSession, append(common, auth, perm(kind.UsersGet))...)
}

// UsersList handles GET /api/v1/users
//
// Query params:
//   - limit: int (optional)
//   - cursor: string (optional)
//   - q: string (optional, currently only used for HTML template; filtering should be done via storage filter factory)
func (a *API) UsersList(w http.ResponseWriter, r *http.Request) {
	mode := response.ModeFromRequest(r)

	if r.Method != http.MethodGet {
		response.NotAllowed(w, r, mode)
		return
	}
	if r.URL.Path != "/api/v1/users" {
		response.NotFound(w, r, mode)
		return
	}

	var (
		limit  int
		cursor = r.URL.Query().Get("cursor")
		q      = r.URL.Query().Get("q")
		filter storage.UserFilter
	)
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}
	if q != "" {
		filter = inmemory.NewUserFilter().Query(q)
	}

	res, err := a.userSVC.List(r.Context(), user.ListQuery{
		Limit:  limit,
		Cursor: cursor,
		Filter: filter,
	})
	if err != nil {
		a.logger.Error().Err(err).Msg("api: list users failed")
		response.Unavailable(w, r, mode)
		return
	}

	items := make([]v1.User, 0, len(res.Items))
	for _, u := range res.Items {
		if u == nil {
			continue
		}
		items = append(items, apimap.User(u))
	}
	response.OK(w, r, mode, &responder.View{
		Data: v1.UserListResponse{
			Items:      items,
			NextCursor: res.NextCursor,
		},
		Component: contentUser.List(res.Items, res.NextCursor, q),
	})
}

// UsersDetail handles GET /api/v1/users/detail/{id}
func (a *API) UsersDetail(w http.ResponseWriter, r *http.Request) {
	mode := response.ModeFromRequest(r)

	if r.Method != http.MethodGet {
		response.NotAllowed(w, r, mode)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/api/v1/users/detail/") {
		response.NotFound(w, r, mode)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/users/detail/")
	if id == "" {
		response.NotFound(w, r, mode)
		return
	}

	u, err := a.userSVC.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			response.NotFound(w, r, mode)
			return
		}
		a.logger.Error().Err(err).Str("user_id", id).Msg("api: get user failed")
		response.Unavailable(w, r, mode)
		return
	}

	apiUser := apimap.User(u)
	response.OK(w, r, mode, &responder.View{
		Data:      apiUser,
		Component: contentUser.Detail(apiUser),
	})
}

// UsersSession handles GET /api/v1/users/session/{id}
//
// Query params:
//   - limit: int (optional)
func (a *API) UsersSession(w http.ResponseWriter, r *http.Request) {
	mode := response.ModeFromRequest(r)

	if r.Method != http.MethodGet {
		response.NotAllowed(w, r, mode)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/api/v1/users/session/") {
		response.NotFound(w, r, mode)
		return
	}

	userID := strings.TrimPrefix(r.URL.Path, "/api/v1/users/session/")
	if userID == "" {
		response.NotFound(w, r, mode)
		return
	}

	var limit int
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}

	res, err := a.sessionSVC.ListByUser(r.Context(), session.ListByUserQuery{
		UserID: userID,
		Limit:  limit,
	})
	if err != nil {
		a.logger.Error().Err(err).Str("user_id", userID).Msg("api: list sessions failed")
		response.Unavailable(w, r, mode)
		return
	}

	items := make([]v1.Session, 0, len(res.Items))
	for _, s := range res.Items {
		if s == nil {
			continue
		}
		items = append(items, apimap.Session(s))
	}
	response.OK(w, r, mode, &responder.View{
		Data: v1.SessionResponse{
			Items: items,
		},
		Component: contentUser.Sessions(items),
	})
}
