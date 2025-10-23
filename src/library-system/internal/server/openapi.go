package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

var swaggerHtml = []byte(`<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui.css" >
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist/swagger-ui-bundle.js"></script>
<script>
  const ui = SwaggerUIBundle({
    url: "/docs/swagger.json",
    dom_id: '#swagger-ui',
    presets: [SwaggerUIBundle.presets.apis],
    layout: "BaseLayout"
  })
</script>
</body>
</html>`)

func (s *Server) InitDocsRoutes() error {
	yamlData, err := ioutil.ReadFile("/home/sanchiko/lab2-rsoi/v4/[inst][v4] Library System.yml")
	if err != nil {
		return err
	}

	var spec interface{}
	if err := yaml.Unmarshal(yamlData, &spec); err != nil {
		return err
	}
	swaggerJSON, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	api := s.GinRouter.Group("/docs")
	api.GET("/swagger.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", swaggerJSON)
	})

	api.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", swaggerHtml)
	})

	return nil
}
