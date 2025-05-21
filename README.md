# OpenSocial ゲームプラットフォーム PoC (Go言語版)

このプロジェクトは、OpenSocial APIを利用したソーシャルゲームプラットフォームの概念実証（Proof of Concept）であり、**Go言語 (Echoフレームワーク)** で実装されています。

## 主な機能

*   **TOPページ**: 利用可能なゲームの一覧を表示します。
*   **ゲーム紹介ページ**: 選択されたゲームの詳細情報を表示します。
*   **ゲーム実行ページ**: OpenSocialガジェット（ゲーム）をホストし実行します。
*   **基本的なOpenSocialコンテナ**:
    *   ガジェットXMLファイルを提供します。
    *   ガジェットが実行できるように、最小限のJavaScript環境（`gadgets.*` および `opensocial.*` のモック）を提供します。
    *   ローカルで提供されるガジェットのレンダリング、または外部URLからの取得が可能です。

## プロジェクト構成 (Go言語版)

*   `go.mod`, `go.sum`: Goモジュールの依存関係定義。
*   `cmd/server/main.go`: メインアプリケーションのエントリーポイント (Echoサーバー)。
*   `cmd/server/renderer.go`: HTMLテンプレートレンダラ。
*   `pkg/`: 共有ライブラリパッケージ。
    *   `pkg/game/game.go`: ゲームデータの定義とロジック。
    *   `pkg/config/` (将来用): 設定関連。
*   `web/`: Web関連ファイル。
    *   `web/templates/`: HTMLテンプレート (`.html`)。
    *   `web/static/`: 静的アセット (CSS, JavaScript, 画像)。
        *   `web/static/js/gadgets.js`: OpenSocial JavaScript APIのモック。
        *   `web/static/icons/`: ゲームアイコンのプレースホルダー。
*   `game_data/`: ゲームガジェットXMLファイルと関連データ。
    *   `sample_gadget.xml`: 基本機能を示すサンプルゲームガジェット。
    *   `mini_gadget.xml`: テスト用の非常にシンプルなガジェット。
*   `cmd/server/main_test.go`: `main` パッケージのテストファイル。
*   `README.md`: このファイル。

## セットアップと実行方法 (Go言語版)

1.  **Go言語のインストール**:
    Goの公式サイト (https://golang.org/) からご使用のOSに合ったGoをインストールしてください (バージョン1.18以上推奨)。

2.  **リポジトリをクローンします** (該当する場合):
    ```bash
    git clone <リポジトリURL>
    cd <リポジトリ名>
    ```

3.  **依存関係の取得**:
    Goモジュールが自動的に依存関係を解決しますが、明示的にダウンロードも可能です。
    ```bash
    go mod download
    ```

4.  **アプリケーションのビルドと実行**:
    *   **ビルド**:
        ```bash
        go build -o social-game-platform ./cmd/server
        ```
        実行可能ファイル `social-game-platform` (または `social-game-platform.exe` on Windows) が生成されます。
    *   **実行**:
        ```bash
        ./social-game-platform
        ```
    または、ビルドせずに直接実行:
    ```bash
    go run ./cmd/server/main.go ./cmd/server/renderer.go 
    # (renderer.go が main パッケージの場合。別パッケージなら不要)
    # main.go と renderer.go が同じ main パッケージなので、まとめて指定する
    ```
    アプリケーションはデフォルトで `http://localhost:1323` で利用可能になります。

5.  **テストの実行**:
    `cmd/server` ディレクトリで以下のコマンドを実行します。
    ```bash
    cd cmd/server
    go test
    ```
    リポジトリルートから特定パッケージのテストを実行する場合:
    ```bash
    go test ./cmd/server/...
    ```


## 実装されているOpenSocial機能 (モック)

*   **ガジェットレンダリング**: `type="html"` のガジェットがサポートされています。コンテンツが抽出され、iframe内にレンダリングされます。
*   **`gadgets.js` API (部分的モック)**:
    *   `gadgets.util.registerOnLoadHandler(callback)`: DOMの準備ができたときにコールバックを実行します。
    *   `gadgets.window.adjustHeight()`: コンソールに出力します。
*   **`opensocial.js` API (部分的モック)**:
    *   `opensocial.newDataRequest()`: 新しいデータリクエストオブジェクトを作成します。
    *   `request.add(opensocial.newFetchPersonRequest('VIEWER'), 'viewer')`: ビューア情報を取得するリクエストを追加します。
    *   `request.send(callback)`: リクエスト送信をシミュレートし、モックデータを返します。
    *   `Person.getDisplayName()`: ビューアのモック表示名を返します。

## 今後の開発アイデア

*   より完全なOpenSocial API実装 (永続化、アクティビティ、友達機能など)。
*   ユーザー認証と管理。
*   データベース統合でゲームやユーザーデータを永続化。
*   より堅牢なガジェットレンダリングとセキュリティ (例: GoベースのHTMLサニタイザ、コンテンツセキュリティポリシー)。
*   UI/UXの改善。
```
