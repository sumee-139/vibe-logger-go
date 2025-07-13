# ⚙️ 設定リファレンス

## 基本設定

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

## 環境変数設定

### 環境変数で設定

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

## メモリオンリーモード

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
| `ProjectName` | `string` | `"default"` | プロジェクト名（ディレクトリ名） |
| `RotationEnabled` | `bool` | `false` | ログローテーションの有効/無効 |
| `MaxRotatedFiles` | `int` | `5` | 保持する古いログファイルの最大数 |

## 環境変数

| 環境変数 | 設定項目 | 例 |
|---------|----------|----|
| `VIBE_LOG_MAX_FILE_SIZE` | MaxFileSize | `1048576` (1MB) |
| `VIBE_LOG_AUTO_SAVE` | AutoSave | `true` / `false` |
| `VIBE_LOG_ENABLE_MEMORY` | EnableMemoryLog | `true` / `false` |
| `VIBE_LOG_MEMORY_LIMIT` | MemoryLogLimit | `500` |
| `VIBE_LOG_FILE_PATH` | FilePath | `logs/app.log` |
| `VIBE_LOG_ENVIRONMENT` | Environment | `production` |
| `VIBE_LOG_PROJECT_NAME` | ProjectName | `my-service` |
| `VIBE_LOG_ROTATION_ENABLED` | RotationEnabled | `true` / `false` |
| `VIBE_LOG_MAX_ROTATED_FILES` | MaxRotatedFiles | `10` |

## ログローテーション設定

### 基本ローテーション

```go
config := &vibelogger.LoggerConfig{
    ProjectName:     "my-app",
    MaxFileSize:     50 * 1024 * 1024, // 50MBで自動ローテーション
    RotationEnabled: true,
    MaxRotatedFiles: 10, // 最大10個の古いファイルを保持
}
```

### アドバンスドローテーション

```go
config := &vibelogger.LoggerConfig{
    ProjectName:     "high-volume-app",
    MaxFileSize:     100 * 1024 * 1024, // 100MB
    RotationEnabled: true,
    MaxRotatedFiles: 50,  // 大量の履歴を保持
    EnableMemoryLog: false, // パフォーマンス重視でメモリログ無効
}
```

## マルチプロジェクト設定

### プロジェクト別ディレクトリ

```go
// プロジェクト別設定
configs := map[string]*vibelogger.LoggerConfig{
    "auth-service": {
        ProjectName:     "auth",
        MaxFileSize:     20 * 1024 * 1024,
        RotationEnabled: true,
        MaxRotatedFiles: 15,
    },
    "payment-service": {
        ProjectName:     "payment", 
        MaxFileSize:     50 * 1024 * 1024,
        RotationEnabled: true,
        MaxRotatedFiles: 30,
        EnableMemoryLog: true, // 支払いサービスはデバッグ用メモリログ有効
    },
    "notification-service": {
        ProjectName:     "notification",
        MaxFileSize:     10 * 1024 * 1024,
        RotationEnabled: true,
        MaxRotatedFiles: 10,
    },
}

// 各サービスのロガー作成
loggers := make(map[string]*vibelogger.Logger)
for serviceName, config := range configs {
    logger, err := vibelogger.CreateFileLoggerWithConfig(serviceName, config)
    if err != nil {
        log.Fatalf("Failed to create logger for %s: %v", serviceName, err)
    }
    loggers[serviceName] = logger
}

// 使用例
loggers["auth-service"].Info("login", "User authentication successful")
loggers["payment-service"].Info("payment", "Payment processed")
loggers["notification-service"].Info("email", "Email notification sent")
```

## パフォーマンス最適化設定

### 高性能設定

```go
// 大量ログ処理時の設定
config := &vibelogger.LoggerConfig{
    MaxFileSize:     100 * 1024 * 1024,  // 100MB (大きめ)
    EnableMemoryLog: false,              // メモリログ無効でパフォーマンス向上
    AutoSave:        true,               // バッファリング有効
    RotationEnabled: true,
    MaxRotatedFiles: 20,
}
```

### 開発用設定

```go
// 開発・デバッグ用設定
config := &vibelogger.LoggerConfig{
    MaxFileSize:     10 * 1024 * 1024,   // 10MB (小さめ)
    EnableMemoryLog: true,               // デバッグ用メモリログ有効
    MemoryLogLimit:  1000,               // 大量のメモリログ
    AutoSave:        true,
    RotationEnabled: true,
    MaxRotatedFiles: 5,
    Environment:     "development",
}
```

## 設定の検証

### 設定妥当性チェック

```go
config := &vibelogger.LoggerConfig{
    MaxFileSize:     -1, // 不正な値
    MemoryLogLimit:  0,  // 不正な値  
}

// 設定検証
if err := config.Validate(); err != nil {
    log.Fatalf("Invalid configuration: %v", err)
}
```

### 設定のデフォルト補完

```go
// 部分的な設定から完全な設定を生成
partialConfig := &vibelogger.LoggerConfig{
    MaxFileSize: 50 * 1024 * 1024,
    // 他の設定は自動でデフォルト値が設定される
}

logger, err := vibelogger.CreateFileLoggerWithConfig("app", partialConfig)
if err != nil {
    log.Fatal(err)
}
```