package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
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

// convert はマークダウン文字列をHTMLに変換します
func convert(markdown string) string {
	if markdown == "" {
		return ""
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return ""
	}

	return buf.String()
}

// Handler はVercel FunctionsのエントリーポイントでHTTPリクエストを処理します
func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// OPTIONSリクエストへの対応
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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
	html := convert(req.Markdown)

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RenderResponse{
		HTML:    html,
		Success: true,
	})
}