package handler

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	ser *service.Service
}

type signInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signUpInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (h *UserHandler) signIn(w http.ResponseWriter, req *http.Request) {
	var in signInInput
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		newErrorResponse(w, err, "Can't decode form", http.StatusBadRequest)
		return
	}

	token, err := h.ser.User.GenerateToken(in.Login, in.Password)
	if err != nil {
		newErrorResponse(w, err, "Can't generate token", http.StatusBadRequest)
		return
	}
	data := map[string]string{
		"token": token,
	}
	jsonData, err := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

func (h *UserHandler) signUp(w http.ResponseWriter, req *http.Request) {
	var in signUpInput
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		newErrorResponse(w, err, "Can't decode form", http.StatusBadRequest)
		return
	}
	user := domain.User{
		Login:    in.Login,
		Password: in.Password,
	}

	err := h.ser.User.CreateUser(user, in.Role)
	if err != nil {
		newErrorResponse(w, err, "Can't create user", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
