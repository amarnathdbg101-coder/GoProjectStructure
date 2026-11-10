package handler

import (
	"encoding/json"
	"net/http"
)

// Response defines the uniform API JSON envelope.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON writes a JSON response with status code and headers.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Success writes a 200 OK JSON success envelope.
func Success(w http.ResponseWriter, msg string, data interface{}) {
	JSON(w, http.StatusOK, Response{Success: true, Message: msg, Data: data})
}

// Created writes a 201 Created JSON envelope.
func Created(w http.ResponseWriter, msg string, data interface{}) {
	JSON(w, http.StatusCreated, Response{Success: true, Message: msg, Data: data})
}

// Error writes a standardized JSON error response.
func Error(w http.ResponseWriter, status int, err string) {
	JSON(w, status, Response{Success: false, Error: err})
}

// Message writes a simple message response.
func Message(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, Response{Success: true, Message: msg})
}
