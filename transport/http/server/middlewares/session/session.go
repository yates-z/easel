package session

import (
	"net/http"

	"github.com/yates-z/easel/auth/authentication/session"
	"github.com/yates-z/easel/transport/http/server"
)

const CookieName = "session_id"
const RedirectUrl = "/login"

func Middleware(sm *session.SessionManager) server.Middleware {

	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) error {
			if ctx.Request.RequestURI == RedirectUrl {
				return next(ctx)
			}
			// Get session_id from cookies
			sessionID, err := ctx.GetCookie(CookieName)
			if err != nil || sessionID == "" {
				ctx.String(http.StatusUnauthorized, "Unauthorized")
				return nil
			}

			// Load the session
			_session, err := sm.GetSession(sessionID)
			if err != nil {
				ctx.String(http.StatusUnauthorized, "Unauthorized")
				return nil
			}

			// Inject session into context
			ctx.Set("session", _session)
			return next(ctx)
		}
	}
}
