package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vote",
	Short: "Vote is a CLI app for upcoming elections",
}

type config struct {
	Address string `yaml:"address"`
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var c config
	if err := c.getConfig(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("config contents: %v", c)

	watch()
}

func (c *config) getConfig() error {
	yml, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("read config file err: %v", err)
	}
	if err = yaml.Unmarshal(yml, c); err != nil {
		return fmt.Errorf("unmarshal config error: %v", err)
	}
	return nil
}

func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					fmt.Println("event:", event)
				case err, ok := <- watcher.Errors:
					if !ok {
						return
					}
					fmt.Printf("error watching for errors: %v", err)
			}
		}
	}()
	if err := watcher.Add("config.yaml"); err != nil {
		fmt.Printf("err watching file: %v", err)
	}
	<- done
}
