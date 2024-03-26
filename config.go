package traefik_api_key_plugin

type Config struct {
	HeaderName            string   `json:"headerName,omitempty"`
	Keys                  []string `json:"keys,omitempty"`
	RemoveHeaderOnSuccess bool     `json:"removeHeaderOnSuccess,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		HeaderName:            "X-Api-Key",
		Keys:                  make([]string, 0),
		RemoveHeaderOnSuccess: false,
	}
}
