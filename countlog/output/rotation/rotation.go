package rotation

import (
	"os"
	"time"
)

type Writer struct {
	file *os.File
}

type NamingStrategy interface {
}

type NameByTime struct {
	Directory string
	Pattern   string
}

type NameByIndex struct {
	Directory string
	Pattern   string
}

type ArchiveStrategy interface {
}

type ArchiveByMove struct {
	Naming NamingStrategy
}

type Compressor interface {
}

type ArchiveByCompression struct {
	RawArchive ArchiveStrategy
	Retention  RetentionStrategy
	Naming     NamingStrategy
	Compressor Compressor
}

type Archive struct {
	Path       string
	ArchivedAt time.Time
	Size       int64
}

type RetentionStrategy interface {
}

type RetainByCount struct {
	MaxArchivesCount int
}

type TriggerStrategy interface {
}

type TriggerByInterval struct {
	Hourly   bool
	Daily    bool
	Weekly   bool
	Monthly  bool
	Location *time.Location
}

type PurgeStrategy interface {
}

type PurgeByDeletion struct {
}

type Config struct {
	WritePath string
	FileMode  os.FileMode
	Trigger   TriggerStrategy
	Archive   ArchiveStrategy
	Retention RetentionStrategy
	Purge     PurgeStrategy
}

func NewWriter(cfg Config) (*Writer, error) {
	fileMode := cfg.FileMode
	if fileMode == 0 {
		fileMode = 0644
	}
	file, err := os.OpenFile(cfg.WritePath, os.O_WRONLY|os.O_APPEND, fileMode)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	file, err = os.OpenFile(cfg.WritePath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileMode)
	if err != nil {
		return nil, err
	}
	return &Writer{file: file}, nil
}

func (writer *Writer) Write(buf []byte) (int, error) {
	return writer.file.Write(buf)
}

func (writer *Writer) Close() error {
	return writer.file.Close()
}
