package index

import (
	rh "github.com/kontza/goatmospi/route_handler"
	"github.com/gin-gonic/gin"
	"github.com/kontza/goatmospi/controller/data"
	"net/http"
)

func init() {
	rh.AddRouteHandler(rh.RoutingData{"GET", "/", indexHandler})
}

// Build the main index.html.
func indexHandler(ctx *gin.Context) {
	prefixPath := ctx.Request.Header.Get("Atmospi-Prefix-Path")
	oldest, newest := data.GetTimestampRange()
	ctx.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"SubDomain":       prefixPath,
			"OldestTimestamp": oldest,
			"LatestTimestamp": newest,
		},
	)
}
