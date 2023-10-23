package model

type Config struct {
	// Path to the server jar
	Url              string `json:"url"`
	Mode             string `json:"mode"`
	PaperVersion     string `json:"paper_version"`
	WaterfallVersion string `json:"waterfall_version"`
	JarName          string
	Server           Server `json:"server"`
	Plugin           Plugin `json:"plugin"`
}

type Server struct {
	MaxMemory       int      `json:"max_memory"`
	MinMemory       int      `json:"min_memory"`
	AnnounceMessage string   `json:"announce_message"`
	RestartTime     []string `json:"restart_time"`
}

type Plugin struct {
	Download bool      `json:"download"`
	Plugins  []Plugins `json:"plugins"`
}

type Plugins struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
