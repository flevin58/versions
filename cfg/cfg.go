package cfg

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//go:embed versions.yaml
var yamlFile []byte

var Data YamlData
var ConfigFile string

type YamlData struct {
	Editor   string
	Commands []Command
}

type Command struct {
	Name        string `yaml:"name"`
	VersionFlag string `yaml:"flag"`
	VersionLine int    `yaml:"line"`
}

func init() {
	var (
		cfgdir string
		err    error
	)

	log.SetFlags(0)

	// Determine path for the "versions.yaml" file
	// Create it if undefined
	if cfgdir, err = os.UserConfigDir(); err != nil {
		log.Fatalln("could not locate the user config folder")
	}
	cfgdir = filepath.Join(cfgdir, "versions")
	os.MkdirAll(cfgdir, 0755)

	// Read the config file if it exists, otherwise copy the embedded yaml
	ConfigFile = filepath.Join(cfgdir, "versions.yaml")
	ConfigData, err := os.ReadFile(ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			os.WriteFile(ConfigFile, yamlFile, 0664)
			ConfigData = yamlFile
		} else {
			log.Fatalf("could not read the config file:%T", err)
		}
	}
	err = yaml.Unmarshal(ConfigData, &Data)
	if err != nil {
		log.Fatalf("could not load the config file: %v", err)
	}
}
