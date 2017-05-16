package admin

import (
	"net/http"

	"github.com/danield21/danield-space/server/controllers/action"
	"github.com/danield21/danield-space/server/controllers/link"
	"github.com/danield21/danield-space/server/controllers/status"
	"github.com/danield21/danield-space/server/controllers/view"
	"github.com/danield21/danield-space/server/handler"
	"golang.org/x/net/context"
)

var SiteInfoHeadersHandler = view.HeaderHandler(http.StatusOK,
	view.Header{"Content-Type", view.HTMLContentType},
)

var SiteInfoPageHandler = handler.Chain(
	view.HTMLHandler,
	handler.ToLink(handler.Chain(
		SiteInfoHeadersHandler,
		SiteInfoPageLink,
		status.LinkAll,
	)),
)

var SiteInfoActionHandler = handler.Chain(
	view.HTMLHandler,
	handler.ToLink(handler.Chain(
		SiteInfoHeadersHandler,
		SiteInfoPageLink,
		action.PutSiteInfoLink,
		status.LinkAll,
	)),
)

func SiteInfoPageLink(h handler.Handler) handler.Handler {
	return func(ctx context.Context, e handler.Environment, w http.ResponseWriter) (context.Context, error) {
		frm := action.Form(ctx)
		ses := handler.Session(ctx)

		user, signedIn := link.User(ses)
		if !signedIn {
			return ctx, status.ErrUnauthorized
		}

		info := e.Repository().SiteInfo().Get(ctx)

		if frm.IsEmpty() {
			frm = action.RepackSiteInfo(info)
		}

		data := struct {
			AdminModel
			action.Result
		}{
			AdminModel: AdminModel{
				BaseModel: view.BaseModel{
					SiteInfo: info,
				},
				User: user,
			},
			Result: action.Result{
				Form: frm,
			},
		}

		return h(link.PageContext(ctx, "page/admin/site-info-manage", data), e, w)
	}
}
