package dto

// struct to store all configs from file
type Config struct {
	Server  ServerConfig
	Aws     AwsConfig
	MongoDb MongoDbConfig
}

// config for main server
type ServerConfig struct {
	ServerPort string `toml:"serverPort"`
	LogLevel   string `toml:"logLevel"`
}

// config for Amazon S3 server
type AwsConfig struct {
	AwsAccessKeyId     string `toml:"awsAccessKeyId"`
	AwsSecretAccessKey string `toml:"awsSecretAccessKey"`
	AwsRegion          string `toml:"awsRegion"`
	AwsBucket          string `toml:"awsBucket"`
}

// config for NoSql DB MongoDb
type MongoDbConfig struct {
}
