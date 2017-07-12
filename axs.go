/*
Axs (pronounced access) implements a simple machine manager. The
configuration file by default is found at $HOME/.axsrc.json. An example of
the configuration file format:

{
	"cc1": {
		"bay-a": {
			"host": "ssh://ben@bay-a:22",
			"serial": "telnet://bay-a:23"
		},
		"bay-b": {
			"host": "ssh://root@192.168.0.1",
			"bmc": "ssh://admin@192.168.0.2"
		}
	}
}
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

func resolveConfigTarget(target string) (*url.URL, error) {
	dest := viper.GetString(target)
	if len(dest) == 0 {
		return nil, fmt.Errorf("target not found")
	}

	return url.Parse(dest)
}

func prepareSSHCommand(target *url.URL) ([]string, error) {
	command := []string{target.Scheme}

	if p := target.Port(); p != "" {
		command = append(command, "-p", p)
	}

	if u := target.User; u != nil && len(u.Username()) > 0 {
		user := u.Username()
		formatted := fmt.Sprintf("%s@%s", user, target.Hostname())
		command = append(command, formatted)
	} else {
		command = append(command, target.Hostname())
	}

	return command, nil
}

func prepareTelnetCommand(target *url.URL) ([]string, error) {
	command := []string{target.Scheme}

	if u := target.User; u != nil && len(u.Username()) > 0 {
		user := u.Username()
		command = append(command, "-l", user)
	}

	command = append(command, target.Hostname())

	if p := target.Port(); len(p) > 0 {
		command = append(command, p)
	}

	return command, nil
}

func prepareTargetCommand(target *url.URL) ([]string, error) {
	switch strings.ToLower(target.Scheme) {
	case "ssh":
		return prepareSSHCommand(target)
	case "telnet":
		return prepareTelnetCommand(target)
	default:
		return nil, fmt.Errorf("unsupported scheme")
	}
}

func initConfig(name string) {
	expanded := os.ExpandEnv(name)
	viper.SetConfigFile(expanded)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] target\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage

	resolve := flag.Bool("resolve", false, "Resolve command.")
	configFile := flag.String("config", "$HOME/.axsrc.json", "Config file.")

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	target := flag.Arg(0)

	initConfig(*configFile)

	url, err := resolveConfigTarget(target)
	if err != nil {
		log.Fatal(err)
	}

	command, err := prepareTargetCommand(url)
	if err != nil {
		log.Fatal(err)
	}

	if *resolve {
		fmt.Println(strings.Join(command, " "))
		return
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
