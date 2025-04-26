// api-gateway/internal/handler/user_handler.go
package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "github.com/gorilla/mux"
    "github.com/sibobbbbbb/backend-engineer-challenge/api-gateway/internal/client"
    "github.com/sibobbbbbb/backend-engineer-challenge/api-gateway/internal/middleware"
    "github.com/sibobbbbbb/backend-engineer-challenge/api-gateway/pkg/response"
)

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Fullname string `json:"fullname"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type UpdateUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Fullname string `json:"fullname"`
}

type UserHandler struct {
    userClient     client.UserClient
    authMiddleware *middleware.AuthMiddleware
}

func NewUserHandler(r *mux.Router, uc client.UserClient, am *middleware.AuthMiddleware) {
    handler := &UserHandler{
        userClient:     uc,
        authMiddleware: am,
    }

    // Public routes (auth)
    r.HandleFunc("/register", handler.Register).Methods("POST")
    r.HandleFunc("/login", handler.Login).Methods("POST")

    // Protected routes (require authentication)
    protected := r.PathPrefix("").Subrouter()
    protected.Use(am.Authenticate)
    protected.HandleFunc("/logout", handler.Logout).Methods("POST")
    protected.HandleFunc("/users", handler.ListUsers).Methods("GET")
    protected.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
    protected.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
    protected.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")
}

// Register menangani permintaan pendaftaran pengguna baru
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    user, err := h.userClient.Register(ctx, req.Username, req.Email, req.Password, req.Fullname)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusCreated, user)
}

// Login menangani permintaan login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    user, token, err := h.userClient.Login(ctx, req.Username, req.Password)
    if err != nil {
        response.Error(w, http.StatusUnauthorized, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, map[string]interface{}{
        "user":  user,
        "token": token,
    })
}

// Logout menangani permintaan logout
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
    tokenHeader := r.Header.Get("Authorization")
    if tokenHeader == "" {
        response.Error(w, http.StatusBadRequest, "Authorization header is required")
        return
    }

    // Format: "Bearer {token}"
    if len(tokenHeader) < 7 || !strings.HasPrefix(tokenHeader, "Bearer ") {
        response.Error(w, http.StatusBadRequest, "Invalid token format")
        return
    }

    token := tokenHeader[7:]
    ctx := r.Context()
    success, err := h.userClient.Logout(ctx, token)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, map[string]bool{"success": success})
}

// GetUser menangani permintaan untuk mendapatkan detail pengguna
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    user, err := h.userClient.GetUser(ctx, id)
    if err != nil {
        response.Error(w, http.StatusNotFound, "User not found")
        return
    }

    response.JSON(w, http.StatusOK, user)
}

// UpdateUser menangani permintaan pembaruan profil pengguna
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var req UpdateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    user, err := h.userClient.UpdateUser(ctx, id, req.Username, req.Email, req.Fullname)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, user)
}

// DeleteUser menangani permintaan penghapusan pengguna
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    success, err := h.userClient.DeleteUser(ctx, id)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, map[string]bool{"success": success})
}

// ListUsers menangani permintaan untuk mendapatkan daftar pengguna
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    queryValues := r.URL.Query()
    pageStr := queryValues.Get("page")
    limitStr := queryValues.Get("limit")

    page := 1
    limit := 10

    if pageStr != "" {
        pageInt, err := strconv.Atoi(pageStr)
        if err == nil && pageInt > 0 {
            page = pageInt
        }
    }

    if limitStr != "" {
        limitInt, err := strconv.Atoi(limitStr)
        if err == nil && limitInt > 0 {
            limit = limitInt
        }
    }

    ctx := r.Context()
    users, total, err := h.userClient.List(ctx, page, limit)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, map[string]interface{}{
        "data":  users,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}