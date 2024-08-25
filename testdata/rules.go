package testdata

import "github.com/pierods/dulcamara"

var ValidRule = dulcamara.Rule{
	name: "validRule",
	response: `port 2222
method POST
path /path/to/url/
--response
a
response
`,
}
