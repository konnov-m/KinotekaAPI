package handler

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"kinoteka/internal/domain"
	"kinoteka/internal/service"
	mock_service "kinoteka/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFilmHandler_actorsList(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockActor)

	tests := []struct {
		name                 string
		addToUrl             string
		mockBehavior         mockBehavior
		titleParam           string
		actorParam           string
		orderByParam         string
		descParam            bool
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Ok",
			addToUrl: ``,
			mockBehavior: func(r *mock_service.MockActor) {
				birthday, _ := time.Parse(time.RFC3339, "1980-12-11T00:00:00Z")
				r.EXPECT().GetActors().Return([]domain.Actor{
					{
						ID:          1,
						Name:        "Райан",
						Surname:     "Томас Гослинг",
						Patronymic:  sql.NullString{String: "", Valid: true},
						Birthday:    birthday,
						Sex:         "m",
						Information: sql.NullString{String: "Томас Гослинг Райан", Valid: true},
					},
				}, nil)
			},
			expectedStatusCode: 200,
			titleParam:         "",
			descParam:          false,
			expectedResponseBody: `[
    {
    "id": 1,
    "name": "Райан",
    "surname": "Томас Гослинг",
    "patronymic": {
        "String": "",
        "Valid": true
    },
    "birthday": "1980-12-11T00:00:00Z",
    "sex": "m",
    "information": {
        "String": "Томас Гослинг Райан",
        "Valid": true
    }
}
]`,
		},
		{
			name:     "Ok with films",
			addToUrl: `?withFilms=true`,
			mockBehavior: func(r *mock_service.MockActor) {
				birthday, _ := time.Parse(time.RFC3339, "1980-12-11T00:00:00Z")
				a := domain.Actor{
					ID:          1,
					Name:        "Райан",
					Surname:     "Томас Гослинг",
					Patronymic:  sql.NullString{String: "", Valid: true},
					Birthday:    birthday,
					Sex:         "m",
					Information: sql.NullString{String: "Томас Гослинг Райан", Valid: true},
				}
				f := domain.Film{
					ID:          6,
					Title:       "Брат",
					Year:        1997,
					Information: sql.NullString{String: "1:40", Valid: true},
					Rating:      sql.NullFloat64{Float64: 8.6, Valid: true},
				}
				r.EXPECT().GetActorsWithFilms().Return([]domain.ActorFilm{
					{
						Actor: a,
						Films: []domain.Film{f},
					},
				}, nil)
			},
			expectedStatusCode: 210,
			titleParam:         "",
			descParam:          false,
			expectedResponseBody: `[
    {
	
    "Actor": {
            "id": 1,
            "name": "Райан",
            "surname": "Томас Гослинг",
            "patronymic": {
                "String": "",
                "Valid": true
            },
            "birthday": "1980-12-11T00:00:00Z",
            "sex": "m",
            "information": {
                "String": "Томас Гослинг Райан",
                "Valid": true
            }
},
	"Films": [
            {
                "id": 6,
                "title": "Брат",
                "year": 1997,
                "information": {
                    "String": "1:40",
                    "Valid": true
                },
                "rating": {
                    "Float64": 8.6,
                    "Valid": true
                }
            }
]

}
]`,
		},
		{
			name:     "Can't get actors",
			addToUrl: ``,
			mockBehavior: func(r *mock_service.MockActor) {
				r.EXPECT().GetActors().Return([]domain.Actor{}, errors.New(""))
			},
			expectedStatusCode:   400,
			titleParam:           "",
			descParam:            false,
			expectedResponseBody: `{"message":"Can't get actors"}`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo)

			services := &service.Service{Actor: repo}
			handler := ActorHandler{services}

			// Init Endpoint
			http.Handle("GET /actor", middlewareLog(http.HandlerFunc(handler.actorsList)))

			// Create Request
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/actor%s", test.addToUrl)
			req := httptest.NewRequest("GET", url, nil)

			// Call your handler directly, passing in the ResponseRecorder and Request
			http.DefaultServeMux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestFilmHandler_createActor(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockActor, actor domain.Actor)
	type mockBehavior2 func(r *mock_service.MockUser, id int64)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		mockBehaviorAdmin    mockBehavior2
		inputBody            string
		inputActor           domain.Actor
		ID                   int64
		birthday             string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			inputBody: `{
    "name": "Райан",
    "surname": "Томас Гослинг",
    "patronymic": {
        "String": "",
        "Valid": true
    },
    "birthday": "1980-12-11T00:00:00Z",
    "sex": "m",
    "information": {
        "String": "Томас Гослинг Райан",
        "Valid": true
    }
}`,
			inputActor: domain.Actor{
				Name:        "Райан",
				Surname:     "Томас Гослинг",
				Patronymic:  sql.NullString{String: "", Valid: true},
				Sex:         "m",
				Information: sql.NullString{String: "Томас Гослинг Райан", Valid: true},
			},
			mockBehavior: func(r *mock_service.MockActor, actor domain.Actor) {
				r.EXPECT().CreateActor(actor).Return(nil)
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   201,
			ID:                   10,
			birthday:             "1980-12-11T00:00:00Z",
			expectedResponseBody: ``,
		},
		{
			name: "Wrong request",
			inputBody: `{
    "information": {
        "String": "02:19",
        "Valid": true
    },
    "rating": {
        "Float64": 9.1,
        "Valid": true
    }
}`,
			inputActor: domain.Actor{
				Information: sql.NullString{String: "02:19", Valid: true},
			},
			mockBehavior: func(r *mock_service.MockActor, actor domain.Actor) {
				r.EXPECT().CreateActor(actor).Return(errors.New("actor is not valid"))
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			ID:                   10,
			expectedResponseBody: `{"message":"Can't create actor"}`,
		},
		{
			name: "Not admin",
			inputBody: `{
    "title": "Бойцовский клуб",
    "year": 1999,
    "information": {
        "String": "02:19",
        "Valid": true
    },
    "rating": {
        "Float64": 9.1,
        "Valid": true
    }
}`,
			inputActor:   domain.Actor{},
			mockBehavior: func(r *mock_service.MockActor, actor domain.Actor) {},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(false, nil)
			},
			expectedStatusCode:   400,
			ID:                   10,
			expectedResponseBody: `{"message":"you don't have enough permissions"}`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			repo2 := mock_service.NewMockUser(c)

			birthday, _ := time.Parse(time.RFC3339, test.birthday)
			test.inputActor.Birthday = birthday

			test.mockBehavior(repo, test.inputActor)
			test.mockBehaviorAdmin(repo2, test.ID)

			services := &service.Service{Actor: repo, User: repo2}
			handler := ActorHandler{services}

			// Init Endpoint
			http.Handle("POST /actor", middlewareLog(http.HandlerFunc(handler.createActor)))

			// Create Request
			w := httptest.NewRecorder()
			ctx := context.WithValue(context.Background(), "userID", test.ID)
			req := httptest.NewRequest("POST", "/actor",
				bytes.NewBufferString(test.inputBody))
			req = req.WithContext(ctx)

			// Call your handler directly, passing in the ResponseRecorder and Request
			http.DefaultServeMux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			if w.Body.String() != test.expectedResponseBody {
				assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
			}
		})
	}
}

