# CLAUDE.md

このファイルは、このリポジトリのコードを扱う際のClaude Code (claude.ai/code)へのガイダンスを提供します。

## プロジェクト概要

MarkDownerは、シンプルで高速なマークダウンビューワーをGo言語で実装したプロジェクトです。KISS、DRY、YAGNI原則に従い、必要最小限の機能に絞って開発されています。

### スコープ
- **含む**: マークダウンファイルの読み込みとHTML変換・表示、複数ファイルのタブ表示、ドラッグ&ドロップ、Web版配布、PWA対応
- **含まない**: 編集機能、ファイル保存、プラグイン機構

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
go install github.com/infoHiroki/MarkDowner@latest

# クロスプラットフォームビルド（Makefile使用）
make build-all
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

# カバレッジレポート生成
make test-coverage
```

### 開発
```bash
# 依存関係のダウンロード
go mod tidy

# コードのフォーマット
go fmt ./...

# 問題のチェック
go vet ./...

# 完全な開発サイクル
make dev
```

## アーキテクチャ

### 全体構成
```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Browser   │────▶│  Go Server   │────▶│  File I/O   │
│  (Client)   │◀────│   (API)      │◀────│             │
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Markdown   │
                    │   Parser     │
                    └──────────────┘
```

### 開発原則

- **KISS (Keep It Simple, Stupid)**: シンプルで理解しやすいコードを維持
- **DRY (Don't Repeat Yourself)**: コードの重複を避け、再利用可能なコンポーネントを作成
- **YAGNI (You Aren't Gonna Need It)**: 必要になるまで機能を追加しない、最小限の実装
- **TDD (Test-Driven Development)**: テスト駆動開発で品質を保証

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
   - 起動時に自動でブラウザを開く

4. **フロントエンド** (`web/index.html`)
   - インラインCSS/JavaScriptを持つシングルページアプリケーション
   - サイドバーレイアウト（左：ファイル管理、右：プレビュー）
   - 複数ファイルのタブ管理
   - ドラッグ&ドロップファイルアップロード
   - 外部依存関係なし

## ディレクトリ構成

```
MarkDowner/
├── main.go              # エントリーポイント（ローカル版）
├── go.mod              # Goモジュール定義
├── go.sum              # 依存関係ロック
├── internal/
│   ├── markdown/       # マークダウン処理
│   │   ├── parser.go
│   │   └── parser_test.go
│   └── handler/        # リクエストハンドラー
│       ├── handler.go
│       └── handler_test.go
├── api/                # Vercel Functions
│   ├── render.go       # マークダウン変換API
│   ├── go.mod
│   └── go.sum
├── public/             # Vercel静的ファイル
│   ├── index.html
│   ├── favicon.svg
│   └── manifest.json
├── web/                # ローカル版静的ファイル
│   ├── index.html
│   └── static/
│       └── favicon.svg
├── docs/               # GitHub Pages用
│   ├── index.html      # プロジェクトサイト
│   └── MOBILE_DESIGN.md # モバイル対応設計書
├── assets/             # アイコンアセット
│   └── icons/
├── vercel.json         # Vercelデプロイ設定
├── CLAUDE.md           # このファイル
├── README.md           # プロジェクトドキュメント
├── LICENSE             # MITライセンス
├── Makefile           # ビルド自動化
└── .github/
    └── workflows/
        └── release.yml # CI/CD設定
```

## API仕様

### POST /api/render
マークダウンをHTMLに変換

**リクエスト:**
```json
{
  "markdown": "# Hello World\nThis is **bold** text."
}
```

**レスポンス（成功）:**
```json
{
  "html": "<h1 id=\"hello-world\">Hello World</h1>\n<p>This is <strong>bold</strong> text.</p>",
  "success": true
}
```

**レスポンス（エラー）:**
```json
{
  "success": false
}
```

## テスト戦略

### ユニットテスト
- パーサー機能: 各種マークダウン記法の変換
- ハンドラー: HTTPリクエスト/レスポンス
- エラーハンドリング

### テスト方針
- 包括的なテストケースを持つテーブル駆動テスト
- テストは通常の操作、エッジケース、エラー状態をカバー
- HTTPハンドラーは`httptest`パッケージでテスト
- マークダウンパーサーは様々なCommonMark構文に対してテスト

## セキュリティ考慮事項
- XSS対策: goldmarkの安全なHTML出力
- パストラバーサル対策: ファイルパス検証
- CORS設定: 同一オリジンのみ許可

## パフォーマンス目標
- 起動時間: < 100ミリ秒
- 1MBファイルの変換: < 50ミリ秒
- メモリ使用量: < 50MB

## 主要な設計決定

- **ゼロ設定**: 設定ファイル不要、すぐに使える
- **単一バイナリ**: 配布が簡単なようにすべてのアセットを埋め込み
- **最小限の依存関係**: マークダウン解析用のgoldmarkのみ
- **標準ライブラリ**: Goの組み込みHTTPサーバーを使用、Webフレームワークなし
- **テスト駆動**: すべてのコンポーネントに包括的なテストカバレッジ

## 制約事項
- 外部依存は最小限に
- 設定ファイル不要（ゼロコンフィグ）
- シングルバイナリで配布可能

## 開発規約

### コミット規約
コミットメッセージは以下の形式で記述すること：

```
<絵文字> <type>: <概要>

<詳細説明（任意）>

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

#### 絵文字とタイプの組み合わせ
- 🎉 `feat`: 新機能の追加
- 🐛 `fix`: バグ修正
- 📚 `docs`: ドキュメントの変更
- 🎨 `style`: コードフォーマットの変更（機能に影響しない）
- ♻️ `refactor`: リファクタリング
- ⚡ `perf`: パフォーマンス改善
- ✅ `test`: テストの追加・修正
- 🔧 `chore`: ビルドプロセスやツールの変更
- 🚀 `deploy`: デプロイ関連の変更
- 📱 `mobile`: モバイル対応
- 🌐 `web`: Web関連の変更

#### コミットの原則
- **アトミック**: 1つのコミットで1つの論理的変更
- **日本語**: コミットメッセージは日本語で記述
- **具体的**: 何を変更したかを明確に記述
- **50文字以内**: 概要は簡潔に

#### 例
```bash
🎉 feat: ドラッグ&ドロップによるファイル並び替え機能を追加

ファイルタブのドラッグ&ドロップで順序変更が可能になりました。
- マウスイベントとタッチイベントの両方に対応
- ドラッグ中の視覚的フィードバックを追加
- 並び替え後のアクティブタブの状態を適切に管理

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## デプロイ

### Web版（Vercel）
プロジェクトはVercelで自動デプロイされます：
- **本番環境**: https://mark-downer-sigma.vercel.app
- **デプロイ方法**: mainブランチへのpushで自動デプロイ
- **プレビュー**: PRごとに自動でプレビューURL生成

### ローカル版
```bash
# 単一バイナリを生成
go build -o markdowner

# 実行
./markdowner
```

## 新機能追加の履歴
- ✅ ドラッグ&ドロップによるファイル並び替え
- ✅ ファイルごとのスクロール位置記憶
- ✅ SVGアイコンセット追加
- ✅ Web版対応（Vercel Functions）
- ✅ PWA基盤（manifest.json）