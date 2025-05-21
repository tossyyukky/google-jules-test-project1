# OpenSocial ゲームプラットフォーム PoC (概念実証)

このプロジェクトは、OpenSocial APIを利用したソーシャルゲームプラットフォームの概念実証です。

## 主な機能

*   **TOPページ**: 利用可能なゲームの一覧を表示します。
*   **ゲーム紹介ページ**: 選択されたゲームの詳細情報を表示します。
*   **ゲーム実行ページ**: OpenSocialガジェット（ゲーム）をホストし実行します。
*   **基本的なOpenSocialコンテナ**:
    *   ガジェットXMLファイルを提供します。
    *   ガジェットが実行できるように、最小限のJavaScript環境（`gadgets.*` および `opensocial.*` のモック）を提供します。
    *   ローカルで提供されるガジェットのレンダリング、または外部URLからの取得が可能です（モックAPIの完全性に依存する制限あり）。

## プロジェクト構成

*   `app.py`: メインのFlaskアプリケーション。
*   `static/`: 静的アセット（CSS、JavaScript、画像）を格納します。
    *   `static/js/gadgets.js`: OpenSocial JavaScript APIのモック。
    *   `static/icons/`: ゲームアイコンのプレースホルダー。
*   `templates/`: FlaskアプリケーションのHTMLテンプレート。
    *   `top.html`: ゲームを一覧表示するメインページ。
    *   `game_intro.html`: ゲーム詳細ページ。
    *   `game_run.html`: ゲームガジェットのiframeをホストするページ。
    *   `gadget_wrapper.html`: `gadgets.js`とガジェットコンテンツを含むHTMLシェル。
*   `game_data/`: ゲームガジェットXMLファイルと関連データを格納します。
    *   `sample_gadget.xml`: 基本機能を示すサンプルゲームガジェット。
    *   `mini_gadget.xml`: テスト用の非常にシンプルなガジェット。
*   `tests/`: PyTestによるユニットテストおよび統合テストを格納します。
    *   `test_app.py`: Flaskアプリケーションのルートとロジックのテスト。
*   `requirements.txt`: Pythonの依存関係リスト。
*   `README.md`: このファイル。

## セットアップと実行方法

1.  **リポジトリをクローンします** (該当する場合)。

2.  **仮想環境を作成します** (推奨):
    ```bash
    python -m venv venv
    source venv/bin/activate  # Windowsの場合: venv\Scripts\activate
    ```

3.  **依存関係をインストールします**:
    ```bash
    pip install -r requirements.txt
    ```

4.  **Flaskアプリケーションを実行します**:
    ```bash
    python app.py
    ```
    アプリケーションは通常 `http://127.0.0.1:5000/` で利用可能になります。

5.  **テストの実行**:
    自動テストを実行するには、PyTestがインストールされていることを確認し（`requirements.txt`に含まれています）、次のコマンドを実行します:
    ```bash
    pytest
    ```

## 実装されているOpenSocial機能 (モック)

*   **ガジェットレンダリング**: `type="html"` のガジェットがサポートされています。コンテンツが抽出され、iframe内にレンダリングされます。
*   **`gadgets.js` API (部分的モック)**:
    *   `gadgets.util.registerOnLoadHandler(callback)`: DOMの準備ができたときにコールバックを実行します。
    *   `gadgets.window.adjustHeight()`: コンソールに出力します（`postMessage`なしでのオリジン間での完全な高さ調整は未実装）。
*   **`opensocial.js` API (部分的モック)**:
    *   `opensocial.newDataRequest()`: 新しいデータリクエストオブジェクトを作成します。
    *   `request.add(opensocial.newFetchPersonRequest('VIEWER'), 'viewer')`: ビューア情報を取得するリクエストを追加します。
    *   `request.send(callback)`: リクエスト送信をシミュレートし、モックデータを返します。
    *   `Person.getDisplayName()`: ビューアのモック表示名を返します。

## 今後の開発アイデア

*   より完全なOpenSocial API実装 (永続化、アクティビティ、友達機能など)。
*   ユーザー認証と管理。
*   ゲームとユーザーデータのデータベース統合。
*   より堅牢なガジェットレンダリングとセキュリティ (例: Cajaなどを使用したサンドボックス化)。
*   データベースからの動的なゲームリスト表示。
*   UI/UXの改善。
