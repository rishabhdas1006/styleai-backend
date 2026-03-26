package common

import "errors"

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailNotFound      = errors.New("email not found")
	ErrMissingAuthToken   = errors.New("missing auth token")
	ErrInvalidAuthToken   = errors.New("invalid auth token")
	ErrCategoryExists     = errors.New("category already exists")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrVariantExists      = errors.New("variant already exists")
	ErrVariantNotFound    = errors.New("variant not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrForbidden          = errors.New("forbidden")

	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrCartNotFound      = errors.New("cart not found")
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrUnauthorized      = errors.New("unauthorized action")
)
