package sync

import (
	"fmt"
	"strings"
	"testing"
)

type sf struct {
	Group

}

func TestSingleflight(t *testing.T) {

	fmt.Println(kebabString("sss-ssss-bhbhSddSdddHddYss"))
}

func kebabString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '-')
		}
		if d != '-' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
