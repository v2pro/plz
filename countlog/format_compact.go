package countlog

import (
	"strings"
	"time"
)

type CompactFormat struct {
	StringLengthCap      int
}

func (format *CompactFormat) FormatLog(event Event) []byte {
	line := []byte{}
	line = append(line, event.Event...)
	for i := 0; i < len(event.Properties); i += 2 {
		k, _ := event.Properties[i].(string)
		if k == "" {
			continue
		}
		v := ""
		if k == "timestamp" {
			tm := event.Properties[i+1].(int64)
			v = time.Unix(tm/1e9, tm%1e9).Format("2006-01-02 15:04:05.999999999")
		} else {
			v = formatV(event.Properties[i+1])
		}

		if v == "" {
			continue
		}
		if event.Level < LevelWarn && format.StringLengthCap > 0 {
			lenCap := len(v)
			if format.StringLengthCap < lenCap {
				lenCap = format.StringLengthCap
				v = v[:lenCap] + "...more, capped"
			}
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
