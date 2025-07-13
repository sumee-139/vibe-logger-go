# 🔍 トラブルシューティング

## よくある問題

### Q: ログファイルが作成されない

**症状**: ログメソッドを呼んでもファイルが作成されない

**原因と解決策**:

```go
// AutoSaveが無効になっていないか確認
config := &vibelogger.LoggerConfig{
    AutoSave: true,  // これが重要
}

// ディレクトリのアクセス権限を確認
logger, err := vibelogger.CreateFileLogger("myapp")
if err != nil {
    // エラーメッセージで権限問題を確認
    log.Printf("Logger creation failed: %v", err)
}
```

### Q: メモリログが制限より少ない

**症状**: MemoryLogLimitを100に設定したのに、50件しか取得できない

**原因と解決策**:

```go
// EnableMemoryLogが有効になっているか確認
config := &vibelogger.LoggerConfig{
    EnableMemoryLog: true,  // メモリログ有効化が必要
    MemoryLogLimit:  100,
}

// メモリログの取得方法を確認
logs := logger.GetMemoryLogs()
fmt.Printf("実際のログ数: %d\n", len(logs))
```

### Q: 環境変数設定が反映されない

**症状**: 環境変数を設定してもデフォルト値が使われる

**原因と解決策**:

```bash
# 環境変数名を確認（VIBE_LOG_プレフィックス必須）
export VIBE_LOG_MAX_FILE_SIZE=1048576
export VIBE_LOG_AUTO_SAVE=true

# 設定確認コマンド
go run -c "config, _ := vibelogger.NewConfigFromEnvironment(); fmt.Printf('%+v', config)"
```

```go
// 環境変数読み込み確認
config, err := vibelogger.NewConfigFromEnvironment()
if err != nil {
    log.Printf("Environment config error: %v", err)
}
fmt.Printf("Loaded config: %+v\n", config)
```

### Q: Path Traversal エラーが発生

**症状**: `path traversal attack detected` エラーが出る

**原因と解決策**:

```go
// 安全なパスを使用
config := &vibelogger.LoggerConfig{
    FilePath: "logs/app.log",  // OK
    // FilePath: "../../../etc/passwd",  // NG - エラーになる
}

// プロジェクト名にも注意
config := &vibelogger.LoggerConfig{
    ProjectName: "my-app",     // OK
    // ProjectName: "../admin",   // NG - エラーになる
}
```

### Q: ログローテーションが動作しない

**症状**: MaxFileSizeを超えてもローテーションされない

**原因と解決策**:

```go
// RotationEnabledが有効になっているか確認
config := &vibelogger.LoggerConfig{
    MaxFileSize:     10 * 1024 * 1024, // 10MB
    RotationEnabled: true,             // これが必要
    MaxRotatedFiles: 5,
}

// 手動でローテーションをテスト
if err := logger.ForceRotation(); err != nil {
    log.Printf("Manual rotation failed: %v", err)
}
```

### Q: メモリ使用量が多すぎる

**症状**: アプリケーションのメモリ使用量が異常に高い

**原因と解決策**:

```go
// メモリログ制限を適切に設定
config := &vibelogger.LoggerConfig{
    EnableMemoryLog: false,  // 本番環境では無効化を検討
    MemoryLogLimit:  100,    // 制限を小さく設定
}

// 定期的なメモリログクリア
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    for range ticker.C {
        logger.ClearMemoryLogs()
    }
}()
```

## パフォーマンス最適化

### 大量ログ処理の最適化

```go
// 高性能設定
config := &vibelogger.LoggerConfig{
    MaxFileSize:     100 * 1024 * 1024,  // 大きなファイルサイズ
    EnableMemoryLog: false,              // メモリログ無効
    RotationEnabled: true,
    MaxRotatedFiles: 10,                 // 適度な履歴保持
}

// バッチ処理での使用例
func processBatch(items []Item, logger *vibelogger.Logger) {
    start := time.Now()
    
    for i, item := range items {
        // 重要なログのみ出力
        if item.IsError() {
            logger.Error("batch_item", "Item processing failed",
                vibelogger.WithContext("item_id", item.ID),
                vibelogger.WithContext("error", item.Error()))
        }
        
        // 進捗ログは間引く
        if i%100 == 0 {
            logger.Info("batch_progress", fmt.Sprintf("Processed %d/%d items", i, len(items)))
        }
    }
    
    duration := time.Since(start)
    logger.Info("batch_complete", "Batch processing completed",
        vibelogger.WithContext("total_items", len(items)),
        vibelogger.WithContext("duration_ms", duration.Milliseconds()))
}
```

### 並行処理での使用

```go
// ゴルーチンセーフなロガー使用
func processWorker(id int, jobs <-chan Job, logger *vibelogger.Logger) {
    for job := range jobs {
        correlationID := fmt.Sprintf("worker-%d-job-%d", id, job.ID)
        
        logger.Info("job_start", "Processing job",
            vibelogger.WithContext("worker_id", id),
            vibelogger.WithContext("job_id", job.ID),
            vibelogger.WithCorrelationID(correlationID))
        
        // ジョブ処理
        if err := job.Process(); err != nil {
            logger.Error("job_error", "Job processing failed",
                vibelogger.WithContext("worker_id", id),
                vibelogger.WithContext("job_id", job.ID),
                vibelogger.WithContext("error", err.Error()),
                vibelogger.WithCorrelationID(correlationID))
            continue
        }
        
        logger.Info("job_complete", "Job processing completed",
            vibelogger.WithContext("worker_id", id),
            vibelogger.WithContext("job_id", job.ID),
            vibelogger.WithCorrelationID(correlationID))
    }
}
```

