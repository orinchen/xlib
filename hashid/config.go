package hashid

import (
	"github.com/sqids/sqids-go"
	"log"
)

type Config struct {
	MinLength uint   `json:"minLength,optional" toml:"minLength"`
	Alphabet  string `json:"alphabet,optional" toml:"alphabet"`
}

func DefaultConfig() *Config {
	return &Config{
		MinLength: 8,
	}
}

// Build ...
func (config Config) Build() *sqids.Sqids {

	options := sqids.Options{
		MinLength: uint8(config.MinLength),
	}

	if config.Alphabet != "" {
		options.Alphabet = config.Alphabet
	}

	sqid, err := sqids.New(options)

	if err != nil {
		log.Panicln("Hashid 初始化失败", err)
	}
	return sqid
}
