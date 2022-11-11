package common

type Config struct {
	Address     string
	CharsPath   string
	PhrasesPath string
	SqliteDB    string
}

var Conf *Config

func InitConfig() {
	if Conf == nil {
		Conf = &Config{
			// env variables
			Address:     GetOrDefault("ADDRESS", ":3000"),
			CharsPath:   GetOrDefault("CHARS_PATH", "Chars.json"),
			PhrasesPath: GetOrDefault("PHRASES_PATH", "Phrases.json"),
			SqliteDB:    GetOrDefault("SQLITE_DB", "vantu.sqlite.db"),
		}
	}
}
