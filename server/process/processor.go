package process

import (
	"context"
	"net/http"

	"github.com/danield21/danield-space/server/form"
	"github.com/gorilla/sessions"
)

type ProcessorFunc func(ctx context.Context, r *http.Request, s *sessions.Session) form.Form

func (prc ProcessorFunc) Process(ctx context.Context, r *http.Request, s *sessions.Session) form.Form {
	return prc(ctx, r, s)
}
