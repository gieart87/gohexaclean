package examples

import (
	"github.com/gieart87/gohexaclean/internal/dto/request"
	pkgErrors "github.com/gieart87/gohexaclean/pkg/errors"
	"github.com/gieart87/gohexaclean/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// This file contains examples of proper error handling using the error mapper
// DO NOT import this package in production code - it's for documentation only

// ExampleHandler demonstrates error handling patterns
type ExampleHandler struct {
	// userService inbound.UserServicePort
}

// ExampleLoginWithErrorMapper shows how to use error mapper in login handler
func (h *ExampleHandler) ExampleLoginWithErrorMapper(c *fiber.Ctx) error {
	var req request.LoginRequest

	// Step 1: Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid request body", err),
		)
	}

	// Step 2: Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewValidationErrorResponse(
				"Validation failed",
				response.ParseValidationErrors(err),
			),
		)
	}

	// Step 3: Call service
	// loginResp, err := h.userService.Login(c.Context(), &req)
	var err error // Mock for example
	if err != nil {
		// ✅ RECOMMENDED: Use error mapper
		// Automatically maps domain.ErrInvalidCredentials to 401 Unauthorized
		appErr := pkgErrors.MapDomainError(err)
		return c.Status(appErr.Code).JSON(
			response.NewErrorResponse(appErr.Message, err),
		)
	}

	return c.JSON(
		response.NewSuccessResponse("Login successful", nil),
	)
}

// ExampleGetUserWithErrorMapper shows how to use error mapper for GET endpoint
func (h *ExampleHandler) ExampleGetUserWithErrorMapper(c *fiber.Ctx) error {
	// Step 1: Parse and validate UUID
	_, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid user ID format", err),
		)
	}

	// Step 2: Call service
	// user, err := h.userService.GetUserByID(c.Context(), id)
	err = nil // Mock for example
	if err != nil {
		// ✅ RECOMMENDED: Use error mapper
		// Automatically maps domain.ErrUserNotFound to 404 Not Found
		appErr := pkgErrors.MapDomainError(err)
		return c.Status(appErr.Code).JSON(
			response.NewErrorResponse(appErr.Message, err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewSuccessResponse("User retrieved successfully", nil),
	)
}

// ExampleCreateUserWithErrorMapper shows error handling for POST endpoint
func (h *ExampleHandler) ExampleCreateUserWithErrorMapper(c *fiber.Ctx) error {
	var req request.CreateUserRequest

	// Step 1: Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse("Invalid request body", err),
		)
	}

	// Step 2: Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewValidationErrorResponse(
				"Validation failed",
				response.ParseValidationErrors(err),
			),
		)
	}

	// Step 3: Call service
	// user, err := h.userService.CreateUser(c.Context(), &req)
	var err error // Mock for example
	if err != nil {
		// ✅ RECOMMENDED: Use error mapper
		// Automatically maps:
		// - domain.ErrUserAlreadyExists to 409 Conflict
		// - other errors to appropriate status codes
		appErr := pkgErrors.MapDomainError(err)
		return c.Status(appErr.Code).JSON(
			response.NewErrorResponse(appErr.Message, err),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		response.NewSuccessResponse("User created successfully", nil),
	)
}

// ExampleWithCustomMessage shows how to use custom error messages
func (h *ExampleHandler) ExampleWithCustomMessage(c *fiber.Ctx) error {
	_, _ = uuid.Parse(c.Params("id"))

	// user, err := h.userService.GetUserByID(c.Context(), id)
	var err error // Mock for example
	if err != nil {
		// ✅ OPTION: Use custom message while keeping proper status code
		appErr := pkgErrors.MapDomainErrorWithCustomMessage(
			err,
			"The requested user could not be found in our system",
		)
		return c.Status(appErr.Code).JSON(
			response.NewErrorResponse(appErr.Message, err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewSuccessResponse("User retrieved successfully", nil),
	)
}

// ExampleGetStatusCodeOnly shows how to get status code without creating AppError
func (h *ExampleHandler) ExampleGetStatusCodeOnly(c *fiber.Ctx) error {
	_, _ = uuid.Parse(c.Params("id"))

	// user, err := h.userService.GetUserByID(c.Context(), id)
	var err error // Mock for example
	if err != nil {
		// ✅ OPTION: Get status code only
		statusCode := pkgErrors.GetHTTPStatusFromDomainError(err)
		return c.Status(statusCode).JSON(
			response.NewErrorResponse("Operation failed", err),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewSuccessResponse("User retrieved successfully", nil),
	)
}

// ExampleManualErrorHandling shows manual error handling (not recommended)
func (h *ExampleHandler) ExampleManualErrorHandling(c *fiber.Ctx) error {
	// var user, err := h.userService.GetUserByID(c.Context(), id)
	var err error // Mock for example

	if err != nil {
		// ⚠️ LEGACY APPROACH: Manual error checking
		// Less maintainable, more verbose, easy to miss cases
		/*
			if errors.Is(err, domain.ErrUserNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(
					response.NewErrorResponse("User not found", err),
				)
			}
			if errors.Is(err, domain.ErrUserAlreadyExists) {
				return c.Status(fiber.StatusConflict).JSON(
					response.NewErrorResponse("User already exists", err),
				)
			}
			if errors.Is(err, domain.ErrInvalidCredentials) {
				return c.Status(fiber.StatusUnauthorized).JSON(
					response.NewErrorResponse("Invalid credentials", err),
				)
			}
			// ... many more cases
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.NewErrorResponse("Internal server error", err),
			)
		*/
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewSuccessResponse("Success", nil),
	)
}

// Summary:
//
// ✅ RECOMMENDED PATTERNS:
// 1. Use pkgErrors.MapDomainError(err) for automatic mapping
// 2. Use pkgErrors.MapDomainErrorWithCustomMessage() for custom messages
// 3. Use pkgErrors.GetHTTPStatusFromDomainError() for status code only
//
// ⚠️ AVOID:
// - Manual error checking with multiple if statements
// - Hardcoding HTTP status codes in service layer
// - Comparing error strings instead of using errors.Is()
//
// 💡 BENEFITS:
// - Centralized error mapping logic
// - Consistent error responses across all endpoints
// - Easy to add new error types
// - Maintainable and testable code
