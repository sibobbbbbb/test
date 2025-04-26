package http

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/domain"
)

type ResponseError struct {
    Message string `json:"message"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    User  *domain.User `json:"user"`
    Token string       `json:"token"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Fullname string `json:"fullname"`
}

type UpdateRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Fullname string `json:"fullname"`
}

type UserHandler struct {
    UserUsecase domain.UserUsecase
}

// NewUserHandler membuat instance baru UserHandler dan mendaftarkan routes
func NewUserHandler(r *mux.Router, uu domain.UserUsecase) {
    handler := &UserHandler{
        UserUsecase: uu,
    }

    // Auth routes
    r.HandleFunc("/register", handler.Register).Methods("POST")
    r.HandleFunc("/login", handler.Login).Methods("POST")
    r.HandleFunc("/logout", handler.Logout).Methods("POST")

    // User routes
    r.HandleFunc("/users", handler.List).Methods("GET")
    r.HandleFunc("/users/{id}", handler.GetByID).Methods("GET")
    r.HandleFunc("/users/{id}", handler.Update).Methods("PUT")
    r.HandleFunc("/users/{id}", handler.Delete).Methods("DELETE")
}

// Register menangani permintaan pendaftaran pengguna baru
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    user, err := h.UserUsecase.Register(ctx, req.Username, req.Email, req.Password, req.Fullname)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, user)
}

// Login menangani permintaan login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    user, token, err := h.UserUsecase.Login(ctx, req.Username, req.Password)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, err.Error())
        return
    }

    resp := LoginResponse{
        User:  user,
        Token: token,
    }

    respondWithJSON(w, http.StatusOK, resp)
}

// Logout menangani permintaan logout
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
    tokenHeader := r.Header.Get("Authorization")
    if tokenHeader == "" {
        respondWithError(w, http.StatusBadRequest, "Authorization header is required")
        return
    }

    // Format: "Bearer {token}"
    if len(tokenHeader) < 7 || tokenHeader[:7] != "Bearer " {
        respondWithError(w, http.StatusBadRequest, "Invalid token format")
        return
    }

    token := tokenHeader[7:]
    ctx := r.Context()
    err := h.UserUsecase.Logout(ctx, token)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// GetByID menangani permintaan untuk mendapatkan detail pengguna
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    user, err := h.UserUsecase.GetByID(ctx, id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "User not found")
        return
    }

    respondWithJSON(w, http.StatusOK, user)
}

// Update menangani permintaan pembaruan profil pengguna
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var req UpdateRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    ctx := r.Context()
    user, err := h.UserUsecase.Update(ctx, id, req.Username, req.Email, req.Fullname)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, user)
}

// Delete menangani permintaan penghapusan pengguna
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ctx := r.Context()
    err := h.UserUsecase.Delete(ctx, id)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// List menangani permintaan untuk mendapatkan daftar pengguna
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
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
    users, total, err := h.UserUsecase.List(ctx, page, limit)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]interface{}{
        "data":  users,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// Helper functions
func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, ResponseError{Message: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}