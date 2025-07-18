# MarkDowner モバイル対応設計書

## 📱 プロジェクト概要

### 目的
MarkDownerをモバイルデバイスで最適に動作するよう設計・実装し、PWA対応によりネイティブアプリライクな体験を提供する。

### 現状分析
**現在のモバイル対応度: C+**
- ✅ viewport設定済み
- ✅ 基本的なタッチ操作対応
- ✅ ファイルアイコンとドラッグハンドル追加済み
- ⚠️ サイドバーレイアウトがモバイルに最適化されていない
- ⚠️ タッチ領域が小さい（ボタン等）
- ⚠️ ドラッグ&ドロップがモバイルで制限的

### 目標
**モバイル対応度: A** を達成し、以下を実現：
- モバイルファーストUI設計
- PWA対応（インストール可能）
- オフライン基本機能
- 快適なタッチ操作

## 🏗️ アーキテクチャ設計

### レスポンシブブレークポイント
```css
/* Mobile First Design */
/* Default: 320px+ (スマートフォン) */
@media (min-width: 768px)  { /* タブレット */ }
@media (min-width: 1024px) { /* デスクトップ */ }
@media (min-width: 1200px) { /* 大画面 */ }
```

### レイアウト戦略

#### Desktop（現在の実装）
```
┌─────────────┬──────────────────────┐
│   Sidebar   │     Main Content     │
│   (300px)   │       (flex: 1)      │
│             │                      │
└─────────────┴──────────────────────┘
```

#### Mobile（新設計）
```
┌────────────────────────────────────┐
│         Header Bar                 │
├────────────────────────────────────┤
│                                    │
│         Main Content               │
│         (Preview Area)             │
│                                    │
│                                    │
├────────────────────────────────────┤
│    File Management Panel           │
│    (Collapsible)                   │
└────────────────────────────────────┘
```

## 🎨 UI/UX設計

### 1. ヘッダーバー（モバイル専用）
```html
<header class="mobile-header">
  <button class="hamburger-menu">☰</button>
  <h1 class="app-title">MarkDowner</h1>
  <button class="file-upload-btn">+</button>
</header>
```

**機能:**
- ハンバーガーメニュー：ファイル管理パネル開閉
- アップロードボタン：ファイル選択ダイアログ

### 2. ファイル管理パネル（下部スライドアップ）
```html
<div class="file-panel">
  <div class="panel-handle"></div>
  <div class="file-upload-area">...</div>
  <div class="file-tabs-horizontal">...</div>
</div>
```

**特徴:**
- 下からスライドアップ
- スワイプで開閉
- 横スクロール可能なファイルタブ

### 3. タッチ操作最適化

#### タッチ領域サイズ（Apple HIG準拠）
```css
.touch-target {
  min-height: 44px;  /* iOS推奨 */
  min-width: 44px;
  padding: 12px;
}
```

#### ジェスチャー対応
- **スワイプ左右**: ファイル切り替え
- **スワイプ上**: ファイルパネル展開
- **スワイプ下**: ファイルパネル折りたたみ
- **ピンチズーム**: プレビュー拡大縮小

## 🔧 技術設計

### PWA実装

#### 1. マニフェストファイル（manifest.json）
```json
{
  "name": "MarkDowner - Markdown Viewer",
  "short_name": "MarkDowner",
  "description": "シンプルで高速なマークダウンビューワー",
  "start_url": "/",
  "display": "standalone",
  "orientation": "portrait-primary",
  "theme_color": "#3498db",
  "background_color": "#2c3e50",
  "icons": [
    {
      "src": "/icons/icon-192x192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/icons/icon-512x512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}
```

#### 2. Service Worker（sw.js）
```javascript
const CACHE_NAME = 'markdowner-v1';
const urlsToCache = [
  '/',
  '/web/index.html',
  '/api/render'  // API呼び出しもキャッシュ
];

// キャッシュ戦略: Cache First
self.addEventListener('fetch', event => {
  event.respondWith(
    caches.match(event.request)
      .then(response => response || fetch(event.request))
  );
});
```

### CSS設計（モバイル対応）

#### 1. レスポンシブレイアウト
```css
/* モバイル用レイアウト */
@media (max-width: 767px) {
  body {
    flex-direction: column;
  }
  
  .sidebar {
    transform: translateY(100%);
    position: fixed;
    bottom: 0;
    width: 100%;
    height: 60vh;
    transition: transform 0.3s ease;
  }
  
  .sidebar.open {
    transform: translateY(0);
  }
  
  .mobile-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background: #2c3e50;
    color: white;
  }
}

/* タブレット用レイアウト */
@media (min-width: 768px) and (max-width: 1023px) {
  .sidebar {
    width: 250px;
  }
}
```

#### 2. タッチ対応コンポーネント
```css
.file-tab-mobile {
  min-height: 44px;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  border-radius: 8px;
  margin-bottom: 8px;
}

.close-btn-mobile {
  min-width: 44px;
  min-height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
}
```

### JavaScript設計（タッチ機能）

