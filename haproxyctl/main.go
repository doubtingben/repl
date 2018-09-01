package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	haproxy "github.com/bcicen/go-haproxy"
)

var prompt = "\ncli > "
var haproxyAddr = "unix://dev/haproxy/run/haproxy.sock"
var helpCommand = regexp.MustCompile(`^help$`)
var noopCommand = regexp.MustCompile(`^(?:noop.*)?$`)
var showstatsCommand = regexp.MustCompile(`^show\s?stats$`)
var showinfoCommand = regexp.MustCompile(`^show\s?info$`)
var showCommand = regexp.MustCompile(`^show\s(\w+)\s?(\w+)?\s?(\w+)?$`)

func main() {
	input := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(input)
	scanner.Split(cliSplitFunc)
	printHeader()
	for scanner.Scan() {
		command := strings.TrimSuffix(scanner.Text(), "\n")
		switch {
		case helpCommand.MatchString(command):
			fmt.Println(showHelp())
		case noopCommand.MatchString(command):
			fmt.Printf("noop")
		case showstatsCommand.MatchString(command):
			stats, err := showstats()
			if err != nil {
				fmt.Println("An error occurred: " + err.Error())
			}
			for _, i := range stats {
				fmt.Printf("%s %s: %s\n", i.PxName, i.SvName, i.Status)
			}
		case showinfoCommand.MatchString(command):
			results, err := showinfo()
			if err != nil {
				fmt.Println("An error occurred: " + err.Error())
			}
			fmt.Printf("%s", results)
		case showCommand.MatchString(command):
			stats, err := showstats()
			if err != nil {
				fmt.Println("An error occurred: " + err.Error())
			}
			fmt.Println(showParse(command, stats))
		default:
			fmt.Printf("wtf, %s", command)
		}
		fmt.Print(prompt)
	}
}

func showParse(command string, stats []*haproxy.Stat) string {
	result := ""
	for _, i := range stats {

		fmt.Printf("%s %s: %s\n", i.PxName, i.SvName, i.Status)
	}
	return result
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

func showHelp() string {
	return "Good luck!"
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
