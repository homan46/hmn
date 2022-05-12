package config

type Config struct {
	Storage struct {
		Type string `json:"type"`
		Path string `json:"path"`
	} `json:"storage"`
	Server struct {
		UseHttps     bool     `json:"use_https"`
		TlsKey       string   `json:"tls_key"`
		TlsCert      string   `json:"tls_cert"`
		ListenOn     string   `json:"listen_on"`
		AllowOrigins []string `json:"allow_origins"`
		CookieSceret string   `json:"cookie_sceret"`
	} `json:"server"`
}
