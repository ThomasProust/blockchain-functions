package env

import (
	"context"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/spf13/viper"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// Config structure of the app configuration
type Config struct {
	ProjectID string `mapstructure:"project_id"`
	KeyPath   string `mapstructure:"keyPath"`
	Ethereum  ChainConfig
	Bitcoin   ChainConfig
}

// ChainConfig configuration of each chain with its name and the supported currencies
type ChainConfig struct {
	Chain         string `mapstructure:"chain"`
	Endpoint      string `mapstructure:"endpoint,omitempty"`
	Confirmations int    `mapstructure:"confirmations"`
	GasStation    string `mapstructure:"gas_station,omitempty"`
	Currencies    []*CurrencyConfig
}

// CurrencyConfig structure of the configuration of each supported currency
type CurrencyConfig struct {
	Name              string `mapstructure:"name"`
	Decimals          int    `mapstructure:"decimals"`
	Address           string `mapstructure:"address,omitempty"`
	TransferSignature string `mapstructure:"transfer_signature,omitempty"`
}

// Constants for project ids
const (
	DEVELOP    string = "black-stream-292507"
	PRODUCTION string = "soteria-production"
)

// InitEnvVars initialize env variables
func InitConfig() *Config {

	gcpProject := os.Getenv("GCP_PROJECT")
	keyPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	file := "staging"
	if gcpProject == PRODUCTION {
		file = "prod"
	}
	config, errConf := loadConfig(file)
	if errConf != nil {
		log.Fatal(errConf)
	}

	config.KeyPath = keyPath
	config.Ethereum.Endpoint = requestGCPEndpointSecret(config.ProjectID)

	return &config

}

func loadConfig(file string) (config Config, err error) {
	viper.AddConfigPath("./config/")
	viper.SetConfigName(file)
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func requestGCPEndpointSecret(projectID string) (endpoint string) {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	// Build the request.
	ethAccessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/" + projectID + "/secrets/eth_endpoint_watcher/versions/latest",
	}
	// Call the API.
	res, err := client.AccessSecretVersion(ctx, ethAccessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
		return
	}

	return string(res.Payload.Data)
}

func FindCurrency(addr string, currencies []*CurrencyConfig) *CurrencyConfig {
	for _, c := range currencies {
		if addr == c.Address {
			return c
		}
	}
	// return the first element which is always eth if nothing is found
	return currencies[0]
}
