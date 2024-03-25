package middleware

import (
	"embed"
	"net/http"
)

type AssetHandlerArgs func(http.ResponseWriter, *http.Request, embed.FS)

func Assets(next AssetHandlerArgs, assetsFs embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r, assetsFs)
	}
}