## セキュリティ関連

### セキュアな設定

```go
// セキュリティを重視した設定
config := &vibelogger.LoggerConfig{
    ProjectName:     "secure-app",
    FilePath:        "",  // デフォルトパスを使用（安全）
    MaxFileSize:     50 * 1024 * 1024,  // 適度なサイズ制限
    MaxRotatedFiles: 10,                // 履歴制限
    Environment:     "production",
}

// 機密情報をログに含めない
func processUser(user User, logger *vibelogger.Logger) {
    logger.Info("user_login", "User login attempt",
        vibelogger.WithContext("user_id", user.ID),
        vibelogger.WithContext("username", user.Username),
        // パスワードやトークンは絶対にログに含めない
        // vibelogger.WithContext("password", user.Password), // NG
    )
}
```

### ログの機密情報マスキング

```go
// 機密情報マスキング関数
func maskSensitiveData(data string) string {
    // クレジットカード番号のマスキング
    ccRegex := regexp.MustCompile(`\d{4}-?\d{4}-?\d{4}-?\d{4}`)
    data = ccRegex.ReplaceAllString(data, "****-****-****-****")
    
    // メールアドレスの部分マスキング
    emailRegex := regexp.MustCompile(`([a-zA-Z0-9._%+-]{1,3})[a-zA-Z0-9._%+-]*@`)
    data = emailRegex.ReplaceAllString(data, "${1}***@")
    
    return data
}

// 使用例
func logUserAction(action string, userInput string, logger *vibelogger.Logger) {
    maskedInput := maskSensitiveData(userInput)
    
    logger.Info("user_action", "User performed action",
        vibelogger.WithContext("action", action),
        vibelogger.WithContext("input", maskedInput))
}
```

## デバッグ方法

### ログ出力の確認

```go
// デバッグ用の詳細ログ
func debugLogger() {
    config := &vibelogger.LoggerConfig{
        EnableMemoryLog: true,
        MemoryLogLimit:  50,
        Environment:     "debug",
    }
    
    logger := vibelogger.NewLoggerWithConfig("debug", config)
    
    // テストログを出力
    logger.Info("debug_test", "Debug logging test")
    logger.Warn("debug_warning", "This is a warning")
    logger.Error("debug_error", "This is an error")
    
    // メモリログの確認
    logs := logger.GetMemoryLogs()
    for i, entry := range logs {
        fmt.Printf("[%d] %s: %s - %s\n", i, entry.Level, entry.Operation, entry.Message)
    }
}
```

### ファイル出力の確認

```bash
# ログファイルの存在確認
ls -la logs/

# ログファイルの内容確認（最新10行）
tail -10 logs/default/app_20250713_*.log

# JSONログの整形表示
tail -1 logs/default/app_20250713_*.log | jq .
```

### 設定の診断

```go
// 設定診断関数
func diagnoseConfig(config *vibelogger.LoggerConfig) {
    fmt.Printf("=== Logger Configuration Diagnosis ===\n")
    fmt.Printf("ProjectName: %s\n", config.ProjectName)
    fmt.Printf("MaxFileSize: %d bytes (%.2f MB)\n", 
        config.MaxFileSize, float64(config.MaxFileSize)/1024/1024)
    fmt.Printf("AutoSave: %t\n", config.AutoSave)
    fmt.Printf("EnableMemoryLog: %t\n", config.EnableMemoryLog)
    fmt.Printf("MemoryLogLimit: %d\n", config.MemoryLogLimit)
    fmt.Printf("RotationEnabled: %t\n", config.RotationEnabled)
    fmt.Printf("MaxRotatedFiles: %d\n", config.MaxRotatedFiles)
    fmt.Printf("Environment: %s\n", config.Environment)
    
    if config.FilePath != "" {
        fmt.Printf("Custom FilePath: %s\n", config.FilePath)
    }
    
    // 設定検証
    if err := config.Validate(); err != nil {
        fmt.Printf("❌ Configuration Error: %v\n", err)
    } else {
        fmt.Printf("✅ Configuration is valid\n")
    }
    fmt.Printf("=====================================\n")
}
```

## サポート情報

### ログファイルの場所

- **デフォルト**: `logs/default/`
- **プロジェクト別**: `logs/{project_name}/`
- **カスタム**: `config.FilePath` で指定した場所

### ファイル命名規則

- **通常ファイル**: `{service_name}_{YYYYMMDD}_{HHMMSS}.log`
- **ローテーションファイル**: `{service_name}_{YYYYMMDD}_{HHMMSS}.log.1`, `.2`, ...

### バージョン確認

```go
import "github.com/sumee-139/vibe-logger-go"

fmt.Printf("vibe-logger-go version: %s\n", vibelogger.GetVersion())
fmt.Printf("Version info: %+v\n", vibelogger.GetVersionInfo())
```