func TestFilmHandler_updateActor(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockActor, actor domain.Actor)
	type mockBehavior2 func(r *mock_service.MockUser, id int64)

	tests := []struct {
		name                 string
		addToUrl             string
		mockBehavior         mockBehavior
		mockBehaviorAdmin    mockBehavior2
		inputBody            string
		inputActor           domain.Actor
		ID                   int64
		birthday             string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Ok",
			addToUrl: "/1",
			inputBody: `{
    "name": "Райан",
    "surname": "Томас Гослинг",
    "patronymic": {
        "String": "",
        "Valid": true
    },
    "birthday": "1980-12-11T00:00:00Z",
    "sex": "m",
    "information": {
        "String": "Томас Гослинг Райан",
        "Valid": true
    }
}`,
			inputActor: domain.Actor{
				ID:          1,
				Name:        "Райан",
				Surname:     "Томас Гослинг",
				Patronymic:  sql.NullString{String: "", Valid: true},
				Sex:         "m",
				Information: sql.NullString{String: "Томас Гослинг Райан", Valid: true},
			},
			mockBehavior: func(r *mock_service.MockActor, actor domain.Actor) {
				r.EXPECT().UpdateActor(actor).Return(nil)
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   201,
			ID:                   10,
			birthday:             "1980-12-11T00:00:00Z",
			expectedResponseBody: ``,
		},
		{
			name:     "Wrong request",
			addToUrl: "/1",
			inputBody: `{
    "information": {
        "String": "02:19",
        "Valid": true
    }
}`,
			inputActor: domain.Actor{
				ID:          1,
				Information: sql.NullString{String: "02:19", Valid: true},
			},
			mockBehavior: func(r *mock_service.MockActor, actor domain.Actor) {
				r.EXPECT().UpdateActor(actor).Return(errors.New("actor is not valid"))
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			ID:                   10,
			expectedResponseBody: `{"message":"Can't update actor"}`,
		},
		{
			name:     "Not admin",
			addToUrl: "/1",
			inputBody: `{
    "title": "Бойцовский клуб",
    "year": 1999,
    "information": {
        "String": "02:19",
        "Valid": true
    },
    "rating": {
        "Float64": 9.1,
        "Valid": true
    }
}`,
			inputActor:   domain.Actor{},
			mockBehavior: func(r *mock_service.MockActor, actor domain.Actor) {},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(false, nil)
			},
			expectedStatusCode:   400,
			ID:                   10,
			expectedResponseBody: `{"message":"you don't have enough permissions"}`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			repo2 := mock_service.NewMockUser(c)

			birthday, _ := time.Parse(time.RFC3339, test.birthday)
			test.inputActor.Birthday = birthday

			test.mockBehavior(repo, test.inputActor)
			test.mockBehaviorAdmin(repo2, test.ID)

			services := &service.Service{Actor: repo, User: repo2}
			handler := ActorHandler{services}

			// Init Endpoint
			http.Handle("PUT /actor/{id}", middlewareLog(http.HandlerFunc(handler.updateActor)))

			// Create Request
			w := httptest.NewRecorder()
			ctx := context.WithValue(context.Background(), "userID", test.ID)
			url := fmt.Sprintf("/actor%s", test.addToUrl)
			req := httptest.NewRequest("PUT", url,
				bytes.NewBufferString(test.inputBody))
			req = req.WithContext(ctx)

			// Call your handler directly, passing in the ResponseRecorder and Request
			http.DefaultServeMux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			if w.Body.String() != test.expectedResponseBody {
				assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
			}
		})
	}
}

