# 📚 導入チュートリアル

## 🏗️ Step 1: プロジェクトセットアップ

### 1.1 基本的なプロジェクト構成

```bash
# 新しいGoプロジェクトを作成
mkdir my-project && cd my-project
go mod init my-project

# vibe-logger-goを追加
go get github.com/sumee-139/vibe-logger-go@v1.0.0
```

### 1.2 推奨ディレクトリ構造

```
my-project/
├── main.go           # メインアプリケーション
├── go.mod           # Go モジュール定義
├── go.sum           # 依存関係のチェックサム
├── config/          # 設定ファイル
├── logs/            # ログディレクトリ（自動作成）
│   ├── default/     # デフォルトプロジェクト
│   └── my-service/  # カスタムプロジェクト
└── internal/        # 内部パッケージ
```

## 🔧 Step 2: 基本的な使い方

### 2.1 最小構成での導入

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/sumee-139/vibe-logger-go"
)

func main() {
    // デフォルト設定でロガー作成
    logger, err := vibelogger.CreateFileLogger("my-service")
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer logger.Close()
    
    // 基本的なログ出力
    logger.Info("init", "Service initialization started")
    
    // エラーハンドリング例
    if err := doSomething(); err != nil {
        logger.Error("operation", "Failed to execute operation", 
            vibelogger.WithContext("error", err.Error()))
        return
    }
    
    logger.Info("shutdown", "Service shutdown completed")
}

func doSomething() error {
    // ダミーの処理
    return fmt.Errorf("simulated error")
}
```

### 2.2 Webアプリケーションでの使用例

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "github.com/sumee-139/vibe-logger-go"
)

type App struct {
    logger *vibelogger.Logger
}

func NewApp() (*App, error) {
    // Webアプリケーション用設定
    config := &vibelogger.LoggerConfig{
        ProjectName:     "web-api",
        MaxFileSize:     20 * 1024 * 1024, // 20MB
        RotationEnabled: true,
        MaxRotatedFiles: 5,
        EnableMemoryLog: true,
        MemoryLogLimit:  500, // エラー分析用
    }
    
    logger, err := vibelogger.CreateFileLoggerWithConfig("api-server", config)
    if err != nil {
        return nil, fmt.Errorf("failed to create logger: %w", err)
    }
    
    return &App{logger: logger}, nil
}

func (app *App) handleRequest(w http.ResponseWriter, r *http.Request) {
    correlationID := r.Header.Get("X-Correlation-ID")
    if correlationID == "" {
        correlationID = generateCorrelationID()
    }
    
    // リクエスト開始ログ
    app.logger.Info("request_start", "Handling HTTP request",
        vibelogger.WithContext("method", r.Method),
        vibelogger.WithContext("path", r.URL.Path),
        vibelogger.WithContext("remote_addr", r.RemoteAddr),
        vibelogger.WithCorrelationID(correlationID))
    
    start := time.Now()
    
    // ビジネスロジック処理
    response := map[string]interface{}{
        "message": "Hello, World!",
        "timestamp": time.Now().UTC(),
    }
    
    // レスポンス送信
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Correlation-ID", correlationID)
    json.NewEncoder(w).Encode(response)
    
    duration := time.Since(start)
    
    // リクエスト完了ログ
    app.logger.Info("request_complete", "HTTP request completed",
        vibelogger.WithContext("status_code", 200),
        vibelogger.WithContext("duration_ms", duration.Milliseconds()),
        vibelogger.WithCorrelationID(correlationID))
}

func generateCorrelationID() string {
    return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

func main() {
    app, err := NewApp()
    if err != nil {
        log.Fatalf("Failed to initialize app: %v", err)
    }
    defer app.logger.Close()
    
    http.HandleFunc("/", app.handleRequest)
    
    app.logger.Info("server_start", "Starting HTTP server", 
        vibelogger.WithContext("port", 8080))
    
    if err := http.ListenAndServe(":8080", nil); err != nil {
        app.logger.Error("server_error", "HTTP server failed", 
            vibelogger.WithContext("error", err.Error()))
    }
}
```

## ⚙️ Step 3: 設定のカスタマイズ

### 3.1 環境別設定管理

```go
package config

import (
    "os"
    "strconv"
    
    "github.com/sumee-139/vibe-logger-go"
)

func GetLoggerConfig() (*vibelogger.LoggerConfig, error) {
    env := getEnv("APP_ENV", "development")
    
    switch env {
    case "production":
        return getProductionConfig()
    case "staging":
        return getStagingConfig()
    default:
        return getDevelopmentConfig()
    }
}

func getProductionConfig() (*vibelogger.LoggerConfig, error) {
    return &vibelogger.LoggerConfig{
        ProjectName:     getEnv("APP_NAME", "production-app"),
        MaxFileSize:     getEnvInt64("LOG_MAX_SIZE", 100*1024*1024), // 100MB
        RotationEnabled: true,
        MaxRotatedFiles: getEnvInt("LOG_MAX_FILES", 20),
        EnableMemoryLog: false, // パフォーマンス重視
        Environment:     "production",
    }, nil
}

func getDevelopmentConfig() (*vibelogger.LoggerConfig, error) {
    return &vibelogger.LoggerConfig{
        ProjectName:     getEnv("APP_NAME", "dev-app"),
        MaxFileSize:     getEnvInt64("LOG_MAX_SIZE", 10*1024*1024), // 10MB
        RotationEnabled: true,
        MaxRotatedFiles: 5,
        EnableMemoryLog: true, // デバッグ用
        MemoryLogLimit:  1000,
        Environment:     "development",
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
    if value := os.Getenv(key); value != "" {
        if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
            return parsed
        }
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    return int(getEnvInt64(key, int64(defaultValue)))
}
```

