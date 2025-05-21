package main

import (
	"encoding/xml" // Already added
	"html/template" // For template.HTML
	"io/ioutil"    // For file reading and HTTP response body reading
	"net/http"     // For HTTP client and status codes
	"net/url"      // For URL parsing
	"strings"      // For strings.HasPrefix

	"opensocial-platform-go/pkg/game" // Added for game data and logic

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// GadgetModule represents the top-level structure of a gadget XML.
type GadgetModule struct {
	XMLName     xml.Name    `xml:"Module"`
	ModulePrefs ModulePrefs `xml:"ModulePrefs"`
	Content     Content     `xml:"Content"`
}

// ModulePrefs contains metadata about the gadget.
type ModulePrefs struct {
	XMLName xml.Name `xml:"ModulePrefs"`
	Title   string   `xml:"title,attr"`
	// Add other attributes like 'height', 'scrolling', etc. if needed.
}

// Content contains the actual gadget content, typically HTML.
type Content struct {
	XMLName xml.Name `xml:"Content"`
	Type    string   `xml:"type,attr"`
	Body    string   `xml:",cdata"` // CDATA section
}

// GameInfo struct and gamesData map are now in pkg/game/game.go

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register custom renderer
	// Ensure "web/templates" is the correct path relative to where the binary is run
	// Or use an absolute path or path relative to the executable.
	// For development, running `go run cmd/server/main.go` from the project root means "web/templates" is correct.
	renderer := NewTemplateRenderer("web/templates")
	e.Renderer = renderer

	// Static files (CSS, JS, Icons)
	// The URL path prefix, and the directory relative to where the app is run.
	e.Static("/static", "web/static") 
	// This will make files in web/static/css available under /static/css, etc.
	// So, in HTML: /static/style.css, /static/js/gadgets.js, /static/icons/sample_game_icon.png

	// Serve gadget XML files from game_data directory
	e.GET("/gadgets/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		// Security: In a real app, validate 'filename' to prevent path traversal.
		// For example, ensure it does not contain ".." or "/"
		// and that it ends with ".xml".
		// filePath := filepath.Join("game_data", filename) // Safer path construction
		filePath := "game_data/" + filename // Simpler for now as per instruction

		// Optional: Check if file exists and is of expected type.
		// _, err := os.Stat(filePath)
		// if os.IsNotExist(err) {
		// 	return echo.NewHTTPError(http.StatusNotFound, "Gadget XML not found")
		// }

		// c.File() attempts to set Content-Type based on extension.
		// For XML, it should typically set "application/xml" or "text/xml".
		return c.File(filePath)
	})

	// Routes
	e.GET("/", func(c echo.Context) error {
		data := map[string]interface{}{
			"games": game.GetAllGames(), // Use game package
		}
		return c.Render(http.StatusOK, "top.html", data)
	})

	e.GET("/game/:game_id", func(c echo.Context) error {
		gameID := c.Param("game_id")
		g, ok := game.GetGameByID(gameID) // Use game package
		if !ok {
			// Consider rendering a 404.html template if you have one
			return echo.NewHTTPError(http.StatusNotFound, "Game not found: "+gameID)
		}
		data := map[string]interface{}{
			"game_id": gameID,
			"game":    g,
		}
		return c.Render(http.StatusOK, "game_intro.html", data)
	})

	e.GET("/game/:game_id/play", func(c echo.Context) error {
		gameID := c.Param("game_id")
		g, ok := game.GetGameByID(gameID) // Use game package
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "Game not found: "+gameID)
		}
		data := map[string]interface{}{
			"game_id": gameID,
			"game":    g,
		}
		return c.Render(http.StatusOK, "game_run.html", data)
	})

	// Start server
	e.GET("/gadgets/ifr", func(c echo.Context) error {
		gadgetURLParam := c.QueryParam("url")
		if gadgetURLParam == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Gadget URL not provided.")
		}

		var xmlBytes []byte
		var err error

		// URLが相対パスか絶対パス（外部URL）か判定
		parsedGadgetURL, err := url.Parse(gadgetURLParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid gadget URL format.")
		}

		if !parsedGadgetURL.IsAbs() || (parsedGadgetURL.Scheme != "http" && parsedGadgetURL.Scheme != "https") {
			// 相対パスまたはローカルファイルシステムへのパスと見なす (例: /gadgets/sample_gadget.xml)
			if !strings.HasPrefix(gadgetURLParam, "/gadgets/") {
				 return echo.NewHTTPError(http.StatusBadRequest, "Invalid local gadget path.")
			}
			filePath := "game_data/" + strings.TrimPrefix(gadgetURLParam, "/gadgets/")
			xmlBytes, err = ioutil.ReadFile(filePath)
			if err != nil {
				c.Logger().Errorf("Failed to read local gadget XML %s: %v", filePath, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read gadget XML.")
			}
		} else {
			// 絶対URL (外部ガジェット)
			// Basic SSRF protection: only allow http and https schemes
			if parsedGadgetURL.Scheme != "http" && parsedGadgetURL.Scheme != "https" {
				c.Logger().Errorf("Unsupported scheme for external gadget XML %s", gadgetURLParam)
				return echo.NewHTTPError(http.StatusBadRequest, "Unsupported URL scheme.")
			}
			
			resp, err := http.Get(gadgetURLParam) // In production, use a client with timeout
			if err != nil {
				c.Logger().Errorf("Failed to fetch external gadget XML %s: %v", gadgetURLParam, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch gadget XML.")
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				c.Logger().Errorf("External gadget XML fetch failed %s with status %d", gadgetURLParam, resp.StatusCode)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch gadget XML (status).")
			}
			xmlBytes, err = ioutil.ReadAll(resp.Body) // In production, limit read size
			if err != nil {
				c.Logger().Errorf("Failed to read external gadget XML body %s: %v", gadgetURLParam, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read gadget XML body.")
			}
		}

		var gadget GadgetModule
		err = xml.Unmarshal(xmlBytes, &gadget)
		if err != nil {
			c.Logger().Errorf("Failed to parse gadget XML from %s: %v", gadgetURLParam, err)
			// c.Logger().Debugf("XML content for parsing error: %s", string(xmlBytes)) // Careful with large XMLs
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse gadget XML.")
		}

		if strings.ToLower(gadget.Content.Type) != "html" {
			return echo.NewHTTPError(http.StatusBadRequest, "Gadget content type must be HTML.")
		}
		
		// template.HTML を使ってHTMLエスケープを防ぐ
		htmlContent := template.HTML(gadget.Content.Body)

		data := map[string]interface{}{
			"GadgetHTMLContent": htmlContent,
		}
		return c.Render(http.StatusOK, "gadget_wrapper.html", data)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

// The hello handler is no longer needed as per instructions.
// func hello(c echo.Context) error {
// 	return c.String(http.StatusOK, "Hello, World from Echo!")
// }
