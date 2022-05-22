package mock

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gostack-labs/bytego"
	"github.com/gostack-labs/bytego/middleware/logger"
	"github.com/gostack-labs/bytego/middleware/recovery"
)

type MockAPI struct {
	Title      string            `json:"title,omitempty"`
	HTTPMethod string            `json:"http_method,omitempty"`
	Path       string            `json:"path,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Responses  []MockAPIResponse `json:"response,omitempty"`
}

type MockAPIResponse struct {
	IsMock   bool   `json:"is_mock,omitempty"`
	HTTPCode int    `json:"http_code,omitempty"`
	Body     string `json:"body,omitempty"`
}

type MockServer struct {
	app      *bytego.App
	addr     string
	mockApis []MockAPI
}

func New(addr string) *MockServer {
	if addr == "" {
		addr = "localhost:8001"
	}
	app := bytego.New()
	app.Use(recovery.New(), logger.New())
	return &MockServer{
		app:  app,
		addr: addr,
	}
}

func (m *MockServer) InitFiles(dir string) error {
	return fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fullPath := filepath.Join(dir, path)
		if filepath.Ext(fullPath) != ".mocks" {
			return nil
		}
		bytes, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}
		var mockApis []MockAPI
		if err = json.Unmarshal(bytes, &mockApis); err != nil {
			return fmt.Errorf("error file %s: %v", fullPath, err)
		}
		if len(mockApis) > 0 {
			m.mockApis = append(m.mockApis, mockApis...)
		}
		return nil
	})
}

func (m *MockServer) InitMockApis(mockApis []MockAPI) {
	m.mockApis = append(m.mockApis, mockApis...)
}

func (s *MockServer) mock() {
	fmt.Println("Mock apis count:", len(s.mockApis))
	handler := func(api MockAPI) bytego.HandlerFunc {
		return func(c *bytego.Ctx) error {
			var mockResponse *MockAPIResponse
			for _, resp := range api.Responses {
				if resp.IsMock {
					mockResponse = &resp
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
			if _, err := c.Response.WriteString(mockResponse.Body); err != nil {
				return err
			}
			return nil
		}
	}
	for _, api := range s.mockApis {
		if api.HTTPMethod == "ANY" {
			s.app.Any(api.Path, handler(api))
		} else {
			s.app.Handle(api.HTTPMethod, api.Path, handler(api))
		}
	}
}

func (s *MockServer) Serve() error {
	s.mock()
	fmt.Println("Mock server listen:", s.addr)
	return s.app.Run(s.addr)
}
