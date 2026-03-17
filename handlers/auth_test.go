package handlers

import (
	"blog/auth"
	"blog/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type errorResponse struct {
	Error string `json:"error"`
}

func setupTest(t *testing.T) (*Handler, *gin.Engine, *pgxpool.Pool) {
	t.Helper()

	dbURL := "postgres://app1:app@localhost:5432/db?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	h := &Handler{DB: pool}
	r := gin.Default()

	return h, r, pool
}

func performJSONRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func decodeJSON[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()

	var v T
	err := json.Unmarshal(w.Body.Bytes(), &v)
	if err != nil {
		t.Fatalf("failed to decode response body: %v; body: %s", err, w.Body.String())
	}

	return v
}

func createTestUser(t *testing.T, pool *pgxpool.Pool, username, email, passwordHash string) int {
	t.Helper()

	var id int
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO users (username, email, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		username, email, passwordHash,
	).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}

	return id
}

func deleteTestUser(t *testing.T, pool *pgxpool.Pool, id int) {
	t.Helper()

	_, err := pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegister_Validation(t *testing.T) {
	h, r, pool := setupTest(t)
	defer pool.Close()

	r.POST("/auth/register", h.Register)

	tests := []struct {
		name          string
		body          string
		wantStatus    int
		wantContains  []string
		notEmptyError bool
	}{
		{
			name: "invalid json fields",
			body: `{
				"user2name": "test_reg",
				"emai2l": "testreg@test.com",
				"passw2ord": "testreg123"
			}`,
			wantStatus: http.StatusBadRequest,
			wantContains: []string{
				"RegisterRequest.Username",
				"RegisterRequest.Email",
				"RegisterRequest.Password",
			},
			notEmptyError: true,
		},
		{
			name: "empty username",
			body: `{
				"username": "",
				"email": "testreg@test.com",
				"password": "testreg123"
			}`,
			wantStatus: http.StatusBadRequest,
			wantContains: []string{
				"RegisterRequest.Username",
			},
			notEmptyError: true,
		},
		{
			name: "empty email",
			body: `{
				"username": "test_reg",
				"email": "",
				"password": "testreg123"
			}`,
			wantStatus: http.StatusBadRequest,
			wantContains: []string{
				"RegisterRequest.Email",
			},
			notEmptyError: true,
		},
		{
			name: "empty password",
			body: `{
				"username": "test_reg",
				"email": "testreg@test.com",
				"password": ""
			}`,
			wantStatus: http.StatusBadRequest,
			wantContains: []string{
				"RegisterRequest.Password",
			},
			notEmptyError: true,
		},
		{
			name: "invalid email",
			body: `{
				"username": "test_reg",
				"email": "testregtest12312com",
				"password": "testreg123"
			}`,
			wantStatus: http.StatusBadRequest,
			wantContains: []string{
				"mail: missing '@' or angle-addr",
			},
			notEmptyError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performJSONRequest(r, http.MethodPost, "/auth/register", tt.body)

			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d, body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			resp := decodeJSON[errorResponse](t, w)

			if tt.notEmptyError && resp.Error == "" {
				t.Fatal("expected non-empty error")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(resp.Error, want) {
					t.Fatalf("expected error to contain %q, got %q", want, resp.Error)
				}
			}
		})
	}
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	h, r, pool := setupTest(t)
	defer pool.Close()

	username := "test_reg_exists"
	email := "test_exists@test.com"
	password := "test123"

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	id := createTestUser(t, pool, username, email, passwordHash)
	defer deleteTestUser(t, pool, id)

	r.POST("/auth/register", h.Register)

	body := fmt.Sprintf(`{
		"username": "%s",
		"email": "%s",
		"password": "%s"
	}`, username, email, password)

	w := performJSONRequest(r, http.MethodPost, "/auth/register", body)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusConflict, w.Code, w.Body.String())
	}

	resp := decodeJSON[errorResponse](t, w)

	if !strings.Contains(resp.Error, "SQLSTATE 23505") {
		t.Fatalf("expected duplicate key error, got: %q", resp.Error)
	}
}

