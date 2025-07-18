# Vibe Logger Go

## プロジェクト概要
AI駆動開発のためのGo言語向けログライブラリ。LLMが理解しやすい構造化されたログを生成し、効率的なデバッグと問題解決を支援します。オリジナルのvibe-logger（https://github.com/fladdict/vibe-logger）のGo版実装です。

## Memory Bank構造（軽量化済み）
このプロジェクトでは効率的な階層化Memory Bankシステムを使用しています：

### コア（常時参照）
- 現在の状況: @.claude/core/current.md
- 次のアクション: @.claude/core/next.md
- プロジェクト概要: @.claude/core/overview.md
- クイックテンプレート: @.claude/core/templates.md

### コンテキスト（必要時参照）
- 技術詳細: @.claude/context/tech.md
- 履歴・決定事項: @.claude/context/history.md

### アーカイブ（定期整理）
- 完了済み情報: @.claude/archive/

## カスタムコマンド

### 基本コマンド
- `/project:plan` - 作業計画立案
- `/project:act` - 計画に基づく実装実行
- `/project:focus` - 現在タスクに集中
- `/project:daily` - 日次振り返り（3分以内）

### 専門化モード
- `/debug:start` - デバッグ特化モード
- `/feature:plan` - 新機能設計モード
- `/review:check` - コードレビューモード

### タグ検索
- タグ形式: `#tag_name` でMemory Bank内検索
- 主要タグ: #urgent #bug #feature #completed

## 開発規約（Core Development Rules）

### 1. パッケージ管理
- **推奨ツール**: go mod（Go標準）
- **インストール**: `go get package` 形式
- **実行**: `go run .` 形式
- **禁止事項**: 
  - 不要なvendoring
  - 古いGOPATHスタイル
  - 非推奨パッケージの使用

### 2. コード品質基準
- **型注釈**: 全ての関数・変数にGoDoc形式のコメント
- **ドキュメント**: パブリックAPI・複雑な処理に必須
- **関数設計**: 単一責任・小さな関数を心がける
- **エラーハンドリング**: Go標準のエラーハンドリングパターン
- **既存パターン**: 必ず既存コードのパターンに従う
- **行長制限**: 120文字

### 3. テスト要件
- **テストフレームワーク**: Go標準のtestingパッケージ
- **カバレッジ目標**: 重要な機能は80%以上
- **必須テストケース**: 
  - エッジケース（境界値・異常値）
  - エラーハンドリング
  - 新機能には対応するテスト
  - バグ修正には回帰テスト

### 4. Git/PR規約

#### コミットメッセージ
- **基本形式**: `[prefix]: [変更内容]`
- **prefix一覧**:
  - `feat`: 新機能追加
  - `fix`: バグ修正
  - `docs`: ドキュメント更新
  - `style`: フォーマット・空白等（動作変更なし）
  - `refactor`: リファクタリング（機能変更なし）
  - `test`: テスト追加・修正
  - `chore`: ビルド・依存関係・設定変更

#### 必須トレーラー
```bash
# バグ報告ベースの修正
git commit -m "fix: resolve memory leak in logger" --trailer "Reported-by: Username"

# GitHub Issue関連
git commit -m "feat: add file logger" --trailer "Github-Issue: #123"
```

#### Pull Request規約
- **タイトル**: コミットメッセージと同様の形式
- **説明要件**:
  - **背景**: なぜこの変更が必要か
  - **変更内容**: 何を変更したか（高レベル）
  - **影響範囲**: どこに影響するか
  - **テスト**: どのようにテストしたか
- **レビュー**:
  - 適切なレビュアーを指定
  - セルフレビューを先に実施
- **禁止事項**:
  - `Co-authored-by` 等のツール言及禁止
  - 単純な作業ログの羅列

## 実行コマンド一覧

### 基本開発フロー
```bash
# プロジェクトセットアップ（初回のみ）
go mod init vibe-logger-go
go mod tidy

# 開発・テスト
go run .                     # 実行
go test ./...                # 全テスト実行
go test -v ./...             # 詳細テスト実行
go test -cover ./...         # カバレッジ付きテスト

# 品質チェック
go fmt ./...                 # コードフォーマット適用
go vet ./...                 # 静的解析
golint ./...                 # リントチェック（要インストール）
go mod verify                # 依存関係検証

# ビルド・リリース
go build .                   # ビルド
go build -o bin/vibe-logger  # 実行ファイル生成
go install                   # インストール
```

### パッケージ管理
```bash
go get [package-name]        # 依存関係追加
go mod tidy                  # 依存関係整理
go mod download              # 依存関係ダウンロード
go list -m all               # 全依存関係表示
```

