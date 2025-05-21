package game

// GameInfo は個々のゲームに関する情報を保持します。
type GameInfo struct {
	Name         string
	Description  string
	Developer    string
	Version      string
	Icon         string // 元のフルパス (例: "static/icons/sample_game_icon.png")
	IconFilename string // 表示用のファイル名 (例: "sample_game_icon.png")
	GadgetURL    string // ガジェットXMLへのパスまたはURL
}

// AllGames はプラットフォームで利用可能なすべてのゲームのマップです。
// キーはゲームID（URLで使用される識別子）です。
var AllGames = map[string]GameInfo{
	"sample_game": {
		Name:         "Sample Game",
		Description:  "これはデモンストレーション用のサンプルゲームです。",
		Developer:    "Platform Team",
		Version:      "1.0",
		Icon:         "static/icons/sample_game_icon.png",
		IconFilename: "sample_game_icon.png", // アイコンファイル名
		GadgetURL:    "/gadgets/sample_gadget.xml",
	},
	"another_game": {
		Name:         "Another Exciting Game",
		Description:  "スリリングな冒険が待っています！",
		Developer:    "Creative Studio",
		Version:      "0.9beta",
		Icon:         "static/icons/another_game_icon.png",
		IconFilename: "another_game_icon.png",
		GadgetURL:    "http://www.google.com/ig/modules/googletalk.xml", // 外部ガジェットの例
	},
	"mini_game": {
		Name:         "Mini Test Gadget",
		Description:  "コンテナテスト用の非常にシンプルなガジェット。",
		Developer:    "Platform Team",
		Version:      "0.1",
		Icon:         "", // アイコンなし
		IconFilename: "",
		GadgetURL:    "/gadgets/mini_gadget.xml",
	},
}

// GetGameByID は指定されたIDのゲーム情報を返します。
// 見つからない場合は、nil と false を返します。
func GetGameByID(id string) (*GameInfo, bool) {
	game, ok := AllGames[id]
	if !ok {
		return nil, false
	}
	return &game, true // マップから取得した値はコピーなので、ポインタを返す場合は注意が必要だが、ここではGameInfoが大きくないので値コピーで問題なし。
						// もしGameInfoが非常に大きい場合や、マップの値を直接変更したい場合はポインタのマップ map[string]*GameInfo を使う。
}

// GetAllGames はすべてのゲームのマップを返します。
func GetAllGames() map[string]GameInfo {
	return AllGames
}