func TestFilmHandler_deleteActor(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockActor, id int64)
	type mockBehavior2 func(r *mock_service.MockUser, id int64)

	tests := []struct {
		name                 string
		addToUrl             string
		mockBehavior         mockBehavior
		mockBehaviorAdmin    mockBehavior2
		UserId               int64
		ActorId              int64
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Ok",
			addToUrl: "/1",
			mockBehavior: func(r *mock_service.MockActor, id int64) {
				r.EXPECT().DeleteActor(id).Return(nil)
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   204,
			UserId:               10,
			ActorId:              1,
			expectedResponseBody: ``,
		},
		{
			name:         "Not admin",
			addToUrl:     "/1",
			mockBehavior: func(r *mock_service.MockActor, id int64) {},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(false, nil)
			},
			expectedStatusCode:   400,
			UserId:               10,
			ActorId:              1,
			expectedResponseBody: `{"message":"you don't have enough permissions"}`,
		},
		{
			name:     "Can't delete",
			addToUrl: "/1",
			mockBehavior: func(r *mock_service.MockActor, id int64) {
				r.EXPECT().DeleteActor(id).Return(errors.New(""))
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			UserId:               10,
			ActorId:              1,
			expectedResponseBody: `{"message":"Can't delete actor"}`,
		},
		{
			name:         "Bad url",
			addToUrl:     "/asd",
			mockBehavior: func(r *mock_service.MockActor, id int64) {},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			UserId:               10,
			ActorId:              1,
			expectedResponseBody: `{"message":"Can't parse id from path"}`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			repo2 := mock_service.NewMockUser(c)

			test.mockBehavior(repo, test.ActorId)
			test.mockBehaviorAdmin(repo2, test.UserId)

			services := &service.Service{Actor: repo, User: repo2}
			handler := ActorHandler{services}

			// Init Endpoint
			http.Handle("DELETE /actor/{id}", middlewareLog(http.HandlerFunc(handler.deleteActor)))

			// Create Request
			w := httptest.NewRecorder()
			ctx := context.WithValue(context.Background(), "userID", test.UserId)
			url := fmt.Sprintf("/actor%s", test.addToUrl)
			req := httptest.NewRequest("DELETE", url,
				nil)
			req = req.WithContext(ctx)

			// Call your handler directly, passing in the ResponseRecorder and Request
			http.DefaultServeMux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			if w.Body.String() != test.expectedResponseBody {
				assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
			}
		})
	}
}

