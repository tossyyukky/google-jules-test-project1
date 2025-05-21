package main

import (
	"html/template"
	"io"
	"path/filepath" // path/filepath を追加

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new TemplateRenderer
func NewTemplateRenderer(templatesDir string) *TemplateRenderer {
	// "web/templates/" ディレクトリ内のすべての "*.html" ファイルを解析
	// ParseGlob はエラーを返す可能性があるので、適切に処理する
	// 実際には起動時に一度だけ解析し、エラーなら Fatal させるのが一般的
	pattern := filepath.Join(templatesDir, "*.html") // OS非依存のパス結合
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		// 起動時にパースエラーがあれば Fatal で終了させる
		// e.Logger.Fatal() などを使うために main.go 側で初期化するか、ここで panic もありうる
		panic("Failed to parse templates: " + err.Error())
	}
	return &TemplateRenderer{
		templates: tmpl,
	}
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// data が nil の場合、空の map を渡すことでテンプレート内での `if .Data` のようなチェックがエラーにならないようにする
	if data == nil {
		data = make(map[string]interface{})
	}
	return t.templates.ExecuteTemplate(w, name, data)
}
