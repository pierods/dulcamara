package testdata

var ValidRule = `port 2222
method POST
path /path/to/url/
--response
a
response
`

var ValidRuleWithSpaces = ` port 2222   
  method     POST   
path /path/to/url/     
--response    
a
response
`
