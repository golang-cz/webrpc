// hello-webrpc v1.0.0 1769ce5a249c8ed4e4dab8320a7d67779eae0664
// --
// Code generated by webrpc-gen with golang generator. DO NOT EDIT.
//
// webrpc-gen -schema=hello-api.ridl -target=golang -pkg=main -server -out=./server/hello_api.gen.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const WebrpcHeader = "Webrpc"

const WebrpcHeaderValue = "webrpc;gen-golang@v0.16.0;hello-webrpc@v1.0.0"

// WebRPC description and code-gen version
func WebRPCVersion() string {
	return "v1"
}

// Schema version of your RIDL schema
func WebRPCSchemaVersion() string {
	return "v1.0.0"
}

// Schema hash generated from your RIDL schema
func WebRPCSchemaHash() string {
	return "1769ce5a249c8ed4e4dab8320a7d67779eae0664"
}

type WebrpcGenVersions struct {
	WebrpcGenVersion string
	CodeGenName      string
	CodeGenVersion   string
	SchemaName       string
	SchemaVersion    string
}

func VersionFromHeader(h http.Header) (*WebrpcGenVersions, error) {
	if h.Get(WebrpcHeader) == "" {
		return nil, fmt.Errorf("header is empty or missing")
	}

	versions, err := parseWebrpcGenVersions(h.Get(WebrpcHeader))
	if err != nil {
		return nil, fmt.Errorf("webrpc header is invalid: %w", err)
	}

	return versions, nil
}

func parseWebrpcGenVersions(header string) (*WebrpcGenVersions, error) {
	versions := strings.Split(header, ";")
	if len(versions) < 3 {
		return nil, fmt.Errorf("expected at least 3 parts while parsing webrpc header: %v", header)
	}

	_, webrpcGenVersion, ok := strings.Cut(versions[0], "@")
	if !ok {
		return nil, fmt.Errorf("webrpc gen version could not be parsed from: %s", versions[0])
	}

	tmplTarget, tmplVersion, ok := strings.Cut(versions[1], "@")
	if !ok {
		return nil, fmt.Errorf("tmplTarget and tmplVersion could not be parsed from: %s", versions[1])
	}

	schemaName, schemaVersion, ok := strings.Cut(versions[2], "@")
	if !ok {
		return nil, fmt.Errorf("schema name and schema version could not be parsed from: %s", versions[2])
	}

	return &WebrpcGenVersions{
		WebrpcGenVersion: webrpcGenVersion,
		CodeGenName:      tmplTarget,
		CodeGenVersion:   tmplVersion,
		SchemaName:       schemaName,
		SchemaVersion:    schemaVersion,
	}, nil
}

//
// Common types
//

type Kind uint32

const (
	Kind_USER  Kind = 1
	Kind_ADMIN Kind = 2
)

var Kind_name = map[uint32]string{
	1: "USER",
	2: "ADMIN",
}

var Kind_value = map[string]uint32{
	"USER":  1,
	"ADMIN": 2,
}

func (x Kind) String() string {
	return Kind_name[uint32(x)]
}

func (x Kind) MarshalText() ([]byte, error) {
	return []byte(Kind_name[uint32(x)]), nil
}

func (x *Kind) UnmarshalText(b []byte) error {
	*x = Kind(Kind_value[string(b)])
	return nil
}

func (x *Kind) Is(values ...Kind) bool {
	if x == nil {
		return false
	}
	for _, v := range values {
		if *x == v {
			return true
		}
	}
	return false
}

type Empty struct {
}

type User struct {
	ID        uint64     `json:"id" db:"id"`
	Username  string     `json:"USERNAME" db:"username"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:",omitempty"`
	DeletedAt *time.Time `json:"deletedAt"`
}

var (
	methods = map[string]method{
		"/rpc/ExampleService/Ping": {
			Name:        "Ping",
			Service:     "ExampleService",
			Annotations: map[string]string{},
		},
		"/rpc/ExampleService/GetUser": {
			Name:        "GetUser",
			Service:     "ExampleService",
			Annotations: map[string]string{},
		},
	}
)

var WebRPCServices = map[string][]string{
	"ExampleService": {
		"Ping",
		"GetUser",
	},
}

//
// Server types
//

type ExampleService interface {
	Ping(ctx context.Context) (bool, error)
	GetUser(ctx context.Context, userID uint64) (*User, error)
}

//
// Client types
//

