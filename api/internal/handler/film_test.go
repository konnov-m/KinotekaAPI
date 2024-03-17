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
)

func TestFilmHandler_film(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockFilm, s1, s2, s3 string, b bool)

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
			mockBehavior: func(r *mock_service.MockFilm, titleParam, actorParam, orderByParam string, descParam bool) {
				r.EXPECT().GetFilmsSort(orderByParam, descParam).Return([]domain.Film{
					{
						ID:          10,
						Title:       "Человек-паук 3",
						Year:        2007,
						Information: sql.NullString{String: "2:19", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.1, Valid: true},
					},
				}, nil)
			},
			expectedStatusCode: 200,
			titleParam:         "",
			actorParam:         "",
			orderByParam:       "",
			descParam:          false,
			expectedResponseBody: `[
    {
        "id": 10,
        "title": "Человек-паук 3",
        "year": 2007,
        "information": {
            "String": "2:19",
            "Valid": true
        },
        "rating": {
            "Float64": 8.1,
            "Valid": true
        }
    }
]`,
		},
		{
			name:     "Ok with title",
			addToUrl: `?title=паук`,
			mockBehavior: func(r *mock_service.MockFilm, titleParam, actorParam, orderByParam string, descParam bool) {
				r.EXPECT().GetFilmsLike(titleParam).Return([]domain.Film{
					{
						ID:          10,
						Title:       "Человек-паук 3",
						Year:        2007,
						Information: sql.NullString{String: "2:19", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.1, Valid: true},
					},
				}, nil)
			},
			expectedStatusCode: 200,
			titleParam:         "паук",
			actorParam:         "",
			orderByParam:       "",
			descParam:          false,
			expectedResponseBody: `[
    {
        "id": 10,
        "title": "Человек-паук 3",
        "year": 2007,
        "information": {
            "String": "2:19",
            "Valid": true
        },
        "rating": {
            "Float64": 8.1,
            "Valid": true
        }
    }
]`,
		},
		{
			name:     "Ok with title and order by title",
			addToUrl: `?title=паук&orderBy=title`,
			mockBehavior: func(r *mock_service.MockFilm, titleParam, actorParam, orderByParam string, descParam bool) {
				r.EXPECT().GetFilmsSortLike(orderByParam, titleParam, descParam).Return([]domain.Film{
					{
						ID:          8,
						Title:       "Человек-паук",
						Year:        2002,
						Information: sql.NullString{String: "2:01", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.3, Valid: true},
					},
					{
						ID:          9,
						Title:       "Человек-паук 2",
						Year:        2004,
						Information: sql.NullString{String: "2:07", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.2, Valid: true},
					},
					{
						ID:          10,
						Title:       "Человек-паук 3",
						Year:        2007,
						Information: sql.NullString{String: "2:19", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.1, Valid: true},
					},
				}, nil)
			},
			expectedStatusCode: 200,
			titleParam:         "паук",
			actorParam:         "",
			orderByParam:       "title",
			descParam:          false,
			expectedResponseBody: `[
    {
        "id": 8,
        "title": "Человек-паук",
        "year": 2002,
        "information": {
            "String": "2:01",
            "Valid": true
        },
        "rating": {
            "Float64": 8.3,
            "Valid": true
        }
    },
    {
        "id": 9,
        "title": "Человек-паук 2",
        "year": 2004,
        "information": {
            "String": "2:07",
            "Valid": true
        },
        "rating": {
            "Float64": 8.2,
            "Valid": true
        }
    },
    {
        "id": 10,
        "title": "Человек-паук 3",
        "year": 2007,
        "information": {
            "String": "2:19",
            "Valid": true
        },
        "rating": {
            "Float64": 8.1,
            "Valid": true
        }
    }
]`,
		},
		{
			name:     "Ok with title and order by rating",
			addToUrl: `?title=паук&orderBy=title`,
			mockBehavior: func(r *mock_service.MockFilm, titleParam, actorParam, orderByParam string, descParam bool) {
				r.EXPECT().GetFilmsSortLike(orderByParam, titleParam, descParam).Return([]domain.Film{
					{
						ID:          10,
						Title:       "Человек-паук 3",
						Year:        2007,
						Information: sql.NullString{String: "2:19", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.1, Valid: true},
					},
					{
						ID:          9,
						Title:       "Человек-паук 2",
						Year:        2004,
						Information: sql.NullString{String: "2:07", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.2, Valid: true},
					},
					{
						ID:          8,
						Title:       "Человек-паук",
						Year:        2002,
						Information: sql.NullString{String: "2:01", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.3, Valid: true},
					},
				}, nil)
			},
			expectedStatusCode: 200,
			titleParam:         "паук",
			actorParam:         "",
			orderByParam:       "title",
			descParam:          false,
			expectedResponseBody: `[
    {
        "id": 10,
        "title": "Человек-паук 3",
        "year": 2007,
        "information": {
            "String": "2:19",
            "Valid": true
        },
        "rating": {
            "Float64": 8.1,
            "Valid": true
        }
    },
    {
        "id": 9,
        "title": "Человек-паук 2",
        "year": 2004,
        "information": {
            "String": "2:07",
            "Valid": true
        },
        "rating": {
            "Float64": 8.2,
            "Valid": true
        }
    },
    {
        "id": 8,
        "title": "Человек-паук",
        "year": 2002,
        "information": {
            "String": "2:01",
            "Valid": true
        },
        "rating": {
            "Float64": 8.3,
            "Valid": true
        }
    }
]`,
		},
		{
			name:     "Wrong order by With title ",
			addToUrl: `?title=паук&orderBy=titl`,
			mockBehavior: func(r *mock_service.MockFilm, titleParam, actorParam, orderByParam string, descParam bool) {
				r.EXPECT().GetFilmsSortLike(orderByParam, titleParam, descParam).Return([]domain.Film{
					{
						ID:          10,
						Title:       "Человек-паук 3",
						Year:        2007,
						Information: sql.NullString{String: "2:19", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.1, Valid: true},
					},
					{
						ID:          9,
						Title:       "Человек-паук 2",
						Year:        2004,
						Information: sql.NullString{String: "2:07", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.2, Valid: true},
					},
					{
						ID:          8,
						Title:       "Человек-паук",
						Year:        2002,
						Information: sql.NullString{String: "2:01", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.3, Valid: true},
					},
				}, nil)
			},
			expectedStatusCode: 200,
			titleParam:         "паук",
			actorParam:         "",
			orderByParam:       "titl",
			descParam:          false,
			expectedResponseBody: `[
    {
        "id": 10,
        "title": "Человек-паук 3",
        "year": 2007,
        "information": {
            "String": "2:19",
            "Valid": true
        },
        "rating": {
            "Float64": 8.1,
            "Valid": true
        }
    },
    {
        "id": 9,
        "title": "Человек-паук 2",
        "year": 2004,
        "information": {
            "String": "2:07",
            "Valid": true
        },
        "rating": {
            "Float64": 8.2,
            "Valid": true
        }
    },
    {
        "id": 8,
        "title": "Человек-паук",
        "year": 2002,
        "information": {
            "String": "2:01",
            "Valid": true
        },
        "rating": {
            "Float64": 8.3,
            "Valid": true
        }
    }
]`,
		},
		{
			name:     "Wrong order by Without title",
			addToUrl: `?orderBy=titl`,
			mockBehavior: func(r *mock_service.MockFilm, titleParam, actorParam, orderByParam string, descParam bool) {
				r.EXPECT().GetFilmsSort(orderByParam, descParam).Return([]domain.Film{
					{
						ID:          10,
						Title:       "Человек-паук 3",
						Year:        2007,
						Information: sql.NullString{String: "2:19", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.1, Valid: true},
					},
					{
						ID:          9,
						Title:       "Человек-паук 2",
						Year:        2004,
						Information: sql.NullString{String: "2:07", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.2, Valid: true},
					},
					{
						ID:          8,
						Title:       "Человек-паук",
						Year:        2002,
						Information: sql.NullString{String: "2:01", Valid: true},
						Rating:      sql.NullFloat64{Float64: 8.3, Valid: true},
					},
				}, nil)
			},
			expectedStatusCode: 200,
			titleParam:         "",
			actorParam:         "",
			orderByParam:       "titl",
			descParam:          false,
			expectedResponseBody: `[
    {
        "id": 10,
        "title": "Человек-паук 3",
        "year": 2007,
        "information": {
            "String": "2:19",
            "Valid": true
        },
        "rating": {
            "Float64": 8.1,
            "Valid": true
        }
    },
    {
        "id": 9,
        "title": "Человек-паук 2",
        "year": 2004,
        "information": {
            "String": "2:07",
            "Valid": true
        },
        "rating": {
            "Float64": 8.2,
            "Valid": true
        }
    },
    {
        "id": 8,
        "title": "Человек-паук",
        "year": 2002,
        "information": {
            "String": "2:01",
            "Valid": true
        },
        "rating": {
            "Float64": 8.3,
            "Valid": true
        }
    }
]`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.titleParam, test.actorParam, test.orderByParam, test.descParam)

			services := &service.Service{Film: repo}
			handler := FilmHandler{services}

			// Init Endpoint
			http.Handle("GET /film", middlewareLog(http.HandlerFunc(handler.film)))

			// Create Request
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/film%s", test.addToUrl)
			req := httptest.NewRequest("GET", url, nil)

			// Call your handler directly, passing in the ResponseRecorder and Request
			http.DefaultServeMux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestFilmHandler_getFilm(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockFilm, id int64)

	tests := []struct {
		name                 string
		addToUrl             string
		mockBehavior         mockBehavior
		ID                   int64
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Ok",
			addToUrl: `/10`,
			mockBehavior: func(r *mock_service.MockFilm, id int64) {
				r.EXPECT().GetFilm(id).Return(domain.Film{
					ID:          10,
					Title:       "Человек-паук 3",
					Year:        2007,
					Information: sql.NullString{String: "2:19", Valid: true},
					Rating:      sql.NullFloat64{Float64: 8.1, Valid: true},
				}, nil)
			},
			expectedStatusCode: 200,
			ID:                 10,
			expectedResponseBody: `{
        "id": 10,
        "title": "Человек-паук 3",
        "year": 2007,
        "information": {
            "String": "2:19",
            "Valid": true
        },
        "rating": {
            "Float64": 8.1,
            "Valid": true
        }
}`,
		},
		{
			name:                 "Wrong type id",
			addToUrl:             `/asd`,
			mockBehavior:         func(r *mock_service.MockFilm, id int64) {},
			expectedStatusCode:   400,
			ID:                   10,
			expectedResponseBody: `{"message":"Can't parse id from path"}`,
		},
		{
			name:     "Out of range index",
			addToUrl: `/111`,
			mockBehavior: func(r *mock_service.MockFilm, id int64) {
				r.EXPECT().GetFilm(id).Return(domain.Film{}, errors.New(""))
			},
			expectedStatusCode:   400,
			ID:                   111,
			expectedResponseBody: `{"message":"Can't get film"}`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.ID)

			services := &service.Service{Film: repo}
			handler := FilmHandler{services}

			// Init Endpoint
			http.Handle("GET /film/{id}", middlewareLog(http.HandlerFunc(handler.getFilm)))

			// Create Request
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/film%s", test.addToUrl)
			req := httptest.NewRequest("GET", url, nil)

			// Call your handler directly, passing in the ResponseRecorder and Request
			http.DefaultServeMux.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestFilmHandler_createFilm(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockFilm, film domain.Film)
	type mockBehavior2 func(r *mock_service.MockUser, id int64)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		mockBehaviorAdmin    mockBehavior2
		inputBody            string
		inputFilm            domain.Film
		ID                   int64
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
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
			inputFilm: domain.Film{
				Title:       "Бойцовский клуб",
				Year:        1999,
				Information: sql.NullString{String: "02:19", Valid: true},
				Rating:      sql.NullFloat64{Float64: 9.1, Valid: true},
			},
			mockBehavior: func(r *mock_service.MockFilm, film domain.Film) {
				r.EXPECT().CreateFilm(film).Return(nil)
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   201,
			ID:                   10,
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
			inputFilm: domain.Film{
				Information: sql.NullString{String: "02:19", Valid: true},
				Rating:      sql.NullFloat64{Float64: 9.1, Valid: true},
			},
			mockBehavior: func(r *mock_service.MockFilm, film domain.Film) {
				r.EXPECT().CreateFilm(film).Return(errors.New("film is not valid"))
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			ID:                   10,
			expectedResponseBody: `{"message":"Can't create film"}`,
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
			inputFilm: domain.Film{
				Title:       "Бойцовский клуб",
				Year:        1999,
				Information: sql.NullString{String: "02:19", Valid: true},
				Rating:      sql.NullFloat64{Float64: 9.1, Valid: true},
			},
			mockBehavior: func(r *mock_service.MockFilm, film domain.Film) {},
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

			repo := mock_service.NewMockFilm(c)
			repo2 := mock_service.NewMockUser(c)

			test.mockBehavior(repo, test.inputFilm)
			test.mockBehaviorAdmin(repo2, test.ID)

			services := &service.Service{Film: repo, User: repo2}
			handler := FilmHandler{services}

			// Init Endpoint
			http.Handle("POST /film", middlewareLog(http.HandlerFunc(handler.createFilm)))

			// Create Request
			w := httptest.NewRecorder()
			ctx := context.WithValue(context.Background(), "userID", test.ID)
			req := httptest.NewRequest("POST", "/film",
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

func TestFilmHandler_updateFilm(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockFilm, film domain.Film)
	type mockBehavior2 func(r *mock_service.MockUser, id int64)

	tests := []struct {
		name                 string
		addToUrl             string
		mockBehavior         mockBehavior
		mockBehaviorAdmin    mockBehavior2
		inputBody            string
		inputFilm            domain.Film
		ID                   int64
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Ok",
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
			inputFilm: domain.Film{
				ID:          1,
				Title:       "Бойцовский клуб",
				Year:        1999,
				Information: sql.NullString{String: "02:19", Valid: true},
				Rating:      sql.NullFloat64{Float64: 9.1, Valid: true},
			},
			mockBehavior: func(r *mock_service.MockFilm, film domain.Film) {
				r.EXPECT().UpdateFilm(film).Return(nil)
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   201,
			ID:                   10,
			expectedResponseBody: ``,
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
			inputFilm: domain.Film{
				Title:       "Бойцовский клуб",
				Year:        1999,
				Information: sql.NullString{String: "02:19", Valid: true},
				Rating:      sql.NullFloat64{Float64: 9.1, Valid: true},
			},
			mockBehavior: func(r *mock_service.MockFilm, film domain.Film) {},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(false, nil)
			},
			expectedStatusCode:   400,
			ID:                   10,
			expectedResponseBody: `{"message":"you don't have enough permissions"}`,
		},
		{
			name:     "Wrong request",
			addToUrl: "/1",
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
			inputFilm: domain.Film{
				ID:          1,
				Information: sql.NullString{String: "02:19", Valid: true},
				Rating:      sql.NullFloat64{Float64: 9.1, Valid: true},
			},
			mockBehavior: func(r *mock_service.MockFilm, film domain.Film) {
				r.EXPECT().UpdateFilm(film).Return(errors.New("film is not valid"))
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			ID:                   10,
			expectedResponseBody: `{"message":"Can't update film"}`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			repo2 := mock_service.NewMockUser(c)

			test.mockBehavior(repo, test.inputFilm)
			test.mockBehaviorAdmin(repo2, test.ID)

			services := &service.Service{Film: repo, User: repo2}
			handler := FilmHandler{services}

			// Init Endpoint
			http.Handle("PUT /film/{id}", middlewareLog(http.HandlerFunc(handler.updateFilm)))

			// Create Request
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/film%s", test.addToUrl)
			ctx := context.WithValue(context.Background(), "userID", test.ID)
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

func TestFilmHandler_deleteFilm(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockFilm, id int64)
	type mockBehavior2 func(r *mock_service.MockUser, id int64)

	tests := []struct {
		name                 string
		addToUrl             string
		mockBehavior         mockBehavior
		mockBehaviorAdmin    mockBehavior2
		UserId               int64
		FilmId               int64
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Ok",
			addToUrl: "/1",
			mockBehavior: func(r *mock_service.MockFilm, id int64) {
				r.EXPECT().DeleteFilm(id).Return(nil)
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   204,
			UserId:               10,
			FilmId:               1,
			expectedResponseBody: ``,
		},
		{
			name:         "Not admin",
			addToUrl:     "/1",
			mockBehavior: func(r *mock_service.MockFilm, id int64) {},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(false, nil)
			},
			expectedStatusCode:   400,
			UserId:               10,
			FilmId:               1,
			expectedResponseBody: `{"message":"you don't have enough permissions"}`,
		},
		{
			name:     "Can't delete",
			addToUrl: "/1",
			mockBehavior: func(r *mock_service.MockFilm, id int64) {
				r.EXPECT().DeleteFilm(id).Return(errors.New(""))
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			UserId:               10,
			FilmId:               1,
			expectedResponseBody: `{"message":"Can't delete film"}`,
		},
		{
			name:         "Bad url",
			addToUrl:     "/asd",
			mockBehavior: func(r *mock_service.MockFilm, id int64) {},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			UserId:               10,
			FilmId:               1,
			expectedResponseBody: `{"message":"Can't parse id from path"}`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			repo2 := mock_service.NewMockUser(c)

			test.mockBehavior(repo, test.FilmId)
			test.mockBehaviorAdmin(repo2, test.UserId)

			services := &service.Service{Film: repo, User: repo2}
			handler := FilmHandler{services}

			// Init Endpoint
			http.Handle("DELETE /film/{id}", middlewareLog(http.HandlerFunc(handler.deleteFilm)))

			// Create Request
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/film%s", test.addToUrl)
			ctx := context.WithValue(context.Background(), "userID", test.UserId)
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

func TestFilmHandler_addActorsToFilm(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockFilm, filmId int64, actorId []int64)
	type mockBehavior2 func(r *mock_service.MockUser, id int64)

	tests := []struct {
		name                 string
		addToUrl             string
		inputBody            string
		inputData            Data
		mockBehavior         mockBehavior
		mockBehaviorAdmin    mockBehavior2
		UserId               int64
		FilmId               int64
		ActorsId             []int64
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Ok",
			addToUrl: "/1",
			inputBody: `{
  "actors": [
    1, 2
  ]
}`,
			inputData: Data{
				Actors: []int64{1, 2},
			},
			mockBehavior: func(r *mock_service.MockFilm, filmId int64, actorId []int64) {
				r.EXPECT().AddActorToFilm(filmId, actorId).Return(nil)
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   204,
			UserId:               10,
			FilmId:               1,
			ActorsId:             []int64{1, 2},
			expectedResponseBody: ``,
		},
		{
			name:     "Not admin",
			addToUrl: "/1",
			inputBody: `{
  "actors": [
    1, 2
  ]
}`,
			inputData: Data{
				Actors: []int64{1, 2},
			},
			mockBehavior: func(r *mock_service.MockFilm, filmId int64, actorId []int64) {},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(false, nil)
			},
			expectedStatusCode:   400,
			UserId:               10,
			FilmId:               1,
			ActorsId:             []int64{1, 2},
			expectedResponseBody: `{"message":"you don't have enough permissions"}`,
		},
		{
			name:     "Can't add",
			addToUrl: "/1",
			inputBody: `{
  "actors": [
    1, 2
  ]
}`,
			inputData: Data{
				Actors: []int64{1, 2},
			},
			mockBehavior: func(r *mock_service.MockFilm, filmId int64, actorId []int64) {
				r.EXPECT().AddActorToFilm(filmId, actorId).Return(errors.New(""))
			},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			UserId:               10,
			FilmId:               1,
			ActorsId:             []int64{1, 2},
			expectedResponseBody: `{"message":"Can't add actor to film"}`,
		},
		{
			name:     "Bad url",
			addToUrl: "/asd",
			inputBody: `{
  "actors": [
    1, 2
  ]
}`,
			inputData: Data{
				Actors: []int64{1, 2},
			},
			mockBehavior: func(r *mock_service.MockFilm, filmId int64, actorId []int64) {},
			mockBehaviorAdmin: func(r *mock_service.MockUser, id int64) {
				r.EXPECT().IsAdmin(id).Return(true, nil)
			},
			expectedStatusCode:   400,
			UserId:               10,
			FilmId:               1,
			expectedResponseBody: `{"message":"Can't parse id from path"}`,
		},
	}

	for _, test := range tests {
		http.DefaultServeMux = http.NewServeMux()
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			repo2 := mock_service.NewMockUser(c)

			test.mockBehavior(repo, test.FilmId, test.ActorsId)
			test.mockBehaviorAdmin(repo2, test.UserId)

			services := &service.Service{Film: repo, User: repo2}
			handler := FilmHandler{services}

			// Init Endpoint
			http.Handle("POST /film/{id}", middlewareLog(http.HandlerFunc(handler.addActorsToFilm)))

			// Create Request
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/film%s", test.addToUrl)
			ctx := context.WithValue(context.Background(), "userID", test.UserId)
			req := httptest.NewRequest("POST", url,
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
