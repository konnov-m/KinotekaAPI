package handler

import (
	"KinotekaAPI/internal/service"
	mock_service "KinotekaAPI/internal/service/mocks"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
}

func TestHandler_userIdentity(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockUser, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(int64(10), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
		},
		{
			name:                 "Invalid Header Value",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(int64(0), errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"Can't parse token"}`,
		},
	}

	for _, test := range testTable {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.token)

			services := &service.Service{User: repo}
			handler := Handler{ser: services}

			// Init Endpoint
			http.Handle("GET /identity", handler.userIdentity(http.HandlerFunc(emptyHandler)))

			// Init Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/identity", nil)
			req.Header.Set(test.headerName, test.headerValue)

			http.DefaultServeMux.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, w.Code, test.expectedStatusCode)
			if w.Body.String() != test.expectedResponseBody {
				assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
			}
		})
	}
}
