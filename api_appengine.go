// +build appengine

package box

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func getHTTPClient(ctx context.Context) *http.Client {
	return urlfetch.Client(ctx)
}
