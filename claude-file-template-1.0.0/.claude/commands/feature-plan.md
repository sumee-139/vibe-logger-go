# 新機能設計モード

新機能企画・要件定義・設計フェーズに特化したコンテキストで作業を開始します。

## 実行手順

1. **プロジェクト全体の把握**
   - `@.claude/core/overview.md`でプロジェクト目的・ターゲット・成功基準を確認
   - 既存のコア機能・技術スタック・制約を把握

2. **計画立案**
   - `@.claude/core/next.md`で優先度マトリックスと今後のゴールを確認
   - 新機能の位置づけと優先度を評価

3. **機能設計テンプレート使用**
   - `@.claude/core/templates.md`から以下の形式で機能を整理：
     ```
     機能名: [機能名] #feature #new
     目的: [解決したい課題] #purpose
     ユーザーストーリー: [As a... I want... So that...] #story
     受け入れ条件: [完了の定義] #acceptance
     ```

4. **技術検討（必要時のみ）**
   - 技術的制約が重要な場合のみ`@.claude/context/tech.md`を確認
   - 既存アーキテクチャとの整合性を確認

5. **実装計画**
   - 機能を小さなタスクに分解
   - 依存関係・順序・見積もりを整理
   - MVP（最小実用プロダクト）と追加機能を分離

6. **計画記録**
   - 新機能タスクを`next.md`の適切な優先度に追加
   - 設計決定を`history.md`に記録

## 使用タグ
`#feature #new #planning #design #requirements #story #acceptance`

読み込むファイルはoverview.md + next.md + templatesの設計部分に限定し、効率的な機能設計を実行してください。