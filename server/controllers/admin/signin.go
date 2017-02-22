package admin

import (
	"net/http"

	"github.com/danield21/danield-space/server/controllers/action"
	"github.com/danield21/danield-space/server/controllers/link"
	"github.com/danield21/danield-space/server/controllers/status"
	"github.com/danield21/danield-space/server/handler"
	"github.com/danield21/danield-space/server/handler/form"
	"github.com/danield21/danield-space/server/handler/view"
	"github.com/danield21/danield-space/server/repository/siteInfo"
	"golang.org/x/net/context"
)

var SignInHeadersHandler = view.HeaderHandler(http.StatusOK,
	view.Header{"Content-Type", view.HTMLContentType},
)

var SignInPageHandler = handler.Chain(
	view.HTMLHandler,
	handler.ToLink(handler.Chain(
		SignInHeadersHandler,
		SignInPageLink,
		link.Theme,
		status.LinkAll,
	)),
)

var SignInActionHandler = handler.Chain(
	view.HTMLHandler,
	handler.ToLink(handler.Chain(
		SignInHeadersHandler,
		SignInPageLink,
		link.SaveSession,
		action.AuthenicateLink,
		link.Theme,
		status.LinkAll,
	)),
)

func SignInPageLink(h handler.Handler) handler.Handler {
	return func(ctx context.Context, e handler.Environment, w http.ResponseWriter) (context.Context, error) {
		r := handler.Request(ctx)
		info := siteInfo.Get(ctx)
		f := form.AsForm(ctx)

		data := struct {
			handler.BaseModel
			Redirect string
			Form     form.Form
		}{
			BaseModel: handler.BaseModel{
				SiteInfo: info,
			},
			Form:     f,
			Redirect: action.Redirect(r),
		}

		return h(link.PageContext(ctx, "page/admin/signin", data), e, w)
	}
}
