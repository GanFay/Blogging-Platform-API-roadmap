package handlers

import (
	"blog/auth"
	"blog/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func createBlogH(t *testing.T, pool *pgxpool.Pool, authorID string, n int) ([]int, error) {
	t.Helper()

	var postsID []int
	for j := 1; j <= n; j++ {
		var postID int

		err := pool.QueryRow(context.Background(), `

		INSERT INTO posts (author_id, title, content, category)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
			authorID,
			fmt.Sprintf(`title%d`, j),
			fmt.Sprintf(`content%d`, j),
			fmt.Sprintf(`category%d`, j),
		).Scan(&postID)

		postsID = append(postsID, postID)

		if err != nil {
			t.Log("err in createPost")
			return postsID, err

		}
	}

	return postsID, nil
}

func deletePostsH(t *testing.T, pool *pgxpool.Pool, IDs []int) {
	t.Helper()
	for _, i := range IDs {
		_, err := pool.Exec(context.Background(), `DELETE FROM posts WHERE id = $1`, i)
		if err != nil {
			t.Fatal(err)
		}

	}
}

func TestCreateBlog_Validation(t *testing.T) {
	h, r, pool, id := setupTest(t)
	defer pool.Close()
	defer deleteTestUser(t, pool, id)
	jwt, err := auth.GenerateAccessJWT(id)
	if err != nil {
		t.Fatal(err.Error())
	}

	r.POST("/posts", h.AuthMiddleware(), h.CreateBlog)

	testTable := []struct {
		testName string
		body     string
		expected string
		code     int
		auth     bool
	}{
		{
			testName: "Unauthorized",
			body: `{
					"title": "Test1234",
					"content": "test1",
					"category": "test2",
					"tags": ["test"]
				}`,
			expected: "missing authorization header",
			code:     401,
			auth:     false,
		},
		{
			testName: "InvalidJSON",
			body: `{
				"awda": awda
				"title": "Test1",
				"content": "test1",
				"category": "test1",
				"tags": ["test1"]
				}`,
			expected: "JSON can't unmarshal body",
			code:     400,
			auth:     true,
		},
		{
			testName: "InvalidTitle",
			body: `{
				"title": "1",
				"content": "test2",
				"category": "test2",
				"tags": ["test2"]
				}`,
			expected: "Incorrect title. It must be between 3 and 50 characters long.",
			code:     400,
			auth:     true,
		},
		{
			testName: "InvalidContent",
			body: `{
				"title": "Test3",
				"content": "",
				"category": "test3",
				"tags": ["test3"]
				}`,
			expected: "Incorrect content. It must be between 3 and 500 characters long.",
			code:     400,
			auth:     true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.testName, func(t *testing.T) {
			body := strings.NewReader(testCase.body)
			req := httptest.NewRequest(http.MethodPost, "/posts", body)
			req.Header.Set("Content-Type", "application/json")
			if testCase.auth {
				req.Header.Set("Authorization", "Bearer "+jwt)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			var resp map[string]string
			err = json.Unmarshal(w.Body.Bytes(), &resp)
			if err != nil {
				t.Fatal(err.Error())
			}
			if w.Code != testCase.code {
				t.Fatal("test: ", testCase.testName, ", want: ", testCase.code, ", got: ", w.Code)
			}
			if resp["error"] != testCase.expected {
				t.Fatal("test: ", testCase.testName, ", want: ", testCase.expected, ", got: ", resp["error"])
			}
		})
	}
}

func TestCreateBlog_Success(t *testing.T) {
	h, r, pool, id := setupTest(t)
	defer pool.Close()
	defer deleteTestUser(t, pool, id)
	jwt, err := auth.GenerateAccessJWT(id)
	if err != nil {
		t.Fatal(err.Error())
	}
	r.POST("/posts", h.AuthMiddleware(), h.CreateBlog)

	body := `{
		"title": "Test",
		"content": "test",
		"category": "test",
		"tags": ["test"]	
			}`

	req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var resp map[string]string
	if w.Code != http.StatusCreated {
		t.Fatal("want: ", http.StatusCreated, ", got: ", w.Code)
	}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err.Error())
	}
	want := "post created successfully"
	if resp["message"] != want {
		t.Fatal("want: ", want, ", got: ", resp["message"])
	}
	_, err = h.DB.Exec(context.Background(), `DELETE FROM posts WHERE title=$1`, "Test")
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestGetAllPosts_Validation(t *testing.T) {
	_, _, p, id := setupTest(t)
	defer p.Close()
	defer deleteTestUser(t, p, id)
	jwt, err := auth.GenerateAccessJWT(id)
	strID := strconv.Itoa(id)
	postsID, err := createBlogH(t, p, strID, 12)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer deletePostsH(t, p, postsID)

	testTable := []struct {
		name        string
		req         func(r *gin.Engine, h *Handler)
		reqTest     *http.Request
		wantLen     int
		wantBodyErr string
		wantCode    int
		auth        bool
	}{
		{
			name: "DefaultPagination",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts`, h.AuthMiddleware(), h.GetAllPosts)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, "/posts", nil),
			wantLen:     10,
			wantBodyErr: "",
			wantCode:    http.StatusOK,
			auth:        true,
		},
		{
			name: "WithSearchTerm",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts`, h.AuthMiddleware(), h.GetAllPosts)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, `/posts?term=title`, nil),
			wantLen:     -1,
			wantBodyErr: "",
			wantCode:    http.StatusOK,
			auth:        true,
		},
		{
			name: "InvalidLenLimit",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts`, h.AuthMiddleware(), h.GetAllPosts)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, "/posts?limit=1000", nil),
			wantLen:     -1,
			wantBodyErr: "limit is too big",
			wantCode:    http.StatusBadRequest,
			auth:        true,
		},
		{
			name: "InvalidLimit",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts`, h.AuthMiddleware(), h.GetAllPosts)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, "/posts?limit=qwerty12345", nil),
			wantLen:     -1,
			wantBodyErr: "limit must be int",
			wantCode:    http.StatusBadRequest,
			auth:        true,
		},
		{
			name: "InvalidOffset",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts`, h.AuthMiddleware(), h.GetAllPosts)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, "/posts?offset=qwerty12345", nil),
			wantLen:     -1,
			wantBodyErr: "ERROR: invalid input syntax for type bigint: \"qwerty12345\" (SQLSTATE 22P02)",
			wantCode:    500,
			auth:        true,
		},
		{
			name: "Success_EmptyList",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts`, h.AuthMiddleware(), h.GetAllPosts)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, "/posts?offset=99999999", nil),
			wantLen:     0,
			wantBodyErr: "",
			wantCode:    http.StatusOK,
			auth:        true,
		},
		{
			name: "NoAuth",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts`, h.AuthMiddleware(), h.GetAllPosts)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, "/posts", nil),
			wantLen:     -1,
			wantBodyErr: "missing authorization header",
			wantCode:    http.StatusUnauthorized,
			auth:        false,
		},
		{
			name: "ID_Success",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts/:id`, h.AuthMiddleware(), h.GetPoID)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, fmt.Sprintf("/posts/%d", postsID[0]), nil),
			wantLen:     -1,
			wantBodyErr: "",
			wantCode:    http.StatusOK,
			auth:        true,
		},
		{
			name: "ID_Invalid_ID",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts/:id`, h.AuthMiddleware(), h.GetPoID)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, fmt.Sprintf("/posts/%f", 9.5), nil),
			wantLen:     -1,
			wantBodyErr: fmt.Sprintf(`invalid id: %f`, 9.5),
			wantCode:    http.StatusBadRequest,
			auth:        true,
		},
		{
			name: "ID_UserNotFound",
			req: func(r *gin.Engine, h *Handler) {
				r.GET(`/posts/:id`, h.AuthMiddleware(), h.GetPoID)
			},
			reqTest:     httptest.NewRequest(http.MethodGet, fmt.Sprintf("/posts/%d", 2147483647), nil),
			wantLen:     -1,
			wantBodyErr: fmt.Sprintf(`user not found: %d`, 2147483647),
			wantCode:    http.StatusNotFound,
			auth:        true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			h, r, _, id2 := setupTest(t)
			deleteTestUser(t, p, id2)
			testCase.req(r, h)
			w := httptest.NewRecorder()
			if testCase.auth {
				testCase.reqTest.Header.Set("Authorization", "Bearer "+jwt)
			}
			r.ServeHTTP(w, testCase.reqTest)
			if testCase.wantBodyErr == "" && testCase.wantLen != 0 {
				var posts map[string][]models.Post
				err = json.Unmarshal(w.Body.Bytes(), &posts)
				if len(posts["posts"]) != testCase.wantLen && testCase.wantLen != -1 {
					t.Fatal("wrong pagination len")

				}
			} else {
				var posts map[string]string
				err = json.Unmarshal(w.Body.Bytes(), &posts)
				if posts["error"] != testCase.wantBodyErr {

					t.Fatal("wrong error body", w.Code, posts, testCase.wantBodyErr)
				}
			}
			if err != nil {
				t.Fatal(err.Error())
			}
			if w.Code != testCase.wantCode {
				t.Fatal("want: ", testCase.wantCode, ", got: ", w.Code)
			}

		})

	}
}

