package rpm

import "regexp"

var alphs = regexp.MustCompile("([a-zA-Z]+)|([0-9]+)|(~)")
