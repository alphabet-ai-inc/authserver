package api_test

import (
	"backend/api"
	"backend/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// type mockAuth struct {
// 	// Implement the interface
// 	GetTokenFromHeaderAndVerifyFunc func(w http.ResponseWriter, r *http.Request) (string, *Claims, error)
// }

// // Implement the method for the mocked auth
// func (m *mockAuth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
// 	return m.GetTokenFromHeaderAndVerifyFunc(w, r)
// }

// type mockapplication struct {
// 	auth AuthInterface // Use the Auth interface defined earlier
// }

// func (m *mockapplication) authRequired(nextHandler http.HandlerFunc) any {
// 	panic("unimplemented")
// }

// Test for CORS Middleware
func TestEnableCORS(t *testing.T) {
	// Setup test environment with allowed origins
	t.Setenv("ALLOWED_ORIGINS", "http://test-host:3000,http://localhost")

	mockDB := new(repository.MockDBRepo)
	app := api.Application{
		DB: mockDB.DatabaseRepo,
	}

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	corsHandler := app.EnableCORS(mockHandler)

	// Test case 1: Request with allowed origin
	t.Run("Allowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", "http://test-host:3000")
		recorder := httptest.NewRecorder()

		corsHandler.ServeHTTP(recorder, req)

		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "http://test-host:3000", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	// Test case 2: Request with disallowed origin
	t.Run("Disallowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", "http://unknown.com")
		recorder := httptest.NewRecorder()

		corsHandler.ServeHTTP(recorder, req)

		resp := recorder.Result()
		assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	})

	// Test case 3: OPTIONS preflight request
	t.Run("OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		recorder := httptest.NewRecorder()

		corsHandler.ServeHTTP(recorder, req)

		resp := recorder.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
		assert.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS",
			resp.Header.Get("Access-Control-Allow-Methods"))
	})
}

// Test for Authentication Middleware
// func TestAuthRequired(t *testing.T) {
// 	// Create the mock auth instance
// 	mockAuth := &mockAuth{
// 		GetTokenFromHeaderAndVerifyFunc: func(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
// 			return "token", &Claims{jwt.RegisteredClaims{Subject: "1"}}, nil
// 		},
// 	}

// 	app := &mockapplication{
// 		auth: mockAuth, // Use the mock auth
// 	}

// 	// Create a handler to test the next middleware
// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK) // Return 200 OK for the next handler
// 	})

// 	authHandler := app.authRequired(nextHandler)

// 	// Test authorized request
// 	reqAuth := httptest.NewRequest(http.MethodGet, "/", nil)
// 	recorderAuth := httptest.NewRecorder()

// 	authHandler.ServeHTTP(recorderAuth, reqAuth)

// 	// Assert the response for the authorized request
// 	assert.Equal(t, http.StatusOK, recorderAuth.Code)

// 	// Now test an unauthorized request
// 	mockAuth.GetTokenFromHeaderAndVerifyFunc = func(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
// 		return "", nil, errors.New("unauthorized") // Simulate an error
// 	}

// 	reqUnauthorized := httptest.NewRequest(http.MethodGet, "/", nil)
// 	recorderUnauthorized := httptest.NewRecorder()

// 	authHandler.ServeHTTP(recorderUnauthorized, reqUnauthorized)

// 	// Assert the response for the unauthorized request
// 	assert.Equal(t, http.StatusUnauthorized, recorderUnauthorized.Code)
// }
