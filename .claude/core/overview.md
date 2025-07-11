# Project Overview
tags: #overview #project #summary

## 3行サマリー
- **目的**: AI駆動開発向けのGo言語ログライブラリ開発 #purpose #problem
- **ターゲット**: LLMを活用するGo開発者・デバッグを効率化したい開発者 #target #users
- **成功基準**: 構造化ログでデバッグ時間を50%削減 #success #metrics

## プロジェクト基本情報
- **開始日**: 2025-07-11
- **期限**: 2025-10-11 (3ヶ月)
- **現在の進捗**: 15% (基本実装完了)

## コア機能（優先順）
1. **構造化ログ**: JSON/テキスト形式でコンテキスト情報を記録 #feature #core #p1
2. **AI最適化**: LLMが理解しやすいログフォーマット #feature #core #p2
3. **パフォーマンス**: 高速で軽量な実装 #feature #core #p3
4. **拡張性**: カスタムフォーマッター・出力先対応 #feature #core #p4

## 技術スタック
- **言語**: Go 1.19+ #tech #language
- **フレームワーク**: 標準ライブラリベース #tech #framework
- **DB**: ファイルベースログ #tech #database
- **その他**: JSON encoding, context package #tech #tools

## 制約・前提
- Go標準ライブラリを優先使用
- 外部依存を最小限に抑制
- 後方互換性を重視

## 成功基準 & 進捗
- [x] 基本ログ機能実装 - 完了
- [ ] AI最適化フォーマッター - 着手予定
- [ ] パフォーマンステスト - 未着手

## 問題定義
従来のログライブラリはLLMが解析しにくく、AI駆動開発での効率が悪い

## ユーザー価値
AI支援によるデバッグ時間短縮とコード品質向上

詳細な技術情報 → @.claude/context/tech.md
詳細な仕様 → @docs/requirements.md