package swagger

import (
	"reviewsch/docs"
)

func SetupSwagger() {
	docs.SwaggerInfo.Title = "Coupon Service API"
	docs.SwaggerInfo.Description = "This is a coupon service server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
