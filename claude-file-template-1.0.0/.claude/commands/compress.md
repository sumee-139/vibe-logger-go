# Memory Bank圧縮

以下の手順でMemory Bankを軽量化してください：

## 実行手順

1. **アーカイブ対象の特定**
   - `@.claude/core/current.md`で完了済みのタスクを確認
   - `@.claude/context/history.md`で古い決定事項・問題を特定
   - 1ヶ月以上前の情報をアーカイブ候補とする

2. **アーカイブファイルの作成**
   - 完了フェーズ: `archive/phase[N]-completed-[YYYYMMDD].md`
   - 古い決定: `archive/decisions-archived-[YYYYMMDD].md`
   - 振り返り: `archive/retrospective-[YYYYMMDD].md`

3. **現在ファイルのクリーンアップ**
   - `core/current.md`: 完了タスクを削除、現在の情報のみ残す
   - `context/history.md`: 最近3ヶ月の情報のみ残す
   - `context/tech.md`: 現在使用中の技術情報のみ残す

4. **整合性チェック**
   - core/ファイル群が50行以内に収まっているか確認
   - 重要な情報がアーカイブされすぎていないか確認
   - 参照リンク（→ @path）が正しく更新されているか確認

5. **圧縮後の確認**
   - `/project:focus`コマンドで必要な情報にアクセスできるか確認
   - 日常的な作業に必要な情報が不足していないか確認

## 圧縮の目安
- core/ファイル合計: 200行以内
- context/ファイル合計: 400行以内
- 古い情報（1ヶ月以上）は積極的にアーカイブ

圧縮完了後、`core/current.md`に「[日付] Memory Bank圧縮完了」と記録してください。