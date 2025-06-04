# 📅 Today-in-History Bot （今日はなんの日？ボット）

Go 製ワンバイナリで動く自動投稿ツールです。  
**毎週日曜日** に翌週 7 日分 × 3 ソース（計 **21 件**）の  
「今日はなんの日？」を生成し、Discord の指定チャンネルへ投稿します。

---

## ✨ 主な機能

* **三つの情報源**  
  1. **静的 CSV**: リポジトリに同梱した一覧から日付一致の記念日を抽出
  2. **ChatGPT**: OpenAI API に「この日は何の日？」を問い合わせて生成  
  3. **Perplexity**: Perplexity API で同様に生成  
* **ワンバイナリ設計** – `go build` だけで実行ファイルを生成  
* **Discord Webhook 投稿** – Bot アカウント不要  
* **GitHub Actions 対応** – 毎週日曜 09:00 JST に自動実行  
* **ドライランモード** – 送信せず内容をコンソールへ出力可能  
* **簡単拡張** – provider インターフェースを実装すれば別ソースを追加可  

---

## 📂 ディレクトリ構成

```text
.
├── main.go                    # ここだけビルドすれば OK
├── internal/
│   ├── provider/              # static / chatgpt / perplexity
│   └── model/
├── data/
│   └── anniversaries.csv      # 静的リスト（UTF-8, カンマ区切り）
├── .github/workflows/post.yml # 定期実行ワークフロー
└── README.md
