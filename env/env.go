package env

import (
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var Vars Config

func init() {
	readFile(&Vars)
	readEnv(&Vars)
	fmt.Printf("%+v", Vars)
}

type Config struct {
	Api struct {
		Env  string `yaml:"env", envconfig:"API_ENV"`
		Port string `yaml:"port", envconfig:"API_PORT"`
	} `yaml:"api"`

	Importer struct {
		Env     string `yaml:"env", envconfig:"IMPORTER_ENV"`
		PathCsv string `yaml:"path-csv", envconfig:"IMPORTER_PATH_CSV"`
	} `yaml:"importer"`

	Database struct {
		Name string `yaml:"name", envconfig:"DB_NAME"`
		Host string `yaml:"host", envconfig:"DB_HOST"`
	} `yaml:"database"`

	Swapi struct {
		Url string `yaml:"url", envconfig:"SWAPI_URL"`
	} `yaml:"swapi"`
}

func processError(err error) {
	log.Print(err)
}

func readFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}

	err = f.Close()

	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
