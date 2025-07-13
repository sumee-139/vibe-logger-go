# Vibe Logger Go

[![Go Version](https://img.shields.io/badge/go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/sumee-139/vibe-logger-go.svg)](https://github.com/sumee-139/vibe-logger-go/releases)

AI駆動開発のためのGo言語向けログライブラリ。LLMが理解しやすい構造化されたログを生成し、効率的なデバッグと問題解決を支援します。

## 🚀 特徴

- **🤖 AI最適化**: LLMが理解しやすい構造化ログ出力
- **🗂️ マルチプロジェクト対応**: プロジェクト別のディレクトリ構造
- **🔄 ログローテーション**: 自動ファイルローテーションとクリーンアップ
- **⚙️ 柔軟な設定**: 環境変数、カスタム設定、デフォルト設定をサポート
- **💾 メモリログ**: インメモリログ機能でリアルタイム分析
- **🔒 セキュリティ**: Path Traversal防止、入力検証、リソース制限
- **⚡ 高性能**: 軽量で高速、スレッドセーフな実装

## 📦 インストール

```bash
go get github.com/sumee-139/vibe-logger-go@v1.0.0
```

**必要なGo バージョン**: 1.19+

## ⚡ クイックスタート

[📋 完全なクイックスタートガイド](docs/quickstart.md)をご確認ください。

### 🚀 30秒で始める

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

### 🎯 5分で応用

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

## 📚 ドキュメント

### 基本ガイド
- [📋 クイックスタート](docs/quickstart.md) - 30秒で始める基本的な使い方
- [📚 導入チュートリアル](docs/tutorial.md) - 基本からアドバンスドな使い方まで
- [⚙️ 設定リファレンス](docs/configuration.md) - 詳細な設定オプション

### 詳細リファレンス
- [🔧 APIリファレンス](docs/api-reference.md) - 全API関数の詳細
- [💡 ベストプラクティス](docs/best-practices.md) - 推奨パターンとノウハウ
- [🔍 トラブルシューティング](docs/troubleshooting.md) - よくある問題と解決方法

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

## デモとサンプル

### 設定デモ実行
```bash
go run cmd/examples/example_config_demo.go
```

### AIデモ実行  
```bash
go run cmd/examples/example_ai_demo.go
```

### テスト実行
```bash
go test .
```

## バージョン情報

```go
import "github.com/sumee-139/vibe-logger-go"

fmt.Printf("Version: %s\n", vibelogger.GetVersion())
fmt.Printf("Version Info: %+v\n", vibelogger.GetVersionInfo())
```

## ライセンス

MIT License