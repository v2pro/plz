package countlog

import "strings"

type CompactFormat struct {
}

func (format *CompactFormat) FormatLog(event Event) []byte {
	line := []byte{}
	line = append(line, event.Event...)
	for i := 0; i < len(event.Properties); i += 2 {
		k, _ := event.Properties[i].(string)
		if k == "" {
			continue
		}
		v := formatV(event.Properties[i+1])
		if v == "" {
			continue
		}
		v = strings.Replace(v, "\r", `\r`, -1)
		v = strings.Replace(v, "\n", `\n`, -1)
		v = strings.Replace(v, "||", " ", -1)
		line = append(line, "||"...)
		line = append(line, k...)
		line = append(line, '=')
		line = append(line, v...)
	}
	line = append(line, '\n')
	return line
}
