# 🔧 APIリファレンス

## ロガー作成関数

### CreateFileLogger

デフォルト設定でファイルロガーを作成します。

```go
func CreateFileLogger(name string) (*Logger, error)
```

**パラメータ:**
- `name` (string): サービス名またはロガー名

**戻り値:**
- `*Logger`: 作成されたロガーインスタンス
- `error`: エラーが発生した場合

**使用例:**
```go
logger, err := vibelogger.CreateFileLogger("my-service")
if err != nil {
    log.Fatal(err)
}
defer logger.Close()
```

### CreateFileLoggerWithConfig

カスタム設定でファイルロガーを作成します。

```go
func CreateFileLoggerWithConfig(name string, config *LoggerConfig) (*Logger, error)
```

**パラメータ:**
- `name` (string): サービス名またはロガー名
- `config` (*LoggerConfig): カスタム設定

**戻り値:**
- `*Logger`: 作成されたロガーインスタンス
- `error`: エラーが発生した場合

**使用例:**
```go
config := &vibelogger.LoggerConfig{
    MaxFileSize:     50 * 1024 * 1024,
    RotationEnabled: true,
    EnableMemoryLog: true,
}

logger, err := vibelogger.CreateFileLoggerWithConfig("my-service", config)
if err != nil {
    log.Fatal(err)
}
defer logger.Close()
```

### NewLoggerWithConfig

インメモリロガーを作成します（ファイル出力なし）。

```go
func NewLoggerWithConfig(name string, config *LoggerConfig) *Logger
```

**パラメータ:**
- `name` (string): サービス名またはロガー名
- `config` (*LoggerConfig): 設定（AutoSave=falseが推奨）

**戻り値:**
- `*Logger`: 作成されたロガーインスタンス

**使用例:**
```go
config := &vibelogger.LoggerConfig{
    AutoSave:        false,
    EnableMemoryLog: true,
    MemoryLogLimit:  100,
}

logger := vibelogger.NewLoggerWithConfig("memory-logger", config)
```

## 設定関数

### DefaultConfig

デフォルト設定を取得します。

```go
func DefaultConfig() *LoggerConfig
```

**戻り値:**
- `*LoggerConfig`: デフォルト設定

**使用例:**
```go
config := vibelogger.DefaultConfig()
config.MaxFileSize = 20 * 1024 * 1024 // 20MBに変更
logger, err := vibelogger.CreateFileLoggerWithConfig("app", config)
```

### NewConfigFromEnvironment

環境変数から設定を作成します。

```go
func NewConfigFromEnvironment() (*LoggerConfig, error)
```

**戻り値:**
- `*LoggerConfig`: 環境変数から読み込んだ設定
- `error`: 環境変数の解析エラー

**使用例:**
```go
config, err := vibelogger.NewConfigFromEnvironment()
if err != nil {
    log.Fatal("Failed to load config from environment:", err)
}

logger, err := vibelogger.CreateFileLoggerWithConfig("app", config)
```

### LoggerConfig.Validate

設定の妥当性を検証します。

```go
func (c *LoggerConfig) Validate() error
```

**戻り値:**
- `error`: 設定に問題がある場合のエラー

**使用例:**
```go
config := &vibelogger.LoggerConfig{
    MaxFileSize: -1, // 不正な値
}

if err := config.Validate(); err != nil {
    log.Fatalf("Invalid configuration: %v", err)
}
```

## ログ出力メソッド

### Info

情報レベルのログを出力します。

```go
func (l *Logger) Info(operation, message string, options ...LogOption) error
```

**パラメータ:**
- `operation` (string): 操作名またはカテゴリ
- `message` (string): ログメッセージ
- `options` (...LogOption): 追加オプション

**戻り値:**
- `error`: ログ出力エラー

**使用例:**
```go
logger.Info("user_login", "User logged in successfully",
    vibelogger.WithContext("user_id", "12345"),
    vibelogger.WithCorrelationID("req-abc-123"))
```

### Warn

警告レベルのログを出力します。

