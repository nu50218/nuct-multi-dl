package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Config struct {
	BaseUri    string  `json:"base_uri"`
	User       string  `json:"user"`
	Pass       string  `json:"pass"`
	LastUpdate *string `json:"last_update"`
	Sites      []*struct {
		ID  string `json:"id"`
		Out string `json:"out"`
	} `json:"sites"`
}

var (
	configFile  = flag.String("f", "nuct-multi-dl.json", "config file")
	failOnError = flag.Bool("failOnError", false, "fail on error")
)

func init() {
	flag.Parse()
}

func main() {
	f, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	config := &Config{}

	if err := json.NewDecoder(f).Decode(config); err != nil {
		log.Fatal(err)
	}

	for _, site := range config.Sites {
		args := []string{"-user=" + config.User, "-pass=" + config.Pass, "-uri=" + config.BaseUri, "-id=" + site.ID, "-out=" + site.Out}
		if config.LastUpdate != nil {
			args = append(args, "-last_update="+*config.LastUpdate)
		}

		fmt.Println("$", "nuct-dl", args[2:])

		cmd := exec.Command("nuct-dl", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Println(err)
			if *failOnError {
				os.Exit(1)
			}
		}
	}
}
