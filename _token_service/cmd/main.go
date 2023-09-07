package tokenservice

import (
	"senpainikolay/go-internship-smartdata/auth-service/_token_service/controller"
)

func RunTokenValidationService() {
	controller.Serve(":4444")
}