```go
func (l *Logger) Warn(operation, message string, options ...LogOption) error
```

**パラメータ:**
- `operation` (string): 操作名またはカテゴリ
- `message` (string): ログメッセージ
- `options` (...LogOption): 追加オプション

**戻り値:**
- `error`: ログ出力エラー

**使用例:**
```go
logger.Warn("database_slow", "Database query took longer than expected",
    vibelogger.WithContext("query_time", "5.2s"),
    vibelogger.WithContext("query", "SELECT * FROM users"))
```

### Error

エラーレベルのログを出力します（スタックトレース付き）。

```go
func (l *Logger) Error(operation, message string, options ...LogOption) error
```

**パラメータ:**
- `operation` (string): 操作名またはカテゴリ
- `message` (string): ログメッセージ
- `options` (...LogOption): 追加オプション

**戻り値:**
- `error`: ログ出力エラー

**使用例:**
```go
logger.Error("payment_failed", "Payment processing failed",
    vibelogger.WithContext("payment_id", "pay_123"),
    vibelogger.WithContext("error", err.Error()),
    vibelogger.WithContext("amount", 99.99))
```

### LogWithContext

カスタムレベルでコンテキスト付きログを出力します。

```go
func (l *Logger) LogWithContext(level, operation, message string, context map[string]interface{}) error
```

**パラメータ:**
- `level` (string): ログレベル（"INFO", "WARN", "ERROR", "DEBUG"）
- `operation` (string): 操作名またはカテゴリ
- `message` (string): ログメッセージ
- `context` (map[string]interface{}): コンテキスト情報

**戻り値:**
- `error`: ログ出力エラー

**使用例:**
```go
context := map[string]interface{}{
    "user_id": "12345",
    "action": "profile_update",
    "ip_address": "192.168.1.100",
}

logger.LogWithContext("DEBUG", "user_action", "User updated profile", context)
```

## ログオプション

### WithContext

ログにコンテキスト情報を追加します。

```go
func WithContext(key string, value interface{}) LogOption
```

**パラメータ:**
- `key` (string): コンテキストキー
- `value` (interface{}): コンテキスト値

**使用例:**
```go
logger.Info("api_call", "External API called",
    vibelogger.WithContext("endpoint", "/api/users"),
    vibelogger.WithContext("method", "GET"),
    vibelogger.WithContext("response_time", 250))
```

### WithCorrelationID

ログに相関IDを追加します。

```go
func WithCorrelationID(correlationID string) LogOption
```

**パラメータ:**
- `correlationID` (string): 相関ID（リクエスト追跡用）

**使用例:**
```go
correlationID := generateCorrelationID()
logger.Info("request_start", "Processing user request",
    vibelogger.WithCorrelationID(correlationID))
```

## メモリログメソッド

### GetMemoryLogs

メモリに保存されているログエントリを取得します。

```go
func (l *Logger) GetMemoryLogs() []LogEntry
```

**戻り値:**
- `[]LogEntry`: メモリログエントリのスライス

**使用例:**
```go
logs := logger.GetMemoryLogs()
for i, entry := range logs {
    fmt.Printf("[%d] %s: %s - %s\n", i, entry.Level, entry.Operation, entry.Message)
}
```

### ClearMemoryLogs

メモリログをクリアします。

```go
func (l *Logger) ClearMemoryLogs()
```

**使用例:**
```go
// 定期的なメモリログクリア
go func() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    for range ticker.C {
        logger.ClearMemoryLogs()
    }
}()
```

## ローテーションメソッド

### ForceRotation

ログファイルを強制的にローテーションします（同期）。

```go
func (l *Logger) ForceRotation() error
```

**戻り値:**
- `error`: ローテーションエラー

**使用例:**
```go
// 同期ローテーション
if err := logger.ForceRotation(); err != nil {
    log.Printf("Failed to rotate log: %v", err)
}
```

### ForceRotationAsync

ログファイルを非同期でローテーションします。

```go
func (l *Logger) ForceRotationAsync()
```

**使用例:**
```go
// 非同期ローテーション（処理をブロックしない）
logger.ForceRotationAsync()
```

