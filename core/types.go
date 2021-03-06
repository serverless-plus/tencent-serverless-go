package core

import (
	"fmt"
	"net/http"

	"github.com/serverless-plus/tencent-serverless-go/events"
)

// GatewayTimeout returns a dafault Gateway Timeout (504) response
func GatewayTimeout() events.APIGatewayResponse {
	return events.APIGatewayResponse{StatusCode: http.StatusGatewayTimeout}
}

// NewLoggedError generates a new error and logs it to stdout
func NewLoggedError(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Println(err.Error())
	return err
}
