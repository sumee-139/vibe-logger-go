# 💡 ベストプラクティス

## 基本的なログ設計

### 適切な操作名（Operation）の選択

```go
// ✅ Good: 明確で検索可能な操作名
logger.Info("user_login", "User authentication successful")
logger.Info("payment_process", "Payment processing started")
logger.Info("email_send", "Email notification sent")

// ❌ Bad: 曖昧で検索しにくい操作名
logger.Info("process", "Something happened")
logger.Info("event", "User did something")
logger.Info("action", "System response")
```

### 構造化されたメッセージ設計

```go
// ✅ Good: 構造化されたコンテキスト情報
logger.Info("api_call", "External API call completed",
    vibelogger.WithContext("service", "payment-gateway"),
    vibelogger.WithContext("endpoint", "/api/v2/charge"),
    vibelogger.WithContext("method", "POST"),
    vibelogger.WithContext("response_time_ms", 245),
    vibelogger.WithContext("status_code", 200))

// ❌ Bad: 情報がメッセージに埋め込まれている
logger.Info("api_call", "Called payment-gateway /api/v2/charge POST in 245ms with status 200")
```

### 相関IDの活用

```go
// ✅ Good: リクエスト全体を通しての追跡
func handleUserRequest(w http.ResponseWriter, r *http.Request) {
    correlationID := generateCorrelationID()
    
    logger.Info("request_start", "User request started",
        vibelogger.WithCorrelationID(correlationID))
    
    // ビジネスロジック内でも同じ相関IDを使用
    if err := processUserData(userID, correlationID); err != nil {
        logger.Error("process_error", "User data processing failed",
            vibelogger.WithContext("error", err.Error()),
            vibelogger.WithCorrelationID(correlationID))
        return
    }
    
    logger.Info("request_complete", "User request completed",
        vibelogger.WithCorrelationID(correlationID))
}

func processUserData(userID string, correlationID string) error {
    logger.Info("data_process", "Processing user data",
        vibelogger.WithContext("user_id", userID),
        vibelogger.WithCorrelationID(correlationID))
    
    // 処理...
    return nil
}
```

## パフォーマンス最適化

### 適切な設定値の選択

```go
// ✅ Production環境: パフォーマンス重視
func getProductionConfig() *vibelogger.LoggerConfig {
    return &vibelogger.LoggerConfig{
        MaxFileSize:     100 * 1024 * 1024, // 100MB - 大きめでI/O回数を削減
        EnableMemoryLog: false,             // メモリ使用量削減
        RotationEnabled: true,
        MaxRotatedFiles: 20,                // 十分な履歴保持
        Environment:     "production",
    }
}

// ✅ Development環境: デバッグ重視
func getDevelopmentConfig() *vibelogger.LoggerConfig {
    return &vibelogger.LoggerConfig{
        MaxFileSize:     10 * 1024 * 1024,  // 10MB - 小さめでローテーション頻度up
        EnableMemoryLog: true,              // デバッグ用
        MemoryLogLimit:  1000,              // 多めのメモリログ
        RotationEnabled: true,
        MaxRotatedFiles: 5,
        Environment:     "development",
    }
}
```

### ログレベルの適切な使い分け

```go
// ✅ Good: レベルに応じた適切な使い分け
func processOrder(order Order) {
    // INFO: 正常な業務フロー
    logger.Info("order_received", "New order received",
        vibelogger.WithContext("order_id", order.ID),
        vibelogger.WithContext("customer_id", order.CustomerID))
    
    // WARN: 注意が必要だが処理は継続
    if order.Amount > 10000 {
        logger.Warn("high_amount_order", "High amount order detected",
            vibelogger.WithContext("order_id", order.ID),
            vibelogger.WithContext("amount", order.Amount))
    }
    
    // ERROR: 処理が失敗した場合
    if err := validateOrder(order); err != nil {
        logger.Error("order_validation", "Order validation failed",
            vibelogger.WithContext("order_id", order.ID),
            vibelogger.WithContext("error", err.Error()))
        return
    }
    
    logger.Info("order_processed", "Order processing completed",
        vibelogger.WithContext("order_id", order.ID))
}
```

