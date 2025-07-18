package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "正常なマークダウン変換",
			requestBody: map[string]string{
				"markdown": "# Hello World\n\nThis is **bold** text.",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				if resp["success"] != true {
					t.Errorf("Expected success=true, got %v", resp["success"])
				}
				
				html, ok := resp["html"].(string)
				if !ok {
					t.Fatal("html field is not a string")
				}
				
				// HTMLに期待される要素が含まれているか確認
				expectedElements := []string{
					`<h1 id="hello-world">Hello World</h1>`,
					`<strong>bold</strong>`,
				}
				
				for _, elem := range expectedElements {
					if !bytes.Contains([]byte(html), []byte(elem)) {
						t.Errorf("Expected HTML to contain %q", elem)
					}
				}
			},
		},
		{
			name:           "空のリクエストボディ",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if resp["success"] != false {
					t.Errorf("Expected success=false, got %v", resp["success"])
				}
			},
		},
		{
			name: "空のマークダウン",
			requestBody: map[string]string{
				"markdown": "",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if resp["success"] != true {
					t.Errorf("Expected success=true, got %v", resp["success"])
				}
				if resp["html"] != "" {
					t.Errorf("Expected empty HTML, got %v", resp["html"])
				}
			},
		},
	}

	handler := NewRenderHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.requestBody != nil {
				var err error
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/render", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			tt.checkResponse(t, rec.Body.Bytes())
		})
	}
}

func TestRenderHandlerMethods(t *testing.T) {
	handler := NewRenderHandler()
	
	// GETリクエストをテスト
	req := httptest.NewRequest(http.MethodGet, "/api/render", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for GET request, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}