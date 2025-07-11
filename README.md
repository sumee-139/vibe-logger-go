# Vibe Logger Go

AI駆動開発のためのGo言語向けログライブラリ。LLMが理解しやすい構造化されたログを生成し、効率的なデバッグと問題解決を支援します。

## 特徴

- **AI最適化**: LLMが理解しやすい構造化ログ出力
- **柔軟な設定**: 環境変数、カスタム設定、デフォルト設定をサポート
- **メモリログ**: インメモリログ機能でリアルタイム分析
- **セキュリティ**: Path Traversal防止、入力検証、リソース制限
- **高性能**: 軽量で高速、スレッドセーフな実装

## インストール

```bash
go get github.com/sumee-139/vibe-logger-go/pkg/vibelogger
```

## クイックスタート

### 基本的な使用方法

```go
package main

import (
    "log"
    "github.com/sumee-139/vibe-logger-go/pkg/vibelogger"
)

func main() {
    // デフォルト設定でファイルロガー作成
    logger, err := vibelogger.CreateFileLogger("myapp")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()
    
    // 各種ログレベルでメッセージ出力
    logger.Info("app_start", "Application started successfully")
    logger.Warn("config_missing", "Using default configuration")
    logger.Error("db_connection", "Failed to connect to database")
}
```

## Configuration System

### デフォルト設定

```go
// デフォルト設定で簡単開始
logger, err := vibelogger.CreateFileLogger("myapp")
if err != nil {
    log.Fatal(err)
}
defer logger.Close()

logger.Info("operation", "Using default configuration")
```

### カスタム設定

```go
// カスタム設定でロガー作成
config := &vibelogger.LoggerConfig{
    MaxFileSize:     1024 * 1024,  // 1MB制限
    AutoSave:        true,         // 自動ファイル保存
    EnableMemoryLog: true,         // メモリログ有効
    MemoryLogLimit:  100,          // メモリログ100件制限
    FilePath:        "logs/custom.log",
    Environment:     "production",
}

logger, err := vibelogger.CreateFileLoggerWithConfig("myapp", config)
if err != nil {
    log.Fatal(err)
}
defer logger.Close()

logger.Info("operation", "Using custom configuration")

// メモリログの取得
memoryLogs := logger.GetMemoryLogs()
fmt.Printf("Memory logs: %d entries\n", len(memoryLogs))
```

### 環境変数設定

```bash
# 環境変数で設定
export VIBE_LOG_MAX_FILE_SIZE=2097152      # 2MB
export VIBE_LOG_AUTO_SAVE=true
export VIBE_LOG_ENABLE_MEMORY=true
export VIBE_LOG_MEMORY_LIMIT=50
export VIBE_LOG_FILE_PATH=logs/app.log
export VIBE_LOG_ENVIRONMENT=production
```

```go
// 環境変数から設定読み込み
config, err := vibelogger.NewConfigFromEnvironment()
if err != nil {
    log.Fatal("Failed to load config from environment:", err)
}

logger, err := vibelogger.CreateFileLoggerWithConfig("myapp", config)
if err != nil {
    log.Fatal(err)
}
defer logger.Close()

logger.Info("operation", "Using environment configuration")
```

### メモリオンリーモード

```go
// ファイル出力なし、メモリログのみ
config := &vibelogger.LoggerConfig{
    AutoSave:        false,  // ファイル出力無効
    EnableMemoryLog: true,   // メモリログ有効
    MemoryLogLimit:  200,    // 200件制限
}

logger := vibelogger.NewLoggerWithConfig("memory-only", config)

// ログ出力
for i := 1; i <= 10; i++ {
    logger.Info("memory_test", fmt.Sprintf("Log entry %d", i))
}

// メモリログの取得と表示
logs := logger.GetMemoryLogs()
for i, entry := range logs {
    fmt.Printf("%d: %s - %s\n", i+1, entry.Operation, entry.Message)
}

// メモリログクリア
logger.ClearMemoryLogs()
```

## 設定オプション

| オプション | 型 | デフォルト | 説明 |
|-----------|----|-----------|----|
| `MaxFileSize` | `int64` | `10485760` (10MB) | ファイルサイズ上限（バイト）、0は無制限 |
| `AutoSave` | `bool` | `true` | ファイル自動保存の有効/無効 |
| `EnableMemoryLog` | `bool` | `false` | インメモリログの有効/無効 |
| `MemoryLogLimit` | `int` | `1000` | メモリログの最大エントリ数 |
| `FilePath` | `string` | `""` | カスタムログファイルパス |
| `Environment` | `string` | `"development"` | 環境名（dev/prod/test等） |

## 環境変数

