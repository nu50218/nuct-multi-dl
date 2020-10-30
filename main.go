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
	BaseUri string `json:"base_uri"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	Sites   []*struct {
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
		fmt.Println("$", "nuct-dl", "-id="+site.ID, "-out="+site.Out)
		cmd := exec.Command("nuct-dl", "-uri="+config.BaseUri, "-id="+site.ID, "-user="+config.User, "-pass="+config.Pass, "-out="+site.Out)
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
