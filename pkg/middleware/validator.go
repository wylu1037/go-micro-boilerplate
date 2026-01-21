package middleware

import (
	"context"

	"buf.build/go/protovalidate"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func validate(ctx context.Context, reqOrRes any, protoValidate ProtoValidateFunc, onValidationErrCallback OnValidationErrCallback) (err error) {
	message, ok := reqOrRes.(proto.Message)
	if !ok {
		return nil
	}
	err = protoValidate(message)

	if err == nil {
		return nil
	}

	if onValidationErrCallback != nil {
		onValidationErrCallback(ctx, err)
	}
	return status.Error(codes.InvalidArgument, err.Error())
}

type options struct {
	protoValidate           ProtoValidateFunc
	onValidationErrCallback OnValidationErrCallback
}
type Option func(*options)

func evaluateOpts(opts []Option) *options {
	optCopy := &options{}
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type OnValidationErrCallback func(ctx context.Context, err error)
type ProtoValidateFunc func(message proto.Message) error

// WithOnValidationErrCallback registers function that will be invoked on validation error(s).
func WithOnValidationErrCallback(onValidationErrCallback OnValidationErrCallback) Option {
	return func(o *options) {
		o.onValidationErrCallback = onValidationErrCallback
	}
}

// WithProtoValidate validate proto
func WithProtoValidate(v protovalidate.Validator) Option {
	return func(o *options) {
		o.protoValidate = func(msg proto.Message) error {
			return v.Validate(msg)
		}
	}
}

// NewValidatorMiddleware returns a go-micro server.HandlerWrapper that validates incoming messages.
func NewValidatorMiddleware(logger *zap.Logger) server.HandlerWrapper {
	logErr := func(ctx context.Context, err error) {
		logger.Error("middleware: failed to validate", zap.Error(err))
	}
	goValidator, err := protovalidate.New()
	if err != nil {
		logger.Error("middleware: failed to new protovalidate", zap.Error(err))
	}

	opts := []Option{
		WithProtoValidate(goValidator),
		WithOnValidationErrCallback(logErr),
	}
	o := evaluateOpts(opts)

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp any) error {
			// Validate the request body (which is the message)
			if err := validate(ctx, req.Body(), o.protoValidate, o.onValidationErrCallback); err != nil {
				return err
			}
			return fn(ctx, req, rsp)
		}
	}
}
