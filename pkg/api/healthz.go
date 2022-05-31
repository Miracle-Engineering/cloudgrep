package api

import (
	"github.com/gin-gonic/gin"
	"github.com/run-x/cloudgrep/pkg/datastore"
)

func Healthz(c *gin.Context) {
	ds := c.MustGet("datastore").(datastore.Datastore)
	err := ds.Ping()
	if err != nil {
		badRequest(c, "Internal")
		return
	}
	successResponse(c, gin.H{"status": "All good!"})
}
