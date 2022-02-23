package app

type Config struct {
	Address         string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/shortener_storage.csv"`
}
