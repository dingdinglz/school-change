package config

// ConfigJsonModel Config模型
type ConfigJsonModel struct {
	Port   string           `json:"port"`
	SSL    bool             `json:"ssl"`
	Debug  bool             `json:"debug"`
	School string           `json:"school"`
	Limit  ConfigLimitModel `json:"limit"`
}

type ConfigLimitModel struct {
	Change int `json:"change"`
}