#### 1. ジェスチャーハンドリング
```javascript
class TouchGestureHandler {
  constructor() {
    this.touchStartX = 0;
    this.touchStartY = 0;
    this.minSwipeDistance = 50;
  }
  
  handleTouchStart(e) {
    this.touchStartX = e.touches[0].clientX;
    this.touchStartY = e.touches[0].clientY;
  }
  
  handleTouchEnd(e) {
    const touchEndX = e.changedTouches[0].clientX;
    const touchEndY = e.changedTouches[0].clientY;
    
    const deltaX = touchEndX - this.touchStartX;
    const deltaY = touchEndY - this.touchStartY;
    
    if (Math.abs(deltaX) > Math.abs(deltaY)) {
      // 横スワイプ: ファイル切り替え
      if (deltaX > this.minSwipeDistance) {
        this.switchToPreviousFile();
      } else if (deltaX < -this.minSwipeDistance) {
        this.switchToNextFile();
      }
    } else {
      // 縦スワイプ: パネル開閉
      if (deltaY < -this.minSwipeDistance) {
        this.openFilePanel();
      } else if (deltaY > this.minSwipeDistance) {
        this.closeFilePanel();
      }
    }
  }
}
```

#### 2. モバイル専用機能
```javascript
// Web Share API対応
function shareFile() {
  if (navigator.share) {
    navigator.share({
      title: 'MarkDown File',
      text: openFiles[activeFileIndex].content,
      url: window.location.href
    });
  }
}

// デバイス向き変更対応
window.addEventListener('orientationchange', () => {
  setTimeout(() => {
    adjustLayoutForOrientation();
  }, 100);
});
```

## 📱 モバイル専用機能

### 1. ファイル管理
- **カメラ撮影**: 画像からテキスト抽出（OCR）
- **クリップボード**: マークダウンテキストの貼り付け
- **ファイル共有**: Web Share API活用

### 2. UX改善
- **ハプティックフィードバック**: 操作確認
- **ローディングアニメーション**: API呼び出し時
- **プルツーリフレッシュ**: ファイル一覧更新

### 3. オフライン機能
- **マークダウン変換**: Service Worker内で実行
- **ファイルキャッシュ**: IndexedDBに保存
- **オフライン表示**: ネットワーク状態の表示

## 🚀 実装フェーズ

### Phase 1: 基本レスポンシブ対応（2-3時間）
1. CSS media query実装
2. モバイルレイアウト構築
3. タッチ領域最適化

### Phase 2: PWA基本機能（1-2時間）
1. manifest.json作成
2. Service Worker実装
3. アイコン作成・設定

### Phase 3: ジェスチャー対応（2-3時間）
1. タッチイベントハンドリング
2. スワイプ操作実装
3. アニメーション追加

### Phase 4: モバイル専用機能（2-3時間）
1. Web Share API統合
2. オフライン機能拡張
3. UX細部調整

## 📊 成果指標

### パフォーマンス目標
- **Lighthouse モバイルスコア**: 90+
- **First Contentful Paint**: <2秒
- **Time to Interactive**: <3秒

### ユーザビリティ目標
- **タッチ成功率**: 95%+
- **ジェスチャー認識率**: 90%+
- **オフライン機能動作率**: 100%

### 対応デバイス
- **iOS**: Safari 14+
- **Android**: Chrome 80+
- **画面サイズ**: 320px〜
- **向き**: Portrait/Landscape両対応

## 🔧 開発環境

### 必要ツール
- モバイルデバイス実機テスト
- Chrome DevTools（デバイスシミュレーション）
- Lighthouse（パフォーマンス測定）

### テスト戦略
1. **デスクトップ動作確認**: 既存機能の非破壊
2. **モバイル実機テスト**: iOS/Android各種デバイス
3. **PWA動作確認**: インストール・オフライン機能
4. **パフォーマンステスト**: Lighthouse測定

## 📝 実装チェックリスト

### Phase 1: レスポンシブ対応
- [ ] モバイルレイアウト CSS実装
- [ ] タッチ領域サイズ調整
- [ ] ブレークポイント設定
- [ ] フォント・アイコンサイズ調整

### Phase 2: PWA基本機能
- [ ] manifest.json作成
- [ ] Service Worker実装
- [ ] PWAアイコン作成（192x192, 512x512）
- [ ] インストール可能性確認

### Phase 3: ジェスチャー機能
- [ ] タッチイベントハンドラー実装
- [ ] スワイプジェスチャー検出
- [ ] アニメーション効果追加
- [ ] ファイル切り替えジェスチャー

### Phase 4: モバイル専用機能
- [ ] Web Share API統合
- [ ] オフライン機能拡張
- [ ] ハプティックフィードバック
- [ ] プルツーリフレッシュ

### テスト項目
- [ ] iPhone Safari動作確認
- [ ] Android Chrome動作確認
- [ ] タブレット動作確認
- [ ] オフライン機能テスト
- [ ] Lighthouseスコア測定

この設計により、現在のデスクトップ体験を損なうことなく、モバイルユーザーに最適化された体験を提供できます。