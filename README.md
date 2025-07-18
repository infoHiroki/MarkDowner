# MarkDowner - シンプルなマークダウンビューワー

MarkDownerは、Go言語で実装されたシンプルで高速なマークダウンビューワーです。

## 特徴

- 📝 マークダウンの即座プレビュー
- 📁 複数ファイルの同時表示とタブ切り替え
- 🎯 ドラッグ&ドロップファイル読み込み
- 🖥️ サイドバー式インターフェース
- 🚀 高速・軽量（単一バイナリ）
- 🔧 ゼロコンフィグ（設定不要）
- 🌍 自動ブラウザ起動
- 🎨 シンタックスハイライト

## インストール

### 方法1: バイナリ配布
[リリースページ](https://github.com/hirokitakamura/markdowner/releases/latest)から実行ファイルをダウンロード：
- `markdowner-darwin-amd64` (macOS Intel)
- `markdowner-darwin-arm64` (macOS Apple Silicon)
- `markdowner-windows-amd64.exe` (Windows)
- `markdowner-linux-amd64` (Linux)

### 方法2: Go install
```bash
go install github.com/hirokitakamura/markdowner@latest
```

### 方法3: ソースからビルド
```bash
git clone https://github.com/hirokitakamura/markdowner.git
cd markdowner
go build -o markdowner
```

## 使い方

### 基本的な使い方
```bash
./markdowner
```
サーバーが起動し、自動でブラウザが開きます。

### エイリアス設定（推奨）
```bash
# ~/.zshrc に追加（お使いのパスに変更してください）
alias md="/Users/your-username/Documents/Dev/MarkDowner/markdowner"

# 設定を反映
source ~/.zshrc

# 使い方
md  # サーバー起動＋ブラウザ自動オープン
```

### カスタムポート
```bash
PORT=3000 ./markdowner
```

## 機能

### 1. ファイル読み込み
- **ドラッグ&ドロップ**: ファイルを直接ドロップ
- **ファイル選択**: アップロードエリアをクリック
- **複数ファイル**: 同時に複数のMarkdownファイルを開ける

### 2. タブ機能
- ファイルごとにタブ表示
- クリックで簡単切り替え
- ×ボタンでタブを閉じる

### 3. サイドバーインターフェース
- 左: ファイル管理（アップロード・タブ）
- 右: 全画面プレビュー表示

## 対応ファイル形式
- `.md`
- `.markdown`

## 技術スタック

### バックエンド
- **言語**: Go 1.21+
- **Webサーバー**: 標準ライブラリ (net/http)
- **マークダウンパーサー**: goldmark
- **静的ファイル**: go:embed

### フロントエンド
- **HTML5**: セマンティックマークアップ
- **CSS**: モダンレスポンシブデザイン
- **JavaScript**: Vanilla JS（フレームワーク不使用）

### 開発原則
- **KISS**: シンプルで理解しやすい設計
- **DRY**: コードの重複を避けた実装
- **YAGNI**: 必要最小限の機能に絞った開発
- **TDD**: テスト駆動開発

## 開発

### 必要環境
- Go 1.21以上

### コマンド
```bash
# 依存関係のダウンロード
go mod tidy

# テスト実行
go test ./... -v

# ビルド
go build -o markdowner

# フォーマット
go fmt ./...

# 静的解析
go vet ./...

# 全プラットフォーム向けビルド（Makefile使用）
make build-all

# リリース準備
make release
```

### クロスコンパイル
```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o markdowner.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o markdowner-linux

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o markdowner-darwin-amd64

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o markdowner-darwin-arm64
```

## API

### POST /api/render
マークダウンをHTMLに変換

**リクエスト:**
```json
{
  "markdown": "# Hello World\nThis is **bold** text."
}
```

**レスポンス:**
```json
{
  "html": "<h1 id=\"hello-world\">Hello World</h1>\n<p>This is <strong>bold</strong> text.</p>",
  "success": true
}
```

## パフォーマンス

- 起動時間: < 100ms
- 1MBファイルの変換: < 50ms
- メモリ使用量: < 50MB

## ライセンス

MIT License