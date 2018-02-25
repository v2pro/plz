package rotation

import "os"

type TriggerBySize struct {
	SizeInKB int64
}

func (trigger *TriggerBySize) UpdateStat(stat interface{}, file *os.File, buf []byte) (interface{}, bool, error) {
	var size int64
	if stat == nil {
		fileInfo, err := file.Stat()
		if err != nil {
			return nil, false, err
		}
		size = fileInfo.Size()
	} else {
		size = stat.(int64)
		size += int64(len(buf))
	}
	triggered := false
	if size >= trigger.SizeInKB*1024 {
		triggered = true
	}
	return size, triggered, nil
}
