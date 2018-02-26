package rotation

import (
	"os"
	"time"
	"path"
	"sort"
	"github.com/v2pro/plz/clock"
)

type NameByTime struct {
	Directory string
	Pattern   string
}

func (namer *NameByTime) NextFile() (string, error) {
	now := clock.Now()
	return path.Join(namer.Directory, now.Format(namer.Pattern)), nil
}

func (namer *NameByTime) ListFiles() ([]Archive, error) {
	dir, err := os.Open(namer.Directory)
	if err != nil {
		return nil, err
	}
	names, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	var timedFiles timedFiles
	for _, name := range names {
		fileTime, err := time.Parse(namer.Pattern, name)
		if err != nil {
			continue
		}
		timedFiles = append(timedFiles, timedFile{
			file: path.Join(namer.Directory, name),
			time: fileTime,
		})
	}
	sort.Sort(timedFiles)
	archives := make([]Archive, len(timedFiles))
	for i, timedFile := range timedFiles {
		stat, err := os.Stat(timedFile.file)
		if err != nil {
			return nil, err
		}
		archives[i] = Archive{
			Path: timedFile.file,
			ArchivedAt: timedFile.time,
			Size: stat.Size(),
		}
	}
	return archives, nil
}

type timedFile struct {
	file string
	time time.Time
}

type timedFiles []timedFile

func (a timedFiles) Len() int           { return len(a) }
func (a timedFiles) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a timedFiles) Less(i, j int) bool { return a[i].time.Before(a[j].time) }