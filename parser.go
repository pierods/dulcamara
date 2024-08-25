package dulcamara

import (
	"fmt"
	"strings"
)

type Rule struct {
	Name     string
	Response string
}

type endpoint struct {
	rule     string
	port     string
	method   string
	path     string
	response string
}

func ParseRule(r Rule) (endpoint, error) {

	e := endpoint{
		rule: r.Name,
	}

	lines := strings.Split(r.Response, "\n")
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
			return e, fmt.Errorf("invalid rule, skipping %s - invalid line %d:%s", r.Name, l, line)
		}
		switch fields[0] {
		case "port":
			e.port = fields[1]
		case "method":
			e.method = fields[1]
		case "path":
			e.path = fields[1]
		default:
			fmt.Printf("unknown directive %s:%s - skipping\n", r.Name, fields[1])
		}
	}
	err := validateEndpoint(e)
	if err != nil {
		return e, fmt.Errorf("invalid rule, skipping. %s:%v", r.Name, err)

	}
	return e, nil
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
