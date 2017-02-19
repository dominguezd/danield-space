package view

import (
	"errors"
	"net/http"

	"github.com/danield21/danield-space/server/envir"
	"golang.org/x/net/context"
)

var ErrNoPage = errors.New("context does not page")
var ErrNoTheme = errors.New("context does not theme")

type uniqueKey string

const pageKey = uniqueKey("page")
const themeKey = uniqueKey("theme")
const dataKey = uniqueKey("data")

func HTMLHandler(ctx context.Context, e envir.Environment, w http.ResponseWriter) error {
	page, err := Page(ctx)
	if err != nil {
		return err
	}
	theme, err := Theme(ctx)
	if err != nil {
		return err
	}
	data := Data(ctx)

	return e.View(w, theme, page, data)
}

func WithPage(ctx context.Context, page string) {
	context.WithValue(ctx, pageKey, page)
}

func WithTheme(ctx context.Context, theme string) {
	context.WithValue(ctx, themeKey, theme)
}

func WithData(ctx context.Context, data interface{}) {
	context.WithValue(ctx, dataKey, data)
}

func Page(ctx context.Context) (string, error) {
	iPage := ctx.Value(pageKey)
	page, ok := iPage.(string)
	if !ok {
		return "", ErrNoPage
	}
	return page, nil
}

func Theme(ctx context.Context) (string, error) {
	iTheme := ctx.Value(themeKey)
	theme, ok := iTheme.(string)
	if !ok {
		return "", ErrNoTheme
	}
	return theme, nil
}

func Data(ctx context.Context) interface{} {
	return ctx.Value(dataKey)
}
