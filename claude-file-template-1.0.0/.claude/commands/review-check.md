# コードレビューモード

コードレビュー・品質チェック・リファクタリングに特化したコンテキストで作業を開始します。

## 実行手順

1. **過去の知見確認**
   - `@.claude/context/history.md`から過去の決定事項・解決した問題・学びを確認
   - 似たような問題や改善点の履歴を把握

2. **レビューテンプレート使用**
   - `@.claude/core/templates.md`のコードレビューチェックリストを適用：
     ```
     レビュー対象: [ファイル/機能] #review #code
     チェック項目: [確認すべき点] #checklist
     改善提案: [提案内容] #improvement
     ```

3. **品質チェック項目**
   - **動作確認**: 機能が正しく動作するか #functionality #testing
   - **エラーハンドリング**: 例外処理が適切か #error #exception
   - **パフォーマンス**: 性能面で問題ないか #performance #optimization
   - **セキュリティ**: 脆弱性がないか #security #vulnerability
   - **テスト**: テストカバレッジは十分か #testing #coverage

4. **コードスタイル確認**
   - 既存のコーディング規約との整合性
   - 可読性・保守性の観点
   - ドキュメント・コメントの適切性

5. **改善提案**
   - 具体的で実行可能な改善案を提示
   - 優先度（Critical/High/Medium/Low）を付与
   - リファクタリングの影響範囲を評価

6. **レビュー記録**
   - レビュー結果を`current.md`に記録
   - 重要な改善点は`history.md`に記録
   - 学んだパターンは`templates.md`に追加

## 使用タグ
`#review #code #quality #refactor #improvement #checklist #standards`

読み込むファイルはhistory.md + templatesのレビュー部分に限定し、効率的なコードレビューを実行してください。