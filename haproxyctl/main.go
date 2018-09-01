package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	haproxy "github.com/bcicen/go-haproxy"
)

var haproxyAddr = "unix://dev/haproxy/run/haproxy.sock"
var showCommand = regexp.MustCompile(`^show\s(\w+)\s(\w+)$`)

func main() {
	input := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(input)
	scanner.Split(cliSplitFunc)
	printHeader()
	for scanner.Scan() {
		switch text := strings.TrimSuffix(scanner.Text(), "\n"); text {
		case "help":
			fmt.Printf("here's your help: %s", text)
		case "noop":
			fmt.Printf("Here's noop: %s", text)
		case "showinfo":
			text, err := showinfo()
			if err != nil {
				fmt.Println("An error occurred: " + err.Error())
			}
			fmt.Printf("%s", text)
		case "showstats":
			stats, err := showstats()
			if err != nil {
				fmt.Println("An error occurred: " + err.Error())
			}
			for _, i := range stats {
				fmt.Printf("%s %s: %s\n", i.PxName, i.SvName, i.Status)
			}
		case "show":
			stats, err := show(text)
			fmt.Printf("showing: %s\n", text)
			if err != nil {
				fmt.Println("An error occurred: " + err.Error())
			}
			for _, i := range stats {
				fmt.Printf("%s %s: %s\n", i.PxName, i.SvName, i.Status)
			}
		case "": // Print nothing when no input recieved
		default:
			fmt.Printf("wtf, %s", text)
		}
		fmt.Print("\ncli > ")
	}
}

func show(command string) ([]*haproxy.Stat, error) {
	client := &haproxy.HAProxyClient{
		Addr: haproxyAddr,
	}

	result, err := client.Stats()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func showinfo() (string, error) {
	client := &haproxy.HAProxyClient{
		Addr: haproxyAddr,
	}

	result, err := client.RunCommand("show info")
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func showstats() ([]*haproxy.Stat, error) {
	client := &haproxy.HAProxyClient{
		Addr: haproxyAddr,
	}

	result, err := client.Stats()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func printHeader() {
	fmt.Println("Welcome to the cli")
	fmt.Print("cli > ")
}

func cliSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if string(data[len(data)-1]) == "\n" {
		return len(data), data, nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return
}
