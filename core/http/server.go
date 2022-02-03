package graph_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	graph_database "github.com/jeconias/graph-service/core/database"
	"github.com/jeconias/graph-service/core/graph"
)

type GraphHttp struct {
	database *graph_database.Database
	gin      *gin.Engine
}

func (v *GraphHttp) Init(db *graph_database.Database) {

	r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	v.gin = r
	v.database = db

	v.initializeRoutes()
}

func (v *GraphHttp) Run(address string) error {
	return v.gin.Run(address)
}

func (v *GraphHttp) initializeRoutes() {

	v.gin.GET("/graph-raw", func(ginCtx *gin.Context) {
		graph := graph.Graph{}
		graph.Init()

		result, err := v.database.ListVertice()
		if err != nil {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"data": result,
		})
	})

	v.gin.GET("/", func(ginCtx *gin.Context) {
		ginCtx.HTML(http.StatusOK, "index.html", nil)
	})

}