| 環境変数 | 設定項目 | 例 |
|---------|----------|---|
| `VIBE_LOG_MAX_FILE_SIZE` | MaxFileSize | `1048576` (1MB) |
| `VIBE_LOG_AUTO_SAVE` | AutoSave | `true` / `false` |
| `VIBE_LOG_ENABLE_MEMORY` | EnableMemoryLog | `true` / `false` |
| `VIBE_LOG_MEMORY_LIMIT` | MemoryLogLimit | `500` |
| `VIBE_LOG_FILE_PATH` | FilePath | `logs/app.log` |
| `VIBE_LOG_ENVIRONMENT` | Environment | `production` |

## ログ出力形式

```json
{
  "timestamp": "2025-07-11T22:19:04.521669872Z",
  "level": "INFO",
  "operation": "user_login",
  "message": "User logged in successfully",
  "environment": {
    "arch": "amd64",
    "go_version": "go1.21.5",
    "os": "linux",
    "pid": "1234",
    "pwd": "/app"
  },
  "severity": 2,
  "category": "security",
  "searchable": "user_login user_login User logged in successfully",
  "pattern": "unknown_pattern"
}
```

## セキュリティ機能

### Path Traversal防止
```go
// 危険なパスは自動的にブロック
config := &vibelogger.LoggerConfig{
    FilePath: "../../../etc/passwd",  // エラーになる
}
```

### リソース制限
- **ファイルサイズ**: 最大1GB制限
- **メモリログ**: 最大10,000エントリ制限
- **パス長**: 最大255文字制限

### 入力検証
- 環境変数の厳格な検証
- 不正な文字・値の自動拒否
- セキュアなデフォルト値

## APIリファレンス

### ロガー作成関数

```go
// デフォルト設定でファイルロガー作成
func CreateFileLogger(name string) (*Logger, error)

// カスタム設定でファイルロガー作成
func CreateFileLoggerWithConfig(name string, config *LoggerConfig) (*Logger, error)

// インメモリロガー作成（ファイル出力なし）
func NewLoggerWithConfig(name string, config *LoggerConfig) *Logger
```

### 設定関数

```go
// デフォルト設定取得
func DefaultConfig() *LoggerConfig

// 環境変数から設定作成
func NewConfigFromEnvironment() (*LoggerConfig, error)

// 設定検証
func (c *LoggerConfig) Validate() error
```

### ログ出力メソッド

```go
// 情報ログ
func (l *Logger) Info(operation, message string) error

// 警告ログ
func (l *Logger) Warn(operation, message string) error

// エラーログ（スタックトレース付き）
func (l *Logger) Error(operation, message string) error

// コンテキスト付きログ
func (l *Logger) LogWithContext(level, operation, message string, context map[string]interface{}) error
```

### メモリログメソッド

```go
// メモリログ取得
func (l *Logger) GetMemoryLogs() []LogEntry

// メモリログクリア
func (l *Logger) ClearMemoryLogs()

// ロガー終了
func (l *Logger) Close() error
```

## トラブルシューティング

### よくある問題

**Q: ログファイルが作成されない**
```go
// AutoSaveが無効になっていないか確認
config := &vibelogger.LoggerConfig{
    AutoSave: true,  // これが重要
}
```

**Q: メモリログが制限より少ない**
```go
// EnableMemoryLogが有効になっているか確認
config := &vibelogger.LoggerConfig{
    EnableMemoryLog: true,  // メモリログ有効化が必要
    MemoryLogLimit:  100,
}
```

**Q: 環境変数設定が反映されない**
```bash
# 環境変数名を確認（VIBE_LOG_プレフィックス必須）
export VIBE_LOG_MAX_FILE_SIZE=1048576
export VIBE_LOG_AUTO_SAVE=true

# 設定確認
go run -c "config, _ := vibelogger.NewConfigFromEnvironment(); fmt.Printf('%+v', config)"
```

**Q: Path Traversal エラーが発生**
```go
// 安全なパスを使用
config := &vibelogger.LoggerConfig{
    FilePath: "logs/app.log",  // OK
    // FilePath: "../../../etc/passwd",  // NG - エラーになる
}
```

### パフォーマンス最適化

```go
// 大量ログ処理時の設定
config := &vibelogger.LoggerConfig{
    MaxFileSize:     100 * 1024 * 1024,  // 100MB (大きめ)
    EnableMemoryLog: false,              // メモリログ無効でパフォーマンス向上
    AutoSave:        true,               // バッファリング有効
}
```

## デモとサンプル

### 設定デモ実行
```bash
go run example_config_demo.go
```

### AIデモ実行  
```bash
go run example_ai_demo.go
```

### テスト実行
```bash
go test ./pkg/vibelogger/
```

## ライセンス

MIT License