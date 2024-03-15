package handler

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserHandler struct {
	s *service.UserService
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
	if req.Method != http.MethodPost {
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
		log.Printf("Request not supported method")
		return
	}
	var in signInInput

	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.s.GenerateToken(in.Login, in.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
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
	if req.Method != http.MethodPost {
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
		log.Printf("Request not supported method")
		return
	}

	var in signUpInput
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return
	}
	user := domain.User{
		Login:    in.Login,
		Password: in.Password,
	}

	err := h.s.CreateUser(user, in.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("HTTP %d - %s", http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}