type ExampleServiceClient interface {
	Ping(ctx context.Context) (bool, error)
	GetUser(ctx context.Context, userID uint64) (*User, error)
}

//
// Server
//

type WebRPCServer interface {
	http.Handler
}

type exampleServiceServer struct {
	ExampleService
	OnError   func(r *http.Request, rpcErr *WebRPCError)
	OnRequest func(w http.ResponseWriter, r *http.Request) error
}

func NewExampleServiceServer(svc ExampleService) *exampleServiceServer {
	return &exampleServiceServer{
		ExampleService: svc,
	}
}

func (s *exampleServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// In case of a panic, serve a HTTP 500 error and then panic.
		if rr := recover(); rr != nil {
			s.sendErrorJSON(w, r, ErrWebrpcServerPanic.WithCausef("%v", rr))
			panic(rr)
		}
	}()

	w.Header().Set(WebrpcHeader, WebrpcHeaderValue)

	ctx := r.Context()
	ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)
	ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)
	ctx = context.WithValue(ctx, ServiceNameCtxKey, "ExampleService")

	r = r.WithContext(ctx)

	var handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)
	switch r.URL.Path {
	case "/rpc/ExampleService/Ping":
		handler = s.servePingJSON
	case "/rpc/ExampleService/GetUser":
		handler = s.serveGetUserJSON
	default:
		err := ErrWebrpcBadRoute.WithCausef("no WebRPC method defined for path %v", r.URL.Path)
		s.sendErrorJSON(w, r, err)
		return
	}

	if r.Method != "POST" {
		w.Header().Add("Allow", "POST") // RFC 9110.
		err := ErrWebrpcBadMethod.WithCausef("unsupported method %v (only POST is allowed)", r.Method)
		s.sendErrorJSON(w, r, err)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if i := strings.Index(contentType, ";"); i >= 0 {
		contentType = contentType[:i]
	}
	contentType = strings.TrimSpace(strings.ToLower(contentType))

	switch contentType {
	case "application/json":
		if s.OnRequest != nil {
			if err := s.OnRequest(w, r); err != nil {
				rpcErr, ok := err.(WebRPCError)
				if !ok {
					rpcErr = ErrWebrpcEndpoint.WithCause(err)
				}
				s.sendErrorJSON(w, r, rpcErr)
				return
			}
		}

		handler(ctx, w, r)
	default:
		err := ErrWebrpcBadRequest.WithCausef("unsupported Content-Type %q (only application/json is allowed)", r.Header.Get("Content-Type"))
		s.sendErrorJSON(w, r, err)
	}
}

