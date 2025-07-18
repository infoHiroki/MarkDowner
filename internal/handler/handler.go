package handler

import (
	"encoding/json"
	"net/http"
	
	"github.com/hirokitakamura/markdowner/internal/markdown"
)

// RenderRequest はマークダウン変換APIのリクエスト構造体
type RenderRequest struct {
	Markdown string `json:"markdown"`
}

// RenderResponse はマークダウン変換APIのレスポンス構造体
type RenderResponse struct {
	HTML    string `json:"html"`
	Success bool   `json:"success"`
}

// RenderHandler はマークダウンをHTMLに変換するハンドラー
type RenderHandler struct{}

// NewRenderHandler は新しいRenderHandlerを作成します
func NewRenderHandler() http.Handler {
	return &RenderHandler{}
}

// ServeHTTP はHTTPリクエストを処理します
func (h *RenderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// POSTメソッドのみ許可
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RenderResponse{
			Success: false,
		})
		return
	}

	// リクエストボディをデコード
	var req RenderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RenderResponse{
			Success: false,
		})
		return
	}

	// マークダウンをHTMLに変換
	html := markdown.Convert(req.Markdown)

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RenderResponse{
		HTML:    html,
		Success: true,
	})
}