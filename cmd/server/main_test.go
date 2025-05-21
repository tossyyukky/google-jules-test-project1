package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"opensocial-platform-go/pkg/game" // ゲームデータにアクセスするため

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert" // アサーションライブラリ
)

// setupServer はテスト用のEchoサーバーインスタンスを初期化して返します。
// main.go の main() 関数内の初期化ロジックと似ていますが、
// e.Start() は呼び出しません。
func setupServer() *echo.Echo {
	e := echo.New()
	// main.go と同じレンダラーを設定
	renderer := NewTemplateRenderer("../../web/templates") // テストファイルからの相対パス
	e.Renderer = renderer

	// main.go と同じ静的ファイル設定
	e.Static("/static", "../../web/static") // テストファイルからの相対パス

	// ハンドラを登録 (main.go からコピーまたはリファクタして共通化も検討)
	e.GET("/", func(c echo.Context) error {
		data := map[string]interface{}{
			"games": game.GetAllGames(),
		}
		return c.Render(http.StatusOK, "top.html", data)
	})

	e.GET("/game/:game_id", func(c echo.Context) error {
		gameID := c.Param("game_id")
		g, ok := game.GetGameByID(gameID)
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "Game not found")
		}
		data := map[string]interface{}{
			"game_id": gameID,
			"game":    g,
		}
		return c.Render(http.StatusOK, "game_intro.html", data)
	})
	
	e.GET("/game/:game_id/play", func(c echo.Context) error {
		gameID := c.Param("game_id")
		g, ok := game.GetGameByID(gameID)
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "Game not found")
		}
		data := map[string]interface{}{
			"game_id": gameID,
			"game":    g,
		}
		return c.Render(http.StatusOK, "game_run.html", data)
	})

	e.GET("/gadgets/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		filePath := "../../game_data/" + filename // テストファイルからの相対パス
		return c.File(filePath)
	})
	
	// /gadgets/ifr ハンドラも main.go からコピーまたは共通化
	// このテストでは簡略化のため、ifrハンドラのテストは別途詳細に行うこととし、
	// ここでは主要ページとガジェットXML提供のテストに集中します。
	// ifr のテストはXMLパースや外部リクエストのモックが必要になるため複雑。
	// (実際のプロジェクトでは `cmd/server/handlers.go` のように分離して、
	// `handlers_test.go` でテストするのが一般的)
	// For now, we'll omit the /gadgets/ifr handler setup in the test
	// to keep it focused as per the instructions. If it were included, it would be:
	// e.GET("/gadgets/ifr", originalIfrHandler) // Assuming originalIfrHandler is the actual handler from main.go


	return e
}


func TestTopPage(t *testing.T) {
	e := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "利用可能なゲーム")
	assert.Contains(t, rec.Body.String(), "Sample Game") // pkg/game/game.go のデータに基づく
}

func TestGameIntroPage_Valid(t *testing.T) {
	e := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/game/sample_game", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Sample Game - ゲーム詳細")
	assert.Contains(t, rec.Body.String(), "ゲーム開始")
}

func TestGameIntroPage_Invalid(t *testing.T) {
	e := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/game/non_existent_game", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGameRunPage_Valid(t *testing.T) {
	e := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/game/sample_game/play", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "プレイ中: Sample Game")
	// iframe の src が正しいか (URLエンコードに注意)
	// 実際の GadgetURL は "/gadgets/sample_gadget.xml"
	// Goのhtml/templateはパス内の / を %2F ではなく %2f (lowercase) にエスケープする
	expectedIframeSrc := `src="/gadgets/ifr?url=%2fgadgets%2fsample_gadget.xml"` 
	assert.Contains(t, rec.Body.String(), expectedIframeSrc)
}

func TestServeGadgetXML_Valid(t *testing.T) {
	e := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/gadgets/sample_gadget.xml", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	contentType := rec.Header().Get("Content-Type")
	// MimeType for XML can be application/xml or text/xml
	isXML := strings.HasPrefix(contentType, "application/xml") || strings.HasPrefix(contentType, "text/xml")
	assert.True(t, isXML, "Content-Type should be XML, got: "+contentType)
	assert.Contains(t, rec.Body.String(), "<ModulePrefs title=\"Sample Game Gadget\"")
}

func TestServeGadgetXML_Invalid(t *testing.T) {
	e := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/gadgets/non_existent_gadget.xml", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code) // c.File は存在しないファイルに404を返す
}

// TODO: /gadgets/ifr のテストを追加。これにはXMLのモックやHTTPクライアントのモックが必要になる場合がある。
// 例えば、ローカルXMLのパース成功ケース、パース失敗ケース、URLパラメータ欠如ケースなど。
// func TestGadgetIfr_Local_Success(t *testing.T) { ... }
// func TestGadgetIfr_MissingURL(t *testing.T) { ... }
