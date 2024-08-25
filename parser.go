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
	parsingBody := false

	for l, line := range lines {
		if parsingBody {
			e.response += line
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] == "--response" {
			parsingBody = true
			continue
		}
		if len(fields) < 2 {
			fmt.Printf("invalid rule %s - invalid line %d:%s\n", r.name, l, line)
			return
		}
		switch fields[0] {
		case "port":
			e.port = fields[1]
		case "method":
			e.method = fields[1]
		case "path":
			e.path = fields[1]
		default:
			fmt.Printf("unknown directive %s:%s - skipping\n", r.name, fields[1])
		}
	}
	err := validateEndpoint(e)
	if err != nil {
		fmt.Printf("invalid rule, skipping. %s:%v\n", r.name, err)
		return
	}
	deploy(e)
}

func validateEndpoint(e endpoint) error {
	if e.port == "" {
		return fmt.Errorf("missing port")
	}
	if e.response == "" {
		return fmt.Errorf("missing response")
	}
	return nil
}
