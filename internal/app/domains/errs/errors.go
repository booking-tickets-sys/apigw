package errs

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HTTPError represents a structured HTTP error response
type HTTPError struct {
	ErrorType string `json:"error"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Status    int    `json:"-"`
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return e.Message
}

// NewHTTPError creates a new HTTP error
func NewHTTPError(errorType, code, message string, status int) *HTTPError {
	return &HTTPError{
		ErrorType: errorType,
		Code:      code,
		Message:   message,
		Status:    status,
	}
}

// Common HTTP errors
var (
	ErrBadRequest         = NewHTTPError("VALIDATION_ERROR", "BAD_REQUEST", "Invalid request", http.StatusBadRequest)
	ErrUnauthorized       = NewHTTPError("AUTHENTICATION_ERROR", "UNAUTHORIZED", "Authentication required", http.StatusUnauthorized)
	ErrForbidden          = NewHTTPError("AUTHORIZATION_ERROR", "FORBIDDEN", "Access denied", http.StatusForbidden)
	ErrNotFound           = NewHTTPError("NOT_FOUND_ERROR", "RESOURCE_NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrConflict           = NewHTTPError("CONFLICT_ERROR", "RESOURCE_CONFLICT", "Resource conflict", http.StatusConflict)
	ErrInternalServer     = NewHTTPError("INTERNAL_ERROR", "INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrServiceUnavailable = NewHTTPError("SERVICE_ERROR", "SERVICE_UNAVAILABLE", "Service temporarily unavailable", http.StatusServiceUnavailable)
)

// GRPCToHTTPError converts a gRPC error to an appropriate HTTP error
func GRPCToHTTPError(err error) *HTTPError {
	if err == nil {
		return nil
	}

	// Check if it's a gRPC status error
	st, ok := status.FromError(err)
	if !ok {
		// If it's not a gRPC status error, return internal server error
		return ErrInternalServer
	}

	// Map gRPC codes to HTTP errors
	switch st.Code() {
	case codes.OK:
		return nil
	case codes.InvalidArgument:
		return NewHTTPError("VALIDATION_ERROR", "INVALID_ARGUMENT", st.Message(), http.StatusBadRequest)
	case codes.NotFound:
		return NewHTTPError("NOT_FOUND_ERROR", "RESOURCE_NOT_FOUND", st.Message(), http.StatusNotFound)
	case codes.AlreadyExists:
		return NewHTTPError("CONFLICT_ERROR", "RESOURCE_ALREADY_EXISTS", st.Message(), http.StatusConflict)
	case codes.PermissionDenied:
		return NewHTTPError("AUTHORIZATION_ERROR", "PERMISSION_DENIED", st.Message(), http.StatusForbidden)
	case codes.Unauthenticated:
		return NewHTTPError("AUTHENTICATION_ERROR", "UNAUTHENTICATED", st.Message(), http.StatusUnauthorized)
	case codes.ResourceExhausted:
		return NewHTTPError("RATE_LIMIT_ERROR", "RATE_LIMIT_EXCEEDED", st.Message(), http.StatusTooManyRequests)
	case codes.FailedPrecondition:
		return NewHTTPError("PRECONDITION_ERROR", "PRECONDITION_FAILED", st.Message(), http.StatusBadRequest)
	case codes.Aborted:
		return NewHTTPError("CONFLICT_ERROR", "OPERATION_ABORTED", st.Message(), http.StatusConflict)
	case codes.OutOfRange:
		return NewHTTPError("VALIDATION_ERROR", "OUT_OF_RANGE", st.Message(), http.StatusBadRequest)
	case codes.Unimplemented:
		return NewHTTPError("NOT_IMPLEMENTED_ERROR", "METHOD_NOT_IMPLEMENTED", st.Message(), http.StatusNotImplemented)
	case codes.Internal:
		return NewHTTPError("INTERNAL_ERROR", "INTERNAL_SERVER_ERROR", st.Message(), http.StatusInternalServerError)
	case codes.Unavailable:
		return NewHTTPError("SERVICE_ERROR", "SERVICE_UNAVAILABLE", st.Message(), http.StatusServiceUnavailable)
	case codes.DataLoss:
		return NewHTTPError("DATA_ERROR", "DATA_LOSS", st.Message(), http.StatusInternalServerError)
	case codes.DeadlineExceeded:
		return NewHTTPError("TIMEOUT_ERROR", "REQUEST_TIMEOUT", st.Message(), http.StatusRequestTimeout)
	case codes.Canceled:
		return NewHTTPError("CANCELED_ERROR", "REQUEST_CANCELED", st.Message(), http.StatusRequestTimeout)
	default:
		return NewHTTPError("UNKNOWN_ERROR", "UNKNOWN_ERROR", st.Message(), http.StatusInternalServerError)
	}
}

// GetGRPCCode returns the gRPC status code from an error
func GetGRPCCode(err error) codes.Code {
	if st, ok := status.FromError(err); ok {
		return st.Code()
	}
	return codes.Unknown
}
