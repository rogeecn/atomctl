package model

type Transformer struct {
	Ignores struct {
		Jet   []string `mapstructure:"jet"`
		Model []string `mapstructure:"model"`
	} `mapstructure:"ignores"`
	Types map[string]map[string]string `mapstructure:"types"`
}
