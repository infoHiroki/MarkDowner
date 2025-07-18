package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/hirokitakamura/markdowner/internal/handler"
)

//go:embed web/*
var webFiles embed.FS

func main() {
	// ポート番号を環境変数から取得（デフォルト: 8080）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ルーティング設定
	mux := http.NewServeMux()

	// APIエンドポイント
	mux.Handle("/api/render", handler.NewRenderHandler())

	// 静的ファイルの配信
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// ルートパスの場合は index.html を配信
			content, err := webFiles.ReadFile("web/index.html")
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(content)
		} else {
			// その他のパスは通常通り配信
			http.FileServer(http.FS(webFiles)).ServeHTTP(w, r)
		}
	})

	// サーバー設定
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// サーバー起動
	fmt.Printf("MarkDowner server is running on http://localhost:%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")
	
	// ブラウザを自動で開く
	go func() {
		time.Sleep(1 * time.Second) // サーバー起動を待つ
		openBrowser(fmt.Sprintf("http://localhost:%s", port))
	}()
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}

// ブラウザを開く関数
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Printf("Please open your browser and visit: %s\n", url)
		return
	}
	if err != nil {
		fmt.Printf("Could not open browser automatically. Please visit: %s\n", url)
	}
}