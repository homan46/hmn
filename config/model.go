package config

type Config struct {
	Storage struct {
		Type string `json:"type"`
		Path string `json:"path"`
	} `json:"storage"`
}
