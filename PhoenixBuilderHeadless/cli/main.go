package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"phoenixbuilder/session"
	bridge_fmt "phoenixbuilder/session/bridge/fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var consoleInput chan string

type CLI struct {
}

func (cli *CLI) Run(configFile string) {
	fp, err := os.OpenFile(configFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	byteConfig, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	fp.Close()
	config := session.NewConfig()
	err = yaml.Unmarshal(byteConfig, config)
	if err != nil {
		panic(err)
	}

	bridge_fmt.HookFunc = func(s string) {
		fmt.Printf(">: %s", s)
	}
	s := session.NewSession(config)
	terminateChan, err := s.Start()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Session Start Successfully\n")
	updatedByteConfig, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	fp, err = os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	fp.Write(updatedByteConfig)
	fp.Close()

	go func() {
		for {
			select {
			case <-terminateChan:
				return
			case cmd := <-consoleInput:
				if len(cmd) > 0 && cmd[0] == '!' {
					switch cmd {
					case "!stop":
						s.Stop()
						break
					default:
						fmt.Println("unknown system command: ", cmd)
					}
					continue
				}
				s.Execute(cmd)
			}
		}
	}()

	reason := <-terminateChan
	fmt.Printf("Cli Terminated (%v)", reason)
}

func main() {
	consoleInput = make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			cmd, _ := reader.ReadString('\n')
			cmd = strings.TrimSpace(cmd)
			consoleInput <- cmd
		}
	}()

	cli := &CLI{}
	cli.Run("config.yaml")
}
