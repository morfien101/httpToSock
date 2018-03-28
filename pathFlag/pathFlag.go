package pathFlag

import (
	"fmt"
	"strings"
)

// PF This allows us to make a list but repeating the --path flag
type PF []string

// PathFlag is an interface that can be used with the flags package.
// It can also be used to split the flags into something we can use.
type PathFlag interface {
	String() string
	Set(string) error
	Split() (map[string]string, error)
}

// String returns a string representation of the Path Flag
func (p *PF) String() string {
	return strings.Join(*p, ",")
}

// Set consumes the value passed in and adds to the slice.
func (p *PF) Set(value string) error {
	*p = append(*p, value)
	return nil
}

// Split will split the flag into a usable map.
func (p *PF) Split() (map[string]string, error) {
	// --path /_status:GET /_status
	// --path /_info:GET /_info
	output := make(map[string]string)
	for _, v := range *p {
		splitter := strings.Split(v, ":")
		if len(splitter) > 2 {
			return map[string]string{}, fmt.Errorf("path mapping %s has broken down into more than 2 parts: %s", splitter[0], strings.Join(splitter[1:], " - "))
		}
		output[splitter[0]] = splitter[1]
	}
	return output, nil
}
