package admin

import (
	"net/http"

	"github.com/danield21/danield-space/server/controllers/action"
	"github.com/danield21/danield-space/server/controllers/link"
	"github.com/danield21/danield-space/server/controllers/status"
	"github.com/danield21/danield-space/server/controllers/view"
	"github.com/danield21/danield-space/server/handler"
	"github.com/danield21/danield-space/server/repository/account"
	"github.com/danield21/danield-space/server/repository/siteInfo"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

var AccountCreateHeadersHandler = view.HeaderHandler(http.StatusOK,
	view.Header{"Content-Type", view.HTMLContentType},
)

var AccountCreatePageHandler = handler.Chain(
	view.HTMLHandler,
	handler.ToLink(handler.Chain(
		AccountCreateHeadersHandler,
		AccountCreatePageLink,
		status.LinkAll,
	)),
)

var AccountCreateActionHandler = handler.Chain(
	view.HTMLHandler,
	handler.ToLink(handler.Chain(
		AccountCreateHeadersHandler,
		AccountCreatePageLink,
		action.PutAccountLink,
		status.LinkAll,
	)),
)

func AccountCreatePageLink(h handler.Handler) handler.Handler {
	return func(ctx context.Context, e handler.Environment, w http.ResponseWriter) (context.Context, error) {
		ses := handler.Session(ctx)
		req := handler.Request(ctx)
		frm := action.Form(ctx)

		user, signedIn := link.User(ses)
		if !signedIn {
			return ctx, status.ErrUnauthorized
		}

		current, err := account.Get(ctx, user)
		if err != nil {
			log.Warningf(ctx, "AccountCreatePageLink - Unable to verify account %s\n%v", user, err)
			return ctx, status.ErrUnauthorized
		}

		target := req.Form.Get("account")
		if frm.IsEmpty() && target != "" {
			tUser, err := account.Get(ctx, target)
			if err == nil {
				frm = action.AccountToForm(tUser)
			} else {
				log.Warningf(ctx, "Unable to get account %s\n%v", target, err)
			}
		}

		info := siteInfo.Get(ctx)

		data := struct {
			AdminModel
			action.Result
			Super bool
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
			Super: current.Super,
		}

		return h(link.PageContext(ctx, "page/admin/account-create", data), e, w)
	}
}