### 3.2 構造化ロギングパターン

```go
package logging

import (
    "context"
    "fmt"
    
    "github.com/sumee-139/vibe-logger-go"
)

type StructuredLogger struct {
    logger *vibelogger.Logger
    serviceName string
    version string
}

func NewStructuredLogger(serviceName, version string) (*StructuredLogger, error) {
    config := &vibelogger.LoggerConfig{
        ProjectName: serviceName,
        MaxFileSize: 50 * 1024 * 1024,
        RotationEnabled: true,
        EnableMemoryLog: true,
        MemoryLogLimit: 500,
    }
    
    logger, err := vibelogger.CreateFileLoggerWithConfig(serviceName, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create structured logger: %w", err)
    }
    
    return &StructuredLogger{
        logger: logger,
        serviceName: serviceName,
        version: version,
    }, nil
}

// ビジネスイベントログ
func (sl *StructuredLogger) LogBusinessEvent(ctx context.Context, event string, message string, fields map[string]interface{}) {
    options := []vibelogger.LogOption{
        vibelogger.WithContext("service", sl.serviceName),
        vibelogger.WithContext("version", sl.version),
        vibelogger.WithContext("event_type", "business"),
    }
    
    // コンテキストから相関IDを取得
    if correlationID := getCorrelationIDFromContext(ctx); correlationID != "" {
        options = append(options, vibelogger.WithCorrelationID(correlationID))
    }
    
    // カスタムフィールドを追加
    for key, value := range fields {
        options = append(options, vibelogger.WithContext(key, value))
    }
    
    sl.logger.Info(event, message, options...)
}

// パフォーマンスログ
func (sl *StructuredLogger) LogPerformance(ctx context.Context, operation string, durationMs int64, fields map[string]interface{}) {
    options := []vibelogger.LogOption{
        vibelogger.WithContext("service", sl.serviceName),
        vibelogger.WithContext("operation", operation),
        vibelogger.WithContext("duration_ms", durationMs),
        vibelogger.WithContext("event_type", "performance"),
    }
    
    if correlationID := getCorrelationIDFromContext(ctx); correlationID != "" {
        options = append(options, vibelogger.WithCorrelationID(correlationID))
    }
    
    for key, value := range fields {
        options = append(options, vibelogger.WithContext(key, value))
    }
    
    if durationMs > 5000 { // 5秒以上は警告
        sl.logger.Warn(operation+"_slow", fmt.Sprintf("Slow operation detected: %dms", durationMs), options...)
    } else {
        sl.logger.Info(operation+"_perf", fmt.Sprintf("Operation completed in %dms", durationMs), options...)
    }
}

// エラーログ（スタックトレース付き）
func (sl *StructuredLogger) LogError(ctx context.Context, operation string, err error, fields map[string]interface{}) {
    options := []vibelogger.LogOption{
        vibelogger.WithContext("service", sl.serviceName),
        vibelogger.WithContext("error", err.Error()),
        vibelogger.WithContext("event_type", "error"),
    }
    
    if correlationID := getCorrelationIDFromContext(ctx); correlationID != "" {
        options = append(options, vibelogger.WithCorrelationID(correlationID))
    }
    
    for key, value := range fields {
        options = append(options, vibelogger.WithContext(key, value))
    }
    
    sl.logger.Error(operation, fmt.Sprintf("Error occurred: %v", err), options...)
}

func (sl *StructuredLogger) Close() error {
    return sl.logger.Close()
}

func getCorrelationIDFromContext(ctx context.Context) string {
    if id, ok := ctx.Value("correlation_id").(string); ok {
        return id
    }
    return ""
}
```

## 🔄 Step 4: ログローテーション

### 4.1 自動ローテーション設定

```go
config := &vibelogger.LoggerConfig{
    ProjectName:     "my-app",
    MaxFileSize:     50 * 1024 * 1024, // 50MB で自動ローテーション
    RotationEnabled: true,
    MaxRotatedFiles: 10, // 最大10個の古いファイルを保持
}

logger, err := vibelogger.CreateFileLoggerWithConfig("rotated-service", config)
```

### 4.2 手動ローテーション

```go
// 同期ローテーション（処理をブロック）
if err := logger.ForceRotation(); err != nil {
    log.Printf("Failed to rotate log: %v", err)
}

// 非同期ローテーション（処理を継続）
logger.ForceRotationAsync()
```

## 🧪 Step 5: テストでの使用

### 5.1 テスト用メモリログ

```go
func TestBusinessLogic(t *testing.T) {
    // テスト用メモリオンリー設定
    config := &vibelogger.LoggerConfig{
        AutoSave:        false, // ファイル出力なし
        EnableMemoryLog: true,
        MemoryLogLimit:  100,
    }
    
    logger := vibelogger.NewLoggerWithConfig("test", config)
    
    // テスト対象の実行
    businessLogic(logger)
    
    // ログ出力の検証
    logs := logger.GetMemoryLogs()
    assert.Equal(t, 3, len(logs))
    assert.Equal(t, "business_event", logs[0].Operation)
    assert.Equal(t, "INFO", logs[0].Level)
}

func businessLogic(logger *vibelogger.Logger) {
    logger.Info("business_event", "Event occurred")
    logger.Warn("validation", "Validation warning")
    logger.Info("completion", "Process completed")
}
```

## 📋 次のステップ

- [⚙️ 設定リファレンス](configuration.md) - 詳細な設定オプション
- [🔧 APIリファレンス](api-reference.md) - 全API関数の詳細
- [💡 ベストプラクティス](best-practices.md) - 推奨パターンとノウハウ
- [🔍 トラブルシューティング](troubleshooting.md) - よくある問題と解決方法