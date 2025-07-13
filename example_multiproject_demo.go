package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sumee-139/vibe-logger-go/pkg/vibelogger"
)

func main() {
	// Clean up previous demo logs
	os.RemoveAll("demo_logs")

	fmt.Println("=== マルチプロジェクト対応ログ整理機能デモ ===")

	// Demo 1: デフォルトプロジェクト
	fmt.Println("\n1. デフォルトプロジェクト（プロジェクト名未指定）")
	demoDefaultProject()

	// Demo 2: 特定プロジェクト
	fmt.Println("\n2. 特定プロジェクト指定")
	demoSpecificProject()

	// Demo 3: 複数プロジェクト同時使用
	fmt.Println("\n3. 複数プロジェクト同時使用")
	demoMultipleProjects()

	// Demo 4: 環境変数での設定
	fmt.Println("\n4. 環境変数でのプロジェクト設定")
	demoEnvironmentVariable()

	// Demo 5: ローテーション機能との統合
	fmt.Println("\n5. ローテーション機能との統合")
	demoRotationIntegration()

	fmt.Println("\n=== デモ完了 ===")
	fmt.Println("生成されたログファイルを確認: ls -la demo_logs/")
}

func demoDefaultProject() {
	config := &vibelogger.LoggerConfig{
		AutoSave:        true,
		EnableMemoryLog: false,
		ProjectName:     "", // プロジェクト名を指定しない
	}

	// カスタムベースディレクトリを設定
	config.FilePath = "demo_logs/default/app.log"

	logger, err := vibelogger.CreateFileLoggerWithConfig("app", config)
	if err != nil {
		fmt.Printf("ロガー作成エラー: %v\n", err)
		return
	}
	defer logger.Close()

	fmt.Printf("  ファイルパス: %s\n", config.FilePath)

	// ログ出力
	logger.Info("startup", "アプリケーション開始（デフォルトプロジェクト）")
	logger.Warn("config", "デフォルト設定を使用")
	logger.Info("ready", "デフォルトプロジェクト準備完了")

	fmt.Println("  ✓ デフォルトプロジェクトのログ出力完了")
}

func demoSpecificProject() {
	config := &vibelogger.LoggerConfig{
		AutoSave:        true,
		EnableMemoryLog: false,
		ProjectName:     "e-commerce", // 特定プロジェクト名を指定
	}

	logger, err := vibelogger.CreateFileLoggerWithConfig("order-service", config)
	if err != nil {
		fmt.Printf("ロガー作成エラー: %v\n", err)
		return
	}
	defer logger.Close()

	fmt.Printf("  プロジェクト名: %s\n", config.ProjectName)
	fmt.Printf("  生成されるディレクトリ: logs/e-commerce/\n")

	// ログ出力
	logger.Info("service_start", "注文サービス開始")
	logger.Info("database_connect", "データベース接続成功")
	logger.Info("order_process", "注文処理システム準備完了")

	fmt.Println("  ✓ E-commerceプロジェクトのログ出力完了")
}

func demoMultipleProjects() {
	projects := []struct {
		name        string
		serviceName string
		description string
	}{
		{"user-auth", "auth-service", "認証サービス"},
		{"payment", "payment-gateway", "決済ゲートウェイ"},
		{"inventory", "stock-manager", "在庫管理システム"},
	}

	loggers := make([]*vibelogger.Logger, len(projects))

	// 各プロジェクトのロガーを作成
	for i, project := range projects {
		config := &vibelogger.LoggerConfig{
			AutoSave:        true,
			EnableMemoryLog: false,
			ProjectName:     project.name,
		}

		logger, err := vibelogger.CreateFileLoggerWithConfig(project.serviceName, config)
		if err != nil {
			fmt.Printf("  %s のロガー作成エラー: %v\n", project.name, err)
			continue
		}
		loggers[i] = logger

		fmt.Printf("  %s -> logs/%s/\n", project.description, project.name)

		// 各プロジェクト固有のログ出力
		logger.Info("startup", fmt.Sprintf("%s開始", project.description))
		logger.Info("health_check", "ヘルスチェック正常")
		logger.Info("ready", fmt.Sprintf("%s準備完了", project.description))
	}

	// クリーンアップ
	for _, logger := range loggers {
		if logger != nil {
			logger.Close()
		}
	}

	fmt.Println("  ✓ 複数プロジェクトのログ分離完了")
}

func demoEnvironmentVariable() {
	// 環境変数を設定
	os.Setenv("VIBE_LOG_PROJECT_NAME", "env-configured-project")
	defer os.Unsetenv("VIBE_LOG_PROJECT_NAME")

	// 環境変数から設定を読み込み
	config, err := vibelogger.NewConfigFromEnvironment()
	if err != nil {
		fmt.Printf("環境変数設定読み込みエラー: %v\n", err)
		return
	}

	logger, err := vibelogger.CreateFileLoggerWithConfig("env-service", config)
	if err != nil {
		fmt.Printf("ロガー作成エラー: %v\n", err)
		return
	}
	defer logger.Close()

	fmt.Printf("  環境変数 VIBE_LOG_PROJECT_NAME: %s\n", config.ProjectName)
	fmt.Printf("  自動設定されるディレクトリ: logs/env-configured-project/\n")

	// ログ出力
	logger.Info("env_config", "環境変数による設定完了")
	logger.Info("service_init", "環境設定サービス初期化")

	fmt.Println("  ✓ 環境変数でのプロジェクト設定完了")
}

func demoRotationIntegration() {
	config := &vibelogger.LoggerConfig{
		MaxFileSize:     200, // 小さなサイズでローテーションをテスト
		AutoSave:        true,
		RotationEnabled: true,
		MaxRotatedFiles: 3,
		ProjectName:     "rotation-demo", // プロジェクト名を指定
	}

	logger, err := vibelogger.CreateFileLoggerWithConfig("rotation-service", config)
	if err != nil {
		fmt.Printf("ロガー作成エラー: %v\n", err)
		return
	}
	defer logger.Close()

	fmt.Printf("  プロジェクト: %s\n", config.ProjectName)
	fmt.Printf("  ローテーションディレクトリ: logs/rotation-demo/\n")
	fmt.Printf("  最大ファイルサイズ: %d bytes\n", config.MaxFileSize)

	// ローテーションを引き起こすため大量のログを出力
	for i := 0; i < 20; i++ {
		logger.Info("bulk_operation", fmt.Sprintf("大量ログ処理 #%d - ローテーションテスト用のメッセージ", i+1))
		time.Sleep(10 * time.Millisecond) // 少し待機
	}

	// 強制ローテーション
	err = logger.ForceRotation()
	if err != nil {
		fmt.Printf("  強制ローテーションエラー: %v\n", err)
	} else {
		fmt.Println("  ✓ 強制ローテーション実行")
	}

	// ローテーションされたファイル一覧を取得
	rotatedFiles := logger.GetRotatedFiles()
	fmt.Printf("  ローテーションファイル数: %d\n", len(rotatedFiles))
	for _, file := range rotatedFiles {
		fmt.Printf("    - %s\n", file)
	}

	// ローテーション後の新しいログ
	logger.Info("post_rotation", "ローテーション後の新しいログファイル")

	fmt.Println("  ✓ プロジェクト別ローテーション機能動作確認完了")
}