## エラー対応ガイド（Error Resolution）

### 1. 問題解決の標準順序
エラーが発生した際は、以下の順序で対処することで効率的に問題を解決できます：

1. **フォーマットエラー** → `go fmt ./...`
2. **型エラー** → `go build .`
3. **リントエラー** → `go vet ./...`
4. **テストエラー** → `go test ./...`

### 2. よくある問題と解決策

#### フォーマット・リント関連
- **行長エラー**: 適切な箇所で改行、文字列は`+`で分割
- **インポート順序**: `go fmt`で自動修正
- **未使用インポート**: 不要な import を削除

#### 型チェック関連
- **nil pointer エラー**: nil チェックを追加
- **型推論エラー**: 明示的な型宣言を追加
- **インターフェース**: 適切なインターフェース実装を確認

#### テスト関連
- **テスト環境**: 必要な依存関係・設定を確認
- **並行テスト**: goroutine の適切な処理を確認
- **モック**: 外部依存関係の適切なモック化

### 3. ベストプラクティス

#### 開発時の心がけ
- **コミット前**: `go test ./...` で総合チェック
- **最小変更**: 一度に多くの変更を避ける
- **既存パターン**: 既存コードの書き方に合わせる
- **段階的修正**: 大きな変更は小さく分割

#### エラー対応時
- **エラーメッセージを熟読**: 具体的な原因を特定
- **コンテキスト確認**: エラー周辺のコードを理解
- **ドキュメント参照**: Go公式ドキュメント・チーム内資料を確認
- **再現性確認**: 修正後に同じエラーが発生しないか確認

#### 情報収集・質問時
- **環境情報**: OS・言語・ツールバージョンを明記
- **再現手順**: 具体的な操作手順を記録
- **エラーログ**: 完全なエラーメッセージを保存
- **試行錯誤**: 既に試した解決策を記録

## 品質ゲート（Quality Gates）

### 必須チェック項目
開発・デプロイ前に以下の項目を必ず確認してください：

#### コミット前チェック
- [ ] `go fmt ./...` - コードフォーマット適用済み
- [ ] `go vet ./...` - 静的解析警告解消済み
- [ ] `go build .` - ビルド成功
- [ ] `go test ./...` - 全テスト通過
- [ ] Git status確認 - 意図しないファイル変更なし

#### PR作成前チェック
- [ ] `go test -cover ./...` - カバレッジ確認
- [ ] セルフレビュー実施済み
- [ ] 関連ドキュメント更新済み
- [ ] テストケース追加済み（新機能・バグ修正）
- [ ] Breaking changesの文書化（該当時）

#### デプロイ前チェック
- [ ] `go build .` - ビルド成功
- [ ] 統合テスト通過
- [ ] パフォーマンス確認
- [ ] セキュリティチェック
- [ ] ロールバック手順確認

### 自動化レベル

#### 完全自動化（CI/CD）
- コードフォーマット
- リントチェック
- 型チェック
- 単体テスト実行
- ビルド検証

#### 半自動化（人間が開始）
- 統合テスト
- E2Eテスト
- セキュリティスキャン
- パフォーマンステスト

#### 手動確認必須
- コードレビュー
- アーキテクチャ設計確認
- ユーザビリティ確認
- ビジネスロジック妥当性
- データ移行影響確認

### チェックリスト運用
- **日次**: コミット前チェックを習慣化
- **週次**: 品質メトリクス確認
- **月次**: チェック項目の見直し・改善

## データファイル
- `examples/` - サンプルコード
- `testdata/` - テストデータ

## 要求仕様書
詳細な要求仕様は以下を参照：
@docs/requirements.md

## 言語・文書化規約

### 日本語使用ルール
- **Gitコミットメッセージ**: 日本語で記述
- **プルリクエスト説明文**: 日本語で記述  
- **コード内コメント**: 日本語で記述
- **変数名・関数名**: 英語（Go言語標準に従う）
- **ドキュメントファイル**: 日本語（README.md等）

### コミットメッセージ形式
```
[prefix]: [変更内容の日本語説明]

詳細な説明（必要に応じて）
- 変更点1
- 変更点2

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## プロジェクト固有の学習
プロジェクト固有の知見は`.clauderules`ファイルに記録されます。

## Memory Bank使用方針
- **通常時**: coreファイルのみ参照でコンテキスト使用量を最小化
- **詳細必要時**: contextファイルを明示的に指定
- **定期整理**: 古い情報をarchiveに移動してパフォーマンス維持

## GitHub Issue対応
- gh issueを作成・更新する際、必要なラベルは作って構いません。