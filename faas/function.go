package faas

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"time"

	faascontext "github.com/serverless-plus/tencent-serverless-go/context"
	"github.com/serverless-plus/tencent-serverless-go/faas/messages"
)

// Function funciton
type Function struct {
	handler Handler
}

// Ping ping request
func (fn *Function) Ping(req *messages.PingRequest, response *messages.PingResponse) error {
	*response = messages.PingResponse{}
	return nil
}

// Invoke invoke function
func (fn *Function) Invoke(req *messages.InvokeRequest, response *messages.InvokeResponse) error {
	defer func() {
		if err := recover(); err != nil {
			panicInfo := getPanicInfo(err)
			response.Error = &messages.InvokeResponse_Error{
				Message:    panicInfo.Message,
				Type:       getErrorType(err),
				StackTrace: panicInfo.StackTrace,
				ShouldExit: true,
			}
		}
	}()

	deadline := time.Unix(req.Deadline.Seconds, req.Deadline.Nanos).UTC()
	invokeContext, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	lc := &faascontext.FunctionContext{
		RequestID:             req.RequestId,
		InvokedFunctionUnique: req.InvokedFunctionUnique,
		Namespace:             req.Namespace,
		FunctionName:          req.FunctionName,
		FunctionVersion:       req.FunctionVersion,
		MemoryLimitInMb:       req.MemoryLimitInMb,
		TimeLimitInMs:         req.TimeLimitInMs,
	}

	if len(req.ClientContext) > 0 {
		if err := json.Unmarshal(req.ClientContext, &lc.ClientContext); err != nil {
			response.Error = functionErrorResponse(err)
			return nil
		}
	}

	if len(req.Environment) > 0 {
		if err := json.Unmarshal([]byte(req.Environment), &lc.Environment); err != nil {
			response.Error = functionErrorResponse(err)
			return nil
		}
		for key, value := range lc.Environment {
			os.Setenv(key, value)
		}
	}

	invokeContext = faascontext.NewContext(invokeContext, lc)

	payload, err := fn.handler.Invoke(invokeContext, req.Payload)
	if err != nil {
		response.Error = functionErrorResponse(err)
		return nil
	}

	response.Payload = payload

	return nil
}

func getErrorType(err interface{}) string {
	errorType := reflect.TypeOf(err)
	if errorType.Kind() == reflect.Ptr {
		return errorType.Elem().Name()
	}
	return errorType.Name()
}

func functionErrorResponse(invokeError error) *messages.InvokeResponse_Error {
	var errorName string
	if errorType := reflect.TypeOf(invokeError); errorType.Kind() == reflect.Ptr {
		errorName = errorType.Elem().Name()
	} else {
		errorName = errorType.Name()
	}
	return &messages.InvokeResponse_Error{
		Message: invokeError.Error(),
		Type:    errorName,
	}
}
