// client configuration module provides app-wide configuration structure with easy init
package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"internal/adapters/cryptor"
	"internal/domain"
)

var Cc domain.ClientConfig

// NewConfig initializes a Config with default values
func NewClientConfig() domain.ClientConfig {
	return domain.ClientConfig{
		ConfigFilePath: "",
		Endpoint:       domain.ENDPOINT,
		CryptoKeyPath:  "",
		LogLevel:       domain.LOGLEVEL,
	}
}

// custom command line parser to read config file name before flag.Parse() -- iter22 requirement
func ParseAgentConfigFile(cf *domain.ClientConfig) {
	for i, arg := range os.Args {
		if arg == "-config" {
			if i+1 < len(os.Args) {
				cf.ConfigFilePath = strings.TrimSpace(os.Args[i+1])
			}
		}
	}
	if val, found := os.LookupEnv("CONFIG"); found {
		cf.ConfigFilePath = val
	}

	if cf.ConfigFilePath == "" {
		return
	}

	jcf := NewClientConfig()

	file, err := os.Open(cf.ConfigFilePath)
	if err != nil {
		panic(fmt.Sprintf("PANIC: error reading config file: %s", err.Error()))
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&jcf); err != nil {
		panic(fmt.Sprintf("PANIC: error decoding JSON config: %s", err.Error()))
	}

	cf.Endpoint = jcf.Endpoint
	cf.RateLimit = jcf.RateLimit
	cf.CryptoKeyPath = jcf.CryptoKeyPath
	cf.LogLevel = jcf.LogLevel
}

// set agent configuration using command line arguments and/or environment variables
func InitClientConfig() domain.ClientConfig {

	cf := NewClientConfig()
	cf.UseCompression = true    // activate gzip compression
	cf.MaxConnectionRetries = 3 // Connection retries for retriable functions (does not include original request. 0 to disable)
	cf.ConfigFilePath = ""

	//load config from command line or env variable with lowest priority
	ParseAgentConfigFile(&cf)

	//set defaults and read command line
	flag.StringVar(&cf.ConfigFilePath, "config", cf.ConfigFilePath, "path to configuration file in JSON format") //used to pass Parse() check
	flag.StringVar(&cf.Endpoint, "a", cf.Endpoint, "the address:port server endpoint to send metric data")
	flag.StringVar(&cf.CryptoKeyPath, "crypto-key", cf.CryptoKeyPath, "path to public crypto key")
	flag.StringVar(&cf.LogLevel, "v", cf.LogLevel, "log verbosity (log level)")
	flag.Parse()

	//parse env variables
	if val, found := os.LookupEnv("ADDRESS"); found {
		cf.Endpoint = val
	}
	if val, found := os.LookupEnv("CRYPTO_KEY"); found {
		cf.CryptoKeyPath = val
	}
	if val, found := os.LookupEnv("LOG_LEVEL"); found {
		cf.LogLevel = val
	}

	// check for critical missing config entries

	if cf.Endpoint == "" {
		panic("PANIC: endpoint address:port is not set")
	}
	if cf.LogLevel == "" {
		panic("PANIC: log level is not set")
	}

	//set encryption logic
	cf.CryptoKeyPath = strings.TrimSpace(cf.CryptoKeyPath)
	if cf.CryptoKeyPath != "" {
		err := cryptor.LoadPublicKey(cf.CryptoKeyPath)
		if err != nil {
			panic("PANIC: failed to load crypto key " + err.Error())
		}
		cryptor.EnableEncryption(true)
	}

	// rate limiter global en\disable
	cf.UseRateLimit = (cf.RateLimit > 0)

	return cf
}
