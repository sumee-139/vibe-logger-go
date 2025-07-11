# Vibe Logger Go

AI駆動開発のためのGo言語向けログライブラリ。LLMが理解しやすい構造化されたログを生成し、効率的なデバッグと問題解決を支援します。

## 特徴

- **AI最適化**: LLMが理解しやすい構造化ログ出力
- **柔軟なフォーマット**: JSON、テキスト形式をサポート
- **コンテキスト対応**: 構造化されたコンテキスト情報
- **高性能**: 軽量で高速な実装

## インストール

```bash
go get github.com/sumee-139/vibe-logger-go/pkg/vibelogger
```

## 使用方法

```go
package main

import (
    "github.com/sumee-139/vibe-logger-go/pkg/vibelogger"
)

func main() {
    logger := vibelogger.New(vibelogger.Config{
        Level:      vibelogger.LevelInfo,
        Format:     vibelogger.FormatJSON,
        OutputFile: "app.log",
    })
    
    logger.Info("Application started", vibelogger.Context{
        "version": "1.0.0",
        "env":     "production",
    })
}
```

## ライセンス

MIT License