### バッチ処理での効率的なログ出力

```go
// ✅ Good: 効率的なバッチログ
func processBatch(items []Item) {
    batchID := generateBatchID()
    
    logger.Info("batch_start", "Batch processing started",
        vibelogger.WithContext("batch_id", batchID),
        vibelogger.WithContext("total_items", len(items)))
    
    successCount := 0
    errorCount := 0
    
    for i, item := range items {
        if err := processItem(item); err != nil {
            // エラーのみ個別ログ
            logger.Error("item_error", "Item processing failed",
                vibelogger.WithContext("batch_id", batchID),
                vibelogger.WithContext("item_id", item.ID),
                vibelogger.WithContext("error", err.Error()))
            errorCount++
        } else {
            successCount++
        }
        
        // 進捗ログは間引く（100件ごと）
        if (i+1)%100 == 0 {
            logger.Info("batch_progress", "Batch processing progress",
                vibelogger.WithContext("batch_id", batchID),
                vibelogger.WithContext("processed", i+1),
                vibelogger.WithContext("total", len(items)))
        }
    }
    
    // 最終結果をサマリーログ
    logger.Info("batch_complete", "Batch processing completed",
        vibelogger.WithContext("batch_id", batchID),
        vibelogger.WithContext("total_items", len(items)),
        vibelogger.WithContext("success_count", successCount),
        vibelogger.WithContext("error_count", errorCount))
}
```

## セキュリティベストプラクティス

### 機密情報の除外

```go
// ✅ Good: 機密情報を適切にマスキング
func logUserRegistration(user User) {
    logger.Info("user_register", "New user registration",
        vibelogger.WithContext("user_id", user.ID),
        vibelogger.WithContext("username", user.Username),
        vibelogger.WithContext("email", maskEmail(user.Email)),
        // パスワードは絶対にログに出力しない
    )
}

func maskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return "***@***"
    }
    username := parts[0]
    domain := parts[1]
    
    if len(username) <= 3 {
        return "***@" + domain
    }
    return username[:3] + "***@" + domain
}

// ✅ Good: 支払い情報の安全な記録
func logPaymentAttempt(payment Payment) {
    logger.Info("payment_attempt", "Payment processing attempt",
        vibelogger.WithContext("payment_id", payment.ID),
        vibelogger.WithContext("amount", payment.Amount),
        vibelogger.WithContext("currency", payment.Currency),
        vibelogger.WithContext("card_last4", payment.CardLast4), // 下4桁のみ
        // カード番号全体、CVV、有効期限は絶対にログに出力しない
    )
}
```

### 入力検証とログ

```go
// ✅ Good: 安全な入力値のログ記録
func logUserInput(userInput string) {
    // 危険な文字の除去/エスケープ
    safeInput := sanitizeInput(userInput)
    
    logger.Info("user_input", "User input received",
        vibelogger.WithContext("input_length", len(userInput)),
        vibelogger.WithContext("sanitized_input", safeInput))
}

func sanitizeInput(input string) string {
    // SQLインジェクション対策
    input = strings.ReplaceAll(input, "'", "\\'")
    input = strings.ReplaceAll(input, "\"", "\\\"")
    
    // XSS対策
    input = strings.ReplaceAll(input, "<", "&lt;")
    input = strings.ReplaceAll(input, ">", "&gt;")
    
    // 長すぎる入力の切り詰め
    if len(input) > 1000 {
        input = input[:1000] + "...[truncated]"
    }
    
    return input
}
```

## エラーハンドリングパターン

### 段階的エラーログ

