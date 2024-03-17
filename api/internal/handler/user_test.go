package handler

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/service"
	mock_service "KinotekaAPI/internal/service/mocks"
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_signUp(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, user domain.User, role string)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domain.User
		role                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"login": "username", "role": "admin", "password": "qwerty"}`,
			inputUser: domain.User{
				Login:    "username",
				Password: "qwerty",
			},
			role: "admin",
			mockBehavior: func(r *mock_service.MockUser, user domain.User, role string) {
				r.EXPECT().CreateUser(user, role).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: ``,
		},
		{
			name:               "Wrong Input",
			inputBody:          `{"login": "username"}`,
			inputUser:          domain.User{},
			role:               "",
			mockBehavior:       func(r *mock_service.MockUser, user domain.User, role string) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"wrong form. Login and password are required"}
`,
		},
		{
			name:      "Service Error",
			inputBody: `{"login": "username", "role": "admin", "password": "qwerty"}`,
			inputUser: domain.User{
				Login:    "username",
				Password: "qwerty",
			},
			role: "admin",
			mockBehavior: func(r *mock_service.MockUser, user domain.User, role string) {
				r.EXPECT().CreateUser(user, role).Return(errors.New("something went wrong"))
			},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"Can't create user"}
`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser, test.role)

			services := &service.Service{User: repo}
			handler := UserHandler{services}

			// Init Endpoint
			http.Handle("POST /sign-up", middlewareLog(http.HandlerFunc(handler.signUp)))

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(test.inputBody))

			// Call your handler directly, passing in the ResponseRecorder and Request
			http.DefaultServeMux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}

func TestUserHandler_signIn(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, login, password string)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domain.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"login": "username", "password": "qwerty"}`,
			inputUser: domain.User{
				Login:    "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockUser, login, password string) {
				r.EXPECT().GenerateToken(login, password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:               "Wrong Input",
			inputBody:          `{"login": "username"}`,
			inputUser:          domain.User{},
			mockBehavior:       func(r *mock_service.MockUser, login, password string) {},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"wrong form. Login and password are required"}
`,
		},
		{
			name:      "Service Error",
			inputBody: `{"login": "username", "password": "qwerty"}`,
			inputUser: domain.User{
				Login:    "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockUser, login, password string) {
				r.EXPECT().GenerateToken(login, password).Return("", errors.New(""))
			},
			expectedStatusCode: 400,
			expectedResponseBody: `{"message":"Can't generate token"}
`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser.Login, test.inputUser.Password)

			services := &service.Service{User: repo}
			handler := UserHandler{services}

			// Init Endpoint
			http.Handle("POST /sign-in", middlewareLog(http.HandlerFunc(handler.signIn)))

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in",
				bytes.NewBufferString(test.inputBody))

			// Call your handler directly, passing in the ResponseRecorder and Request
			http.DefaultServeMux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}
