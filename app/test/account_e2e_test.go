package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/opchaves/gin-web-app/app/model/fixture"
	"github.com/opchaves/gin-web-app/app/service"
	"github.com/stretchr/testify/assert"
)

type apiResponse struct {
	Data service.RegisterResponse `json:"data"`
}

func TestMain_AccountE2E(t *testing.T) {
	router := SetupTest(t)

	authUser := fixture.GetMockUser()
	cookie := ""

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Register Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"first_name": authUser.FirstName,
					"last_name":  authUser.LastName,
					"email":      authUser.Email,
					"password":   authUser.Password,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &apiResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.FirstName, respBody.Data.FirstName)
				assert.Equal(t, authUser.LastName, respBody.Data.LastName)
				assert.Equal(t, authUser.Email, respBody.Data.Email)
				assert.NotNil(t, respBody.Data.ID)
				assert.NotNil(t, respBody.Data.CreatedAt)
				assert.NotNil(t, respBody.Data.UpdatedAt)
			},
		},
		{
			name: "Login Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": authUser.Password,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &apiResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.FirstName, respBody.Data.FirstName)
				assert.Equal(t, authUser.LastName, respBody.Data.LastName)
				assert.Equal(t, authUser.Email, respBody.Data.Email)
				assert.NotNil(t, respBody.Data.ID)
				assert.NotNil(t, respBody.Data.CreatedAt)
				assert.NotNil(t, respBody.Data.UpdatedAt)

				assert.Contains(t, recorder.Header(), "Set-Cookie")

				cookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Get Account",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/auth/me", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &apiResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.FirstName, respBody.Data.FirstName)
				assert.Equal(t, authUser.LastName, respBody.Data.LastName)
				assert.Equal(t, authUser.Email, respBody.Data.Email)
				assert.NotNil(t, respBody.Data.ID)
				assert.NotNil(t, respBody.Data.CreatedAt)
				assert.NotNil(t, respBody.Data.UpdatedAt)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			tc.setupHeaders(t, request)
			assert.NoError(t, err)
			router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}
