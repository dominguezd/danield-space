package link

import (
	"net/http"

	"github.com/danield21/danield-space/server/envir"
	"github.com/danield21/danield-space/server/repository/theme"
	"github.com/danield21/danield-space/server/service"
	"github.com/danield21/danield-space/server/service/view"
	"golang.org/x/net/context"
)

func Theme(h service.Handler) service.Handler {
	return func(ctx context.Context, e envir.Environment, w http.ResponseWriter) (context.Context, error) {
		r := service.Request(ctx)

		defaultTheme := theme.GetApp(ctx)

		theme := getRequestedTheme(r, defaultTheme)

		return h(view.WithTheme(ctx, theme), e, w)
	}
}

//GetTheme gets the theme. If no theme was specified, then the default theme is given
func getRequestedTheme(r *http.Request, defaultTheme string) string {
	err := r.ParseForm()
	if err != nil {
		return defaultTheme
	}

	if theme, ok := r.Form["theme"]; ok {
		return theme[0]
	}

	return defaultTheme
}