package mock

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gostack-labs/bytego"
	"github.com/gostack-labs/bytego/middleware/logger"
	"github.com/gostack-labs/bytego/middleware/recovery"
)

type MockAPI struct {
	Title      string            `json:"title,omitempty"`
	HTTPMethod string            `json:"http_method,omitempty"`
	Path       string            `json:"path,omitempty"`
	Format     string            `json:"format,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Responses  []MockAPIResponse `json:"response,omitempty"`
}

type MockAPIResponse struct {
	IsMock   bool   `json:"is_mock,omitempty"`
	HTTPCode int    `json:"http_code,omitempty"`
	Body     string `json:"body,omitempty"`
}

type MockServer interface {
	InitFiles(dir string) error
	InitMockApis(mockApis []MockAPI)
	Serve() error
}

func NewMockServer(addr string) MockServer {
	if addr == "" {
		addr = "localhost:8001"
	}

	app := bytego.New()
	app.Use(recovery.New(), logger.New())

	return &mockServer{
		app:  app,
		addr: addr,
	}
}

type mockServer struct {
	app      *bytego.App
	addr     string
	mockApis []MockAPI
}

func (s *mockServer) InitFiles(dir string) error {
	return fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fullPath := filepath.Join(dir, path)
		if filepath.Ext(fullPath) != ".mocks" {
			return nil
		}
		bytes, err := os.ReadFile(filepath.Clean(fullPath))
		if err != nil {
			return err
		}
		var mockApis []MockAPI
		if err = json.Unmarshal(bytes, &mockApis); err != nil {
			return fmt.Errorf("error file %s: %w", fullPath, err)
		}
		if len(mockApis) > 0 {
			s.mockApis = append(s.mockApis, mockApis...)
		}
		return nil
	})
}

func (s *mockServer) InitMockApis(mockApis []MockAPI) {
	s.mockApis = append(s.mockApis, mockApis...)
}

func (s *mockServer) Serve() error {
	s.mock()
	fmt.Println("Mock server listen:", s.addr)
	return s.app.Run(s.addr)
}

func (s *mockServer) handler(api MockAPI) bytego.HandlerFunc {
	return func(c *bytego.Ctx) error {
		var mockResponse *MockAPIResponse
		for _, resp := range api.Responses {
			resp2 := resp
			if resp.IsMock {
				mockResponse = &resp2
				break
			}
		}
		if mockResponse == nil && len(api.Responses) > 0 {
			mockResponse = &api.Responses[0]
		}
		if mockResponse == nil {
			return c.String(http.StatusNoContent, "no mock response")
		}
		for key, value := range api.Headers {
			c.SetHeader(key, value)
		}

		c.Status(mockResponse.HTTPCode)
		body := mockResponse.Body
		if api.Format == "jsonp" {
			callback := c.Query("callback")
			callbackPrefix := "callback("
			if callback != "" && strings.HasPrefix(body, callbackPrefix) {
				body = fmt.Sprintf("%s(%s", callback, body[len(callbackPrefix):])
			}
		}
		if _, err := c.Response.WriteString(body); err != nil {
			return err
		}
		return nil
	}
}

func (s *mockServer) mock() {
	fmt.Println("Mock apis count:", len(s.mockApis))
	for _, api := range s.mockApis {
		if api.HTTPMethod == "ANY" {
			s.app.Any(api.Path, s.handler(api))
		} else {
			s.app.Handle(api.HTTPMethod, api.Path, s.handler(api))
		}
	}
}
