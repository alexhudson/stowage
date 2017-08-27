package main

import "encoding/json"

// StowageMetadata is an internal type
type StowageMetadata struct {
	TypeName string `json:"type"`
	Version  int    `json:"version"`
}

func versionMeta() []byte {
	ret, _ := json.Marshal(StowageMetadata{TypeName: "stowage",
		Version: 1})

	return ret
}