func TestRegister_Success(t *testing.T) {
	h, r, pool := setupTest(t)
	defer pool.Close()

	r.POST("/auth/register", h.Register)

	username := "test_reg_success"
	email := "testregsuccess@test.com"
	password := "testreg123"

	body := fmt.Sprintf(`{
		"username": "%s",
		"email": "%s",
		"password": "%s"
	}`, username, email, password)

	w := performJSONRequest(r, http.MethodPost, "/auth/register", body)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var userID int
	err := h.DB.QueryRow(
		context.Background(),
		"SELECT id FROM users WHERE username = $1",
		username,
	).Scan(&userID)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteTestUser(t, pool, userID)
}

func TestLogin_Validation(t *testing.T) {
	h, r, pool := setupTest(t)
	defer pool.Close()

	passwordHash, err := auth.HashPassword("test123")
	if err != nil {
		t.Fatal(err)
	}

	id := createTestUser(t, pool, "test_log_validation", "test_validation@test.com", passwordHash)
	defer deleteTestUser(t, pool, id)

	r.POST("/auth/login", h.Login)

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantError  string
	}{
		{
			name: "invalid json fields",
			body: `{
				"usern2ame": "test_log",
				"passw2ord": "test123log"
			}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "username is too short or too long",
		},
		{
			name: "empty username",
			body: `{
				"username": "",
				"password": "test123"
			}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "username is too short or too long",
		},
		{
			name: "empty password",
			body: `{
				"username": "test_log",
				"password": ""
			}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "password is too short or too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performJSONRequest(r, http.MethodPost, "/auth/login", tt.body)

			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d, body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			resp := decodeJSON[errorResponse](t, w)

			if resp.Error != tt.wantError {
				t.Fatalf("expected error %q, got %q", tt.wantError, resp.Error)
			}
		})
	}
}

func TestLogin_Success(t *testing.T) {
	h, r, pool := setupTest(t)
	defer pool.Close()

	username := "test_log"
	password := "test123log"
	email := "testlog@test.com"

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	userID := createTestUser(t, pool, username, email, passwordHash)
	defer deleteTestUser(t, pool, userID)

	r.POST("/auth/login", h.Login)

	body := fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, username, password)

	w := performJSONRequest(r, http.MethodPost, "/auth/login", body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	resp := decodeJSON[models.TestLoginResponse](t, w)

	userIDAccessJWT, err := auth.ParseJWTAccess(resp.AccessToken)
	if err != nil {
		t.Fatal(err)
	}

	if userIDAccessJWT != userID {
		t.Fatalf("wrong access token user id: expected %d, got %d", userID, userIDAccessJWT)
	}

	cookies := w.Header()["Set-Cookie"]
	if len(cookies) == 0 {
		t.Fatal("no cookies in response")
	}

	cookie := cookies[0]
	cookieParts := strings.Split(cookie, ";")
	tokenPart := strings.SplitN(cookieParts[0], "=", 2)
	if len(tokenPart) != 2 {
		t.Fatalf("invalid cookie format: %s", cookie)
	}

	userIDRefreshJWT, err := auth.ParseJWTRefresh(tokenPart[1])
	if err != nil {
		t.Fatal(err)
	}

	if userIDRefreshJWT != userID {
		t.Fatalf("wrong refresh token user id: expected %d, got %d", userID, userIDRefreshJWT)
	}
}

func TestLogin_Invalid(t *testing.T) {
	h, r, pool := setupTest(t)
	defer pool.Close()

	r.POST("/auth/login", h.Login)

	body := `{
		"username": "",
		"password": ""
	}`

	w := performJSONRequest(r, http.MethodPost, "/auth/login", body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}
