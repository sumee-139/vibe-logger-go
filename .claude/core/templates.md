# Quick Templates
tags: #templates #quick

## Quick Modes
### `/debug:start` - デバッグ特化モード
```
問題: [何が起きているか] #bug
再現手順: [手順] #reproduce
期待値: [期待される動作] #expected
実際: [実際の動作] #actual
環境: [OS/バージョン] #environment
```

### `/feature:plan` - 新機能設計モード
```
機能名: [機能名] #feature
目的: [解決したい課題] #purpose
ユーザーストーリー: [As a... I want... So that...] #story
受け入れ条件: [完了の定義] #acceptance
```

### `/review:check` - コードレビューモード
```
レビュー対象: [ファイル/機能] #review
チェック項目:
- [ ] 動作確認 #functionality
- [ ] エラーハンドリング #error  
- [ ] パフォーマンス #performance
- [ ] セキュリティ #security
- [ ] テスト #testing
改善提案: [提案内容] #improvement
```

## 基本テンプレート

### Decision Log（@.claude/context/history.mdに記録）
```
[日付] [決定内容] → [理由] #decision
```

### 学習ログ（@.claude/core/current.mdに記録）
```
技術: [学んだ技術] → [どう使えるか] #tech
ツール: [試したツール] → [評価と使用感] #tool
プロセス: [改善したプロセス] → [効果] #process
```

## よく使うパターン

### Git操作パターン
```bash
# 作業前の最新化
git pull origin main && git status

# 機能ブランチ作成
git checkout -b feature/[機能名]

# 変更の確認とコミット
git diff && git add -A && git commit -m "[prefix]: [変更内容]"

# コンフリクト解決
git stash && git pull origin main && git stash pop
```