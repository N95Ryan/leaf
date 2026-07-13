package ui

var (
	ListShortcuts = []Shortcut{
		{Key: "n", Label: "new"},
		{Key: "r", Label: "read"},
		{Key: "e", Label: "edit"},
		{Key: "/", Label: "search"},
		{Key: "t", Label: "sort"},
		{Key: "d", Label: "delete"},
		{Key: "q", Label: "quit"},
	}

	ViewShortcuts = []Shortcut{
		{Key: "i/e", Label: "edit"},
		{Key: "Esc", Label: "back"},
	}

	EditShortcuts = []Shortcut{
		{Key: "Tab", Label: "switch field"},
		{Key: "Ctrl+S", Label: "save"},
		{Key: "Esc", Label: "cancel"},
	}

	SearchInputShortcuts = []Shortcut{
		{Key: "Enter", Label: "search"},
		{Key: "Esc", Label: "cancel"},
	}

	SearchResultsShortcuts = []Shortcut{
		{Key: "Enter", Label: "search"},
		{Key: "j/k", Label: "navigate"},
		{Key: "r", Label: "read"},
		{Key: "Esc", Label: "back"},
	}

	CreateTitleShortcuts = []Shortcut{
		{Key: "Enter", Label: "next"},
		{Key: "Esc", Label: "cancel"},
	}

	CreateContentShortcuts = []Shortcut{
		{Key: "Ctrl+S", Label: "save"},
		{Key: "Esc", Label: "back to title"},
	}
)