func TestUpdatePosts_Valid(t *testing.T) {
	h, r, p, id := setupTest(t)
	defer p.Close()
	defer deleteTestUser(t, p, id)
	IDs, err := createBlogH(t, p, strconv.Itoa(id), 3)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer deletePostsH(t, p, IDs)
	jwt, err := auth.GenerateAccessJWT(id)
	if err != nil {
		t.Fatal(err.Error())
	}
	r.PUT(`/posts/:id`, h.AuthMiddleware(), h.UpdateBlog)
	TestTable := []struct {
		name        string
		body        string
		reqTest     func(body string) *http.Request
		wantBodyErr string
		wantCode    int
		auth        bool
	}{
		{
			name: "Unauthorized",
			body: `{
						"title": "test",
						"content": "test",
						"category": "test",
						"tags": ["test1", "test2"]
					}`,
			reqTest: func(body string) *http.Request {
				return httptest.NewRequest(http.MethodPut, fmt.Sprintf("/posts/%d", IDs[0]), strings.NewReader(body))
			},
			wantBodyErr: "missing authorization header",
			wantCode:    http.StatusUnauthorized,
			auth:        false,
		},
		{
			name: "InvalidId",
			body: `{
						"title": "test",
						"content": "test",
						"category": "test",
						"tags": ["test1", "test2"]
					}`,
			reqTest: func(body string) *http.Request {
				return httptest.NewRequest(http.MethodPut, "/posts/qwe", strings.NewReader(body))
			},
			wantBodyErr: "invalid id qwe",
			wantCode:    http.StatusBadRequest,
			auth:        true,
		},
		{
			name: "Invalid_Json",
			body: `{
						"title": 123,
						"content": 123,
						"category": 123,
						"tags": "qwerty"
					}`,
			reqTest: func(body string) *http.Request {
				return httptest.NewRequest(http.MethodPut, fmt.Sprintf("/posts/%d", IDs[1]), strings.NewReader(body))
			},
			wantBodyErr: "json: cannot unmarshal number into Go struct",
			wantCode:    http.StatusBadRequest,
			auth:        true,
		},
		{
			name: "PostNotFound",
			body: `{
						"title": "test",
						"content": "test",
						"category": "test",
						"tags": ["test1", "test2"]
					}`,
			reqTest: func(body string) *http.Request {
				return httptest.NewRequest(http.MethodPut, "/posts/2147483647", strings.NewReader(body))
			},
			wantBodyErr: "no rows in result set",
			wantCode:    http.StatusBadRequest,
			auth:        true,
		},
		{
			name: "Success",
			body: `{
						"title": "test",
						"content": "test",
						"category": "test",
						"tags": ["test1", "test2"]
					}`,
			reqTest: func(body string) *http.Request {
				return httptest.NewRequest(http.MethodPut, fmt.Sprintf("/posts/%d", IDs[2]), strings.NewReader(body))
			},
			wantBodyErr: "",
			wantCode:    http.StatusOK,
			auth:        true,
		},
	}
	for _, testCase := range TestTable {
		t.Run(testCase.name, func(t *testing.T) {
			req := testCase.reqTest(testCase.body)
			if testCase.auth {
				req.Header.Set("Authorization", "Bearer "+jwt)

			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			var resp map[string]string
			err = json.Unmarshal(w.Body.Bytes(), &resp)
			if err != nil {
				t.Fatal(err.Error())
			}
			switch testCase.wantBodyErr {
			case "":
				if resp["message"] != "successfully updated blog!" {
					t.Fatal("want: successfully updated, got: ", resp["message"])
				}
			default:
				if !strings.Contains(resp["error"], testCase.wantBodyErr) {
					t.Fatal("want: \"", testCase.wantBodyErr, "\", got: ", resp["error"])
				}
				if w.Code != testCase.wantCode {
					t.Fatal("want: ", testCase.wantCode, ", got: ", w.Code)
				}
			}
		})
	}
}

func TestUpdateBlog_NotOwner(t *testing.T) {
	h, r, p, id := setupTest(t)
	defer p.Close()
	defer deleteTestUser(t, p, id)
	IDs, err := createBlogH(t, p, strconv.Itoa(id), 1)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer deletePostsH(t, p, IDs)
	HP, err := auth.HashPassword("user123")
	if err != nil {
		t.Fatal(err.Error())
	}
	id2, err := createTestUser(t, p, "user", "user", HP)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer deleteTestUser(t, p, id2)
	jwt, err := auth.GenerateAccessJWT(id2)

	body := `{
						"title": "test",
						"content": "test",
						"category": "test",
						"tags": ["test1", "test2"]
					}`

	r.PUT(`/posts/:id`, h.AuthMiddleware(), h.UpdateBlog)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/posts/%d", IDs[0]), strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+jwt)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err.Error())
	}
	if w.Code != http.StatusForbidden {
		t.Fatal("got: ", w.Code, ", want: ", http.StatusForbidden)
	}
	if resp["message"] != "not permission" {
		t.Fatal("got: ", resp["message"], ", want: 'not permission'")
	}
}
