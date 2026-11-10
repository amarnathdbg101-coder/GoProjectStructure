// Package response provides standardized, uniform JSON response envelopes for Go REST APIs.
package response

import (
	"encoding/json"
	"net/http"
)

// Response defines the uniform JSON payload envelope.
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// JSON writes any data as JSON with the specified HTTP status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Success writes an HTTP 200 OK standard response envelope.
func Success(w http.ResponseWriter, message string, data any) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created writes an HTTP 201 Created standard response envelope.
func Created(w http.ResponseWriter, message string, data any) {
	JSON(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error writes a standardized error envelope with the provided status code.
func Error(w http.ResponseWriter, status int, errMessage string) {
	JSON(w, status, Response{
		Success: false,
		Error:   errMessage,
	})
}

// Message writes a simple message envelope without extra data payload.
func Message(w http.ResponseWriter, status int, message string) {
	isSuccess := status >= 200 && status < 300
	JSON(w, status, Response{
		Success: isSuccess,
		Message: message,
	})
}