## リソース管理

### Close

ロガーを終了し、リソースを解放します。

```go
func (l *Logger) Close() error
```

**戻り値:**
- `error`: 終了処理エラー

**使用例:**
```go
logger, err := vibelogger.CreateFileLogger("app")
if err != nil {
    log.Fatal(err)
}
defer logger.Close() // 必ず呼び出す
```

## バージョン情報

### GetVersion

ライブラリのバージョン文字列を取得します。

```go
func GetVersion() string
```

**戻り値:**
- `string`: バージョン文字列（例: "1.0.0"）

**使用例:**
```go
fmt.Printf("vibe-logger-go version: %s\n", vibelogger.GetVersion())
```

### GetVersionInfo

詳細なバージョン情報を取得します。

```go
func GetVersionInfo() VersionInfo
```

**戻り値:**
- `VersionInfo`: 詳細なバージョン情報

**使用例:**
```go
info := vibelogger.GetVersionInfo()
fmt.Printf("Version: %s\n", info.Version)
fmt.Printf("Go Version: %s\n", info.GoVersion)
fmt.Printf("User Agent: %s\n", info.UserAgent)
```

### IsStableVersion

安定版かどうかを確認します。

```go
func IsStableVersion() bool
```

**戻り値:**
- `bool`: 安定版の場合true

**使用例:**
```go
if vibelogger.IsStableVersion() {
    fmt.Println("Using stable version")
} else {
    fmt.Println("Using pre-release version")
}
```

### CompareVersion

バージョンを比較します。

```go
func CompareVersion(version string) int
```

**パラメータ:**
- `version` (string): 比較対象のバージョン

**戻り値:**
- `int`: 比較結果（-1: 小さい, 0: 同じ, 1: 大きい）

**使用例:**
```go
result := vibelogger.CompareVersion("0.9.0")
if result > 0 {
    fmt.Println("Current version is newer than 0.9.0")
}
```

## データ構造

### LoggerConfig

ロガーの設定を定義する構造体。

```go
type LoggerConfig struct {
    MaxFileSize     int64  `json:"max_file_size"`
    AutoSave        bool   `json:"auto_save"`
    EnableMemoryLog bool   `json:"enable_memory_log"`
    MemoryLogLimit  int    `json:"memory_log_limit"`
    FilePath        string `json:"file_path"`
    Environment     string `json:"environment"`
    ProjectName     string `json:"project_name"`
    RotationEnabled bool   `json:"rotation_enabled"`
    MaxRotatedFiles int    `json:"max_rotated_files"`
}
```

### LogEntry

ログエントリを表す構造体。

```go
type LogEntry struct {
    Timestamp   time.Time              `json:"timestamp"`
    Level       string                 `json:"level"`
    Operation   string                 `json:"operation"`
    Message     string                 `json:"message"`
    Context     map[string]interface{} `json:"context,omitempty"`
    Environment EnvironmentInfo        `json:"environment"`
    // その他のフィールド...
}
```

### VersionInfo

バージョン情報を表す構造体。

```go
type VersionInfo struct {
    Version    string `json:"version"`
    Major      int    `json:"major"`
    Minor      int    `json:"minor"`
    Patch      int    `json:"patch"`
    GoVersion  string `json:"go_version"`
    UserAgent  string `json:"user_agent"`
}
```

## エラーハンドリング

### よくあるエラー

- `ErrInvalidConfig`: 設定が不正
- `ErrFilePermission`: ファイル権限エラー
- `ErrPathTraversal`: パストラバーサル攻撃検出
- `ErrRotationFailed`: ローテーション失敗

**使用例:**
```go
logger, err := vibelogger.CreateFileLogger("app")
if err != nil {
    switch {
    case errors.Is(err, vibelogger.ErrInvalidConfig):
        log.Fatal("Configuration is invalid")
    case errors.Is(err, vibelogger.ErrFilePermission):
        log.Fatal("File permission denied")
    default:
        log.Fatalf("Unexpected error: %v", err)
    }
}
```