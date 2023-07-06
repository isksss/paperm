package model

type Config struct {
	// Path to the server jar
	PaperVersion string `json:"paper_version"`
	Server       Server `json:"server"`
}

type Server struct {
	MaxMemory int `json:"max_memory"`
	MinMemory int `json:"min_memory"`
}