// +build !appengine

package box

import (
	"net/http"

	"golang.org/x/net/context"
)

func getHTTPClient(ctx context.Context) *http.Client {
	return http.DefaultClient
}
