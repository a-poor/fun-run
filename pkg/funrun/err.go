package funrun

import "fmt"

type MapError map[string]error

func (m MapError) Error() string {
	var s string
	for k, v := range m {
		s += fmt.Sprintf("(%s) %s\n", k, v.Error())
	}
	return s
}
