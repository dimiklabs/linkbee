package dto

import "net/http"

type ServiceError struct {
	ErrorCode   string                 `json:"error_code"`
	Description string                 `json:"description"`
	StatusCode  int                    `json:"-"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

func (e *ServiceError) Error() string {
	return e.Description
}

func (e *ServiceError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

func (e *ServiceError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

func (e *ServiceError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

func (e *ServiceError) IsBadRequest() bool {
	return e.StatusCode == http.StatusBadRequest
}

func (e *ServiceError) IsConflict() bool {
	return e.StatusCode == http.StatusConflict
}

func (e *ServiceError) IsInternalError() bool {
	return e.StatusCode == http.StatusInternalServerError
}

func NewServiceError(errorCode, description string, statusCode int) *ServiceError {
	return &ServiceError{
		ErrorCode:   errorCode,
		Description: description,
		StatusCode:  statusCode,
	}
}

func NewServiceErrorWithData(errorCode, description string, statusCode int, data map[string]interface{}) *ServiceError {
	return &ServiceError{
		ErrorCode:   errorCode,
		Description: description,
		StatusCode:  statusCode,
		Data:        data,
	}
}

func NewInternalError(errorCode, description string) *ServiceError {
	return NewServiceError(errorCode, description, http.StatusInternalServerError)
}

func NewBadRequestError(errorCode, description string) *ServiceError {
	return NewServiceError(errorCode, description, http.StatusBadRequest)
}

func NewUnauthorizedError(errorCode, description string) *ServiceError {
	return NewServiceError(errorCode, description, http.StatusUnauthorized)
}

func NewForbiddenError(errorCode, description string) *ServiceError {
	return NewServiceError(errorCode, description, http.StatusForbidden)
}

func NewNotFoundError(errorCode, description string) *ServiceError {
	return NewServiceError(errorCode, description, http.StatusNotFound)
}

func NewConflictError(errorCode, description string) *ServiceError {
	return NewServiceError(errorCode, description, http.StatusConflict)
}

func NewTooManyRequestsError(errorCode, description string) *ServiceError {
	return NewServiceError(errorCode, description, http.StatusTooManyRequests)
}
