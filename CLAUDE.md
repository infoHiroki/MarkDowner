# CLAUDE.md

このファイルは、このリポジトリのコードを扱う際のClaude Code (claude.ai/code)へのガイダンスを提供します。

## コマンド

### ビルドと実行
```bash
# アプリケーションのビルド
go build -o markdowner

# サーバーの実行（デフォルトポート: 8080）
./markdowner

# カスタムポートで実行
PORT=3000 ./markdowner

# グローバルインストール
go install github.com/hirokitakamura/markdowner@latest
```

### テスト
```bash
# 詳細出力付きで全テストを実行
go test ./... -v

# 特定パッケージのテストを実行
go test ./internal/markdown -v
go test ./internal/handler -v

# 特定のテストを実行
go test -run TestConvert ./internal/markdown -v
go test -run TestRenderHandler ./internal/handler -v
```

### 開発
```bash
# 依存関係のダウンロード
go mod tidy

# コードのフォーマット
go fmt ./...

# 問題のチェック
go vet ./...
```

## アーキテクチャ

これはTDD開発でKISS/DRY/YAGNI原則に従うミニマリストなマークダウンビューアです。

### 開発原則

- **KISS (Keep It Simple, Stupid)**: シンプルで理解しやすいコードを維持
- **DRY (Don't Repeat Yourself)**: コードの重複を避け、再利用可能なコンポーネントを作成
- **YAGNI (You Aren't Gonna Need It)**: 必要になるまで機能を追加しない、最小限の実装

### コアコンポーネント

1. **マークダウンパーサー** (`internal/markdown/`)
   - CommonMark解析用のgoldmarkライブラリをラップ
   - 単一のパブリック関数: `Convert(markdown string) string`
   - すべてのマークダウンからHTMLへの変換を処理

2. **HTTPハンドラー** (`internal/handler/`)
   - `RenderHandler`: `/api/render`へのPOSTリクエストを処理
   - JSONリクエスト/レスポンス形式
   - HTTPメソッドとリクエストボディを検証

3. **Webサーバー** (`main.go`)
   - `//go:embed`を使用してフロントエンドファイルを埋め込み
   - 埋め込みファイルシステムから静的ファイルを提供
   - APIコールを適切なハンドラーにルーティング
   - PORT環境変数でポート設定可能

4. **フロントエンド** (`web/index.html`)
   - インラインCSS/JavaScriptを持つシングルページアプリケーション
   - 300ミリ秒のデバウンス付きリアルタイムプレビュー
   - .md/.markdownファイルのアップロードサポート
   - 外部依存関係なし

### 主要な設計決定

- **ゼロ設定**: 設定ファイル不要、すぐに使える
- **単一バイナリ**: 配布が簡単なようにすべてのアセットを埋め込み
- **最小限の依存関係**: マークダウン解析用のgoldmarkのみ
- **標準ライブラリ**: Goの組み込みHTTPサーバーを使用、Webフレームワークなし
- **テスト駆動**: すべてのコンポーネントに包括的なテストカバレッジ

### API仕様

**POST /api/render**
- リクエスト: `{"markdown": "# Hello World"}`
- レスポンス: `{"html": "<h1 id=\"hello-world\">Hello World</h1>", "success": true}`
- エラーレスポンス: `{"success": false}`

### テスト戦略

- 包括的なテストケースを持つテーブル駆動テスト
- テストは通常の操作、エッジケース、エラー状態をカバー
- HTTPハンドラーは`httptest`パッケージでテスト
- マークダウンパーサーは様々なCommonMark構文に対してテスト

### パフォーマンス目標

- 起動時間: < 100ミリ秒
- 1MBファイルの変換: < 50ミリ秒
- メモリ使用量: < 50MB