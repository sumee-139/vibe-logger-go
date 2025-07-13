# ⚡ クイックスタート

## 🚀 30秒で始める

```go
package main

import (
    "log"
    "github.com/sumee-139/vibe-logger-go"
)

func main() {
    // 1行でロガー作成（自動的に logs/ ディレクトリに保存）
    logger, err := vibelogger.CreateFileLogger("myapp")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()
    
    // 構造化ログを出力
    logger.Info("startup", "Application started successfully")
    logger.Warn("config", "Using default configuration")
    logger.Error("database", "Connection failed", 
        vibelogger.WithContext("host", "localhost:5432"),
        vibelogger.WithContext("timeout", "30s"))
}
```

**これだけで以下が自動的に設定されます:**
- ✅ ファイルへの構造化ログ出力（JSON形式）
- ✅ タイムスタンプと環境情報の自動付与
- ✅ AI最適化されたフィールド（重要度、カテゴリ、検索タグ）
- ✅ スレッドセーフな並行処理対応

## 🎯 5分で応用

```go
// プロジェクト別ログ管理
config := &vibelogger.LoggerConfig{
    ProjectName:     "e-commerce",      // logs/e-commerce/ に保存
    MaxFileSize:     50 * 1024 * 1024,  // 50MB制限
    RotationEnabled: true,              // 自動ローテーション
    MaxRotatedFiles: 10,                // 10ファイルまで保持
}

logger, err := vibelogger.CreateFileLoggerWithConfig("order-service", config)
if err != nil {
    log.Fatal(err)
}
defer logger.Close()

// リッチなコンテキスト情報付きログ
logger.Info("order_created", "New order received",
    vibelogger.WithContext("order_id", "12345"),
    vibelogger.WithContext("user_id", "user789"),
    vibelogger.WithContext("amount", 99.99),
    vibelogger.WithCorrelationID("req-abc-123"))
```

## 📋 次のステップ

- [📚 導入チュートリアル](tutorial.md) - 基本からアドバンスドな使い方まで
- [⚙️ 設定リファレンス](configuration.md) - 詳細な設定オプション
- [🔧 APIリファレンス](api-reference.md) - 全API関数の詳細
- [💡 ベストプラクティス](best-practices.md) - 推奨パターンとノウハウ