func TestFilmHandler_getActor(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockActor, id int64)

	tests := []struct {
		name                 string
		addToUrl             string
		mockBehavior         mockBehavior
		ActorId              int64
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Ok",
			addToUrl: "/1",
			mockBehavior: func(r *mock_service.MockActor, id int64) {
				birthday, _ := time.Parse(time.RFC3339, "1980-12-11T00:00:00Z")
				r.EXPECT().GetActor(id).Return(domain.Actor{
					ID:          1,
					Name:        "Райан",
					Surname:     "Томас Гослинг",
					Patronymic:  sql.NullString{String: "", Valid: true},
					Birthday:    birthday,
					Sex:         "m",
					Information: sql.NullString{String: "Томас Гослинг Райан", Valid: true},
				}, nil)
			},
			expectedStatusCode: 200,
			ActorId:            1,
			expectedResponseBody: `
{
    "id": 1,
    "name": "Райан",
    "surname": "Томас Гослинг",
    "patronymic": {
        "String": "",
        "Valid": true
    },
    "birthday": "1980-12-11T00:00:00Z",
    "sex": "m",
    "information": {
        "String": "Томас Гослинг Райан",
        "Valid": true
    }
}`,
		},
		{
			name:                 "Bad url",
			addToUrl:             "/asd",
			mockBehavior:         func(r *mock_service.MockActor, id int64) {},
			expectedStatusCode:   400,
			ActorId:              1,
			expectedResponseBody: `{"message":"Can't parse id from path"}`,
		},
		{
			name:     "Can't get",
			addToUrl: "/1",
			mockBehavior: func(r *mock_service.MockActor, id int64) {
				r.EXPECT().GetActor(id).Return(domain.Actor{}, errors.New(""))
			},
			expectedStatusCode:   400,
			ActorId:              1,
			expectedResponseBody: `{"message":"Can't get actor"}`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			repo2 := mock_service.NewMockUser(c)

			test.mockBehavior(repo, test.ActorId)

			services := &service.Service{Actor: repo, User: repo2}
			handler := ActorHandler{services}

			// Init Endpoint
			http.Handle("GET /actor/{id}", middlewareLog(http.HandlerFunc(handler.getActor)))

			// Create Request
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/actor%s", test.addToUrl)
			req := httptest.NewRequest("GET", url,
				nil)

			// Call your handler directly, passing in the ResponseRecorder and Request
			http.DefaultServeMux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			if w.Body.String() != test.expectedResponseBody {
				assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
			}
		})
	}
}
