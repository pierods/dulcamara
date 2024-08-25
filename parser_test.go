package dulcamara

import (
	"testing"

	"github.com/pierods/dulcamara/testdata"
)

func Test_parseRule(t *testing.T) {

	endpoint, err := ParseRule(Rule{Response: testdata.ValidRule})
	if err != nil {
		t.Fatal("Should correctly parse a valid rule")
	}
	if endpoint.port != "2222" || endpoint.method != "POST" || endpoint.path != "/path/to/url/" || endpoint.response != "aresponse" {
		t.Fatal("Should correctly parse a valid mock")
	}

	endpoint, err = ParseRule(Rule{Response: testdata.ValidRuleWithSpaces})
	if err != nil {
		t.Fatal("Should correctly parse a valid rule")
	}
	if endpoint.port != "2222" || endpoint.method != "POST" || endpoint.path != "/path/to/url/" || endpoint.response != "aresponse" {
		t.Fatal("Should correctly parse a valid mock")
	}
}