```go
// ✅ Good: エラーの重要度に応じた段階的ログ
func processPayment(payment Payment) error {
    logger.Info("payment_start", "Payment processing started",
        vibelogger.WithContext("payment_id", payment.ID))
    
    // 一時的な問題（リトライ可能）
    if err := validatePaymentMethod(payment); err != nil {
        logger.Warn("payment_validation_warn", "Payment validation warning",
            vibelogger.WithContext("payment_id", payment.ID),
            vibelogger.WithContext("warning", err.Error()))
        // 処理継続
    }
    
    // 致命的な問題（処理停止）
    if err := chargePayment(payment); err != nil {
        logger.Error("payment_charge_error", "Payment charge failed",
            vibelogger.WithContext("payment_id", payment.ID),
            vibelogger.WithContext("error", err.Error()),
            vibelogger.WithContext("amount", payment.Amount))
        return fmt.Errorf("payment charge failed: %w", err)
    }
    
    logger.Info("payment_complete", "Payment processing completed",
        vibelogger.WithContext("payment_id", payment.ID))
    
    return nil
}
```

### エラーの文脈保持

```go
// ✅ Good: エラーチェーンを通した文脈の保持
func processUserOrder(userID string, orderData OrderData) error {
    correlationID := generateCorrelationID()
    
    logger.Info("order_process_start", "User order processing started",
        vibelogger.WithContext("user_id", userID),
        vibelogger.WithCorrelationID(correlationID))
    
    user, err := getUserProfile(userID, correlationID)
    if err != nil {
        logger.Error("user_profile_error", "Failed to get user profile",
            vibelogger.WithContext("user_id", userID),
            vibelogger.WithContext("error", err.Error()),
            vibelogger.WithCorrelationID(correlationID))
        return fmt.Errorf("user profile retrieval failed: %w", err)
    }
    
    order, err := createOrder(user, orderData, correlationID)
    if err != nil {
        logger.Error("order_creation_error", "Failed to create order",
            vibelogger.WithContext("user_id", userID),
            vibelogger.WithContext("error", err.Error()),
            vibelogger.WithCorrelationID(correlationID))
        return fmt.Errorf("order creation failed: %w", err)
    }
    
    logger.Info("order_process_complete", "User order processing completed",
        vibelogger.WithContext("user_id", userID),
        vibelogger.WithContext("order_id", order.ID),
        vibelogger.WithCorrelationID(correlationID))
    
    return nil
}
```

## マイクロサービスでの使用パターン

### サービス間通信のログ

```go
// ✅ Good: サービス間呼び出しの追跡
func callExternalService(serviceURL string, request ServiceRequest, correlationID string) (*ServiceResponse, error) {
    logger.Info("external_call_start", "External service call started",
        vibelogger.WithContext("service_url", serviceURL),
        vibelogger.WithContext("request_id", request.ID),
        vibelogger.WithCorrelationID(correlationID))
    
    start := time.Now()
    
    response, err := makeHTTPCall(serviceURL, request)
    
    duration := time.Since(start)
    
    if err != nil {
        logger.Error("external_call_error", "External service call failed",
            vibelogger.WithContext("service_url", serviceURL),
            vibelogger.WithContext("request_id", request.ID),
            vibelogger.WithContext("duration_ms", duration.Milliseconds()),
            vibelogger.WithContext("error", err.Error()),
            vibelogger.WithCorrelationID(correlationID))
        return nil, err
    }
    
    logger.Info("external_call_complete", "External service call completed",
        vibelogger.WithContext("service_url", serviceURL),
        vibelogger.WithContext("request_id", request.ID),
        vibelogger.WithContext("response_status", response.Status),
        vibelogger.WithContext("duration_ms", duration.Milliseconds()),
        vibelogger.WithCorrelationID(correlationID))
    
    return response, nil
}
```

### 分散トレーシング連携

