package configs

type Supabase struct {
	URL     string `mapstructure:"url"`
	AnonKey string `mapstructure:"anon_key"`
}
