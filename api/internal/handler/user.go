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

type TokenResponse struct {
	Token string `json:"token"`
}

// @Summary SignIn
// @Tags sign
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 201 {object} TokenResponse
// @Failure 400
// @Failure 500
// @Failure default
// @Router /sign-in [post]
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
	jsonData, err := json.Marshal(TokenResponse{Token: token})

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

// @Summary SignUp
// @Tags sign
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body signUpInput true "account info"
// @Success 201 {integer} integer 1
// @Failure 400
// @Failure 500
// @Failure default
// @Router /sign-up [post]
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
