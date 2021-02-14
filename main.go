package main

import (
	"database/sql"
	"flag"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Conn       string            `yaml:"conn"`
	Port       int               `yaml:"port"`
	Publishers map[string]string `yaml:"publishers"`
}

var configPath string
var config *Config

func init() {
	flag.StringVar(&configPath, "c", "config.yaml", "path to the config file")
	flag.Parse()
}

func main() {
	var err error
	config, err = NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	for uuid, view := range config.Publishers {
		http.HandleFunc("/"+uuid, func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet {
				log.Info("report requested: ", uuid)
				report, err2 := getReport(view)
				if err2 != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					log.Error(err2)
					return
				}
				writer.Header().Add("Content-Type", "application/json")
				writer.Write([]byte(report))
			} else {
				writer.WriteHeader(http.StatusMethodNotAllowed)
			}
		})
	}
	http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
}

func getReport(name string) (string, error) {
	sqlStatement := "select * from " + name
	db, err := sql.Open("postgres", config.Conn)
	defer db.Close()
	if err != nil {
		log.Fatal("Could connect to database", err)
	}
	row := db.QueryRow(sqlStatement)
	var json *string
	err = row.Scan(&json)
	if err != nil {
		log.Error(err)
		return "[]", nil
	}

	if json == nil {
		return "[]", nil
	}

	return *json, nil
}

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
