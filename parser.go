package main

import (
	"fmt"
	"strings"
)

type endpoint struct {
	rule     string
	port     string
	method   string
	path     string
	response string
}

func parseRule(r rule) {

	e := endpoint{
		rule: r.name,
	}

	lines := strings.Split(r.response, "\n")
	for l, line := range lines {
		if len(line) == 0 {
			continue
		}
		split := strings.SplitN(line, " ", 2)
		if len(split) < 2 || split[1] == "" {
			fmt.Printf("invalid rule %s: invalid line %d:%s\n", r.name, l, line)
			return
		}
		switch split[0] {
		case "port":
			e.port = strings.TrimSpace(split[1])
		case "method":
			e.method = strings.TrimSpace(split[1])
		case "path":
			e.path = strings.TrimSpace(split[1])
		case "response":
			e.response = split[1]
		}
	}
	err := validateEndpoint(e)
	if err != nil {
		fmt.Println(err)
		return
	}
	deploy(e)
}

func validateEndpoint(e endpoint) error {
	if e.port == "" {
		return fmt.Errorf("missing port")
	}
	return nil
}