func (s *exampleServiceServer) servePingJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Ping")

	// Call service method implementation.
	ret0, err := s.ExampleService.Ping(ctx)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	respPayload := struct {
		Ret0 bool `json:"status"`
	}{ret0}
	respBody, err := json.Marshal(respPayload)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadResponse.WithCausef("failed to marshal json response: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) serveGetUserJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "GetUser")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCausef("failed to read request data: %w", err))
		return
	}
	defer r.Body.Close()

	reqPayload := struct {
		Arg0 uint64 `json:"userID"`
	}{}
	if err := json.Unmarshal(reqBody, &reqPayload); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCausef("failed to unmarshal request data: %w", err))
		return
	}

	// Call service method implementation.
	ret0, err := s.ExampleService.GetUser(ctx, reqPayload.Arg0)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	respPayload := struct {
		Ret0 *User `json:"user"`
	}{ret0}
	respBody, err := json.Marshal(respPayload)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadResponse.WithCausef("failed to marshal json response: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) sendErrorJSON(w http.ResponseWriter, r *http.Request, rpcErr WebRPCError) {
	if s.OnError != nil {
		s.OnError(r, &rpcErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rpcErr.HTTPStatus)

	respBody, _ := json.Marshal(rpcErr)
	w.Write(respBody)
}

func RespondWithError(w http.ResponseWriter, err error) {
	rpcErr, ok := err.(WebRPCError)
	if !ok {
		rpcErr = ErrWebrpcEndpoint.WithCause(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rpcErr.HTTPStatus)

	respBody, _ := json.Marshal(rpcErr)
	w.Write(respBody)
}

//
// Helpers
//

type method struct {
	Name        string
	Service     string
	Annotations map[string]string
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "webrpc context value " + k.name
}

var (
	HTTPResponseWriterCtxKey = &contextKey{"HTTPResponseWriter"}

	HTTPRequestCtxKey = &contextKey{"HTTPRequest"}

	ServiceNameCtxKey = &contextKey{"ServiceName"}

	MethodNameCtxKey = &contextKey{"MethodName"}
)

func ServiceNameFromContext(ctx context.Context) string {
	service, _ := ctx.Value(ServiceNameCtxKey).(string)
	return service
}

func MethodNameFromContext(ctx context.Context) string {
	method, _ := ctx.Value(MethodNameCtxKey).(string)
	return method
}

func RequestFromContext(ctx context.Context) *http.Request {
	r, _ := ctx.Value(HTTPRequestCtxKey).(*http.Request)
	return r
}

func MethodCtx(ctx context.Context) (method, bool) {
	req := RequestFromContext(ctx)
	if req == nil {
		return method{}, false
	}

	m, ok := methods[req.URL.Path]
	if !ok {
		return method{}, false
	}

	return m, true
}

func ResponseWriterFromContext(ctx context.Context) http.ResponseWriter {
	w, _ := ctx.Value(HTTPResponseWriterCtxKey).(http.ResponseWriter)
	return w
}

//
// Errors
//

type WebRPCError struct {
	Name       string `json:"error"`
	Code       int    `json:"code"`
	Message    string `json:"msg"`
	Cause      string `json:"cause,omitempty"`
	HTTPStatus int    `json:"status"`
	cause      error
}

var _ error = WebRPCError{}

func (e WebRPCError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s %d: %s: %v", e.Name, e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s %d: %s", e.Name, e.Code, e.Message)
}

func (e WebRPCError) Is(target error) bool {
	if target == nil {
		return false
	}
	if rpcErr, ok := target.(WebRPCError); ok {
		return rpcErr.Code == e.Code
	}
	return errors.Is(e.cause, target)
}

func (e WebRPCError) Unwrap() error {
	return e.cause
}

func (e WebRPCError) WithCause(cause error) WebRPCError {
	err := e
	err.cause = cause
	err.Cause = cause.Error()
	return err
}

func (e WebRPCError) WithCausef(format string, args ...interface{}) WebRPCError {
	cause := fmt.Errorf(format, args...)
	err := e
	err.cause = cause
	err.Cause = cause.Error()
	return err
}

// Deprecated: Use .WithCause() method on WebRPCError.
func ErrorWithCause(rpcErr WebRPCError, cause error) WebRPCError {
	return rpcErr.WithCause(cause)
}

// Webrpc errors
var (
	ErrWebrpcEndpoint           = WebRPCError{Code: 0, Name: "WebrpcEndpoint", Message: "endpoint error", HTTPStatus: 400}
	ErrWebrpcRequestFailed      = WebRPCError{Code: -1, Name: "WebrpcRequestFailed", Message: "request failed", HTTPStatus: 400}
	ErrWebrpcBadRoute           = WebRPCError{Code: -2, Name: "WebrpcBadRoute", Message: "bad route", HTTPStatus: 404}
	ErrWebrpcBadMethod          = WebRPCError{Code: -3, Name: "WebrpcBadMethod", Message: "bad method", HTTPStatus: 405}
	ErrWebrpcBadRequest         = WebRPCError{Code: -4, Name: "WebrpcBadRequest", Message: "bad request", HTTPStatus: 400}
	ErrWebrpcBadResponse        = WebRPCError{Code: -5, Name: "WebrpcBadResponse", Message: "bad response", HTTPStatus: 500}
	ErrWebrpcServerPanic        = WebRPCError{Code: -6, Name: "WebrpcServerPanic", Message: "server panic", HTTPStatus: 500}
	ErrWebrpcInternalError      = WebRPCError{Code: -7, Name: "WebrpcInternalError", Message: "internal error", HTTPStatus: 500}
	ErrWebrpcClientDisconnected = WebRPCError{Code: -8, Name: "WebrpcClientDisconnected", Message: "client disconnected", HTTPStatus: 400}
	ErrWebrpcStreamLost         = WebRPCError{Code: -9, Name: "WebrpcStreamLost", Message: "stream lost", HTTPStatus: 400}
	ErrWebrpcStreamFinished     = WebRPCError{Code: -10, Name: "WebrpcStreamFinished", Message: "stream finished", HTTPStatus: 200}
)

// Schema errors
var (
	ErrUserNotFound = WebRPCError{Code: 1000, Name: "UserNotFound", Message: "User not found", HTTPStatus: 400}
)
