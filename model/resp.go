package model

import "time"

type Project struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
	Error       string `json:"error"`
}

type Change struct {
	Commit  string `json:"commit"`
	Summary string `json:"summary"`
	Message string `json:"message"`
}

type Downloads struct {
	Application struct {
		Name   string `json:"name"`
		Sha256 string `json:"sha256"`
	} `json:"application"`
	MojangMappings struct {
		Name   string `json:"name"`
		Sha256 string `json:"sha256"`
	} `json:"mojang-mappings"`
}

type BuildProject struct {
	ProjectID   string    `json:"project_id"`
	ProjectName string    `json:"project_name"`
	Version     string    `json:"version"`
	Build       int       `json:"build"`
	Time        time.Time `json:"time"`
	Channel     string    `json:"channel"`
	Promoted    bool      `json:"promoted"`
	Changes     []Change  `json:"changes"`
	Downloads   Downloads `json:"downloads"`
}