```go
// ✅ Good: OpenTelemetryとの連携パターン
func processWithTracing(ctx context.Context, data ProcessData) error {
    span := trace.SpanFromContext(ctx)
    traceID := span.SpanContext().TraceID().String()
    spanID := span.SpanContext().SpanID().String()
    
    logger.Info("process_start", "Processing started",
        vibelogger.WithContext("trace_id", traceID),
        vibelogger.WithContext("span_id", spanID),
        vibelogger.WithContext("data_id", data.ID))
    
    if err := validateData(data); err != nil {
        span.SetStatus(codes.Error, "Data validation failed")
        logger.Error("validation_error", "Data validation failed",
            vibelogger.WithContext("trace_id", traceID),
            vibelogger.WithContext("span_id", spanID),
            vibelogger.WithContext("data_id", data.ID),
            vibelogger.WithContext("error", err.Error()))
        return err
    }
    
    result, err := processData(data)
    if err != nil {
        span.SetStatus(codes.Error, "Data processing failed")
        logger.Error("process_error", "Data processing failed",
            vibelogger.WithContext("trace_id", traceID),
            vibelogger.WithContext("span_id", spanID),
            vibelogger.WithContext("data_id", data.ID),
            vibelogger.WithContext("error", err.Error()))
        return err
    }
    
    logger.Info("process_complete", "Processing completed",
        vibelogger.WithContext("trace_id", traceID),
        vibelogger.WithContext("span_id", spanID),
        vibelogger.WithContext("data_id", data.ID),
        vibelogger.WithContext("result_id", result.ID))
    
    return nil
}
```

## テストでのベストプラクティス

### テスト用ログ設定

```go
// ✅ Good: テスト用の最適化された設定
func setupTestLogger(t *testing.T) *vibelogger.Logger {
    config := &vibelogger.LoggerConfig{
        AutoSave:        false, // ファイル出力なし
        EnableMemoryLog: true,  // メモリログでテスト検証
        MemoryLogLimit:  100,
        Environment:     "test",
    }
    
    logger := vibelogger.NewLoggerWithConfig("test-"+t.Name(), config)
    
    // テスト終了時のクリーンアップ
    t.Cleanup(func() {
        logger.Close()
    })
    
    return logger
}

func TestBusinessLogic(t *testing.T) {
    logger := setupTestLogger(t)
    
    // テスト対象の実行
    err := businessLogic(logger)
    
    // ログ出力の検証
    logs := logger.GetMemoryLogs()
    
    // 期待されるログエントリの確認
    assert.Equal(t, 3, len(logs))
    assert.Equal(t, "business_start", logs[0].Operation)
    assert.Equal(t, "INFO", logs[0].Level)
    assert.Contains(t, logs[0].Message, "started")
    
    // エラーがないことを確認
    assert.NoError(t, err)
}
```

## 運用での推奨事項

### ログローテーション戦略

```go
// ✅ Good: 適切なローテーション設定
func getRotationConfig(environment string) *vibelogger.LoggerConfig {
    switch environment {
    case "production":
        return &vibelogger.LoggerConfig{
            MaxFileSize:     100 * 1024 * 1024, // 100MB
            RotationEnabled: true,
            MaxRotatedFiles: 30,                // 1ヶ月分の履歴
        }
    case "staging":
        return &vibelogger.LoggerConfig{
            MaxFileSize:     50 * 1024 * 1024,  // 50MB
            RotationEnabled: true,
            MaxRotatedFiles: 14,                // 2週間分の履歴
        }
    default:
        return &vibelogger.LoggerConfig{
            MaxFileSize:     10 * 1024 * 1024,  // 10MB
            RotationEnabled: true,
            MaxRotatedFiles: 7,                 // 1週間分の履歴
        }
    }
}
```

### 監視・アラート連携

```go
// ✅ Good: 監視システムとの連携
func setupMonitoringLogger() *vibelogger.Logger {
    config := &vibelogger.LoggerConfig{
        ProjectName:     "monitoring",
        EnableMemoryLog: true,
        MemoryLogLimit:  1000,
    }
    
    logger, _ := vibelogger.CreateFileLoggerWithConfig("monitor", config)
    
    // 定期的なエラーログ監視
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()
        
        for range ticker.C {
            logs := logger.GetMemoryLogs()
            errorCount := 0
            
            for _, log := range logs {
                if log.Level == "ERROR" {
                    errorCount++
                }
            }
            
            // エラー率が高い場合はアラート
            if errorCount > 10 {
                sendAlert(fmt.Sprintf("High error rate detected: %d errors in last minute", errorCount))
            }
        }
    }()
    
    return logger
}
```