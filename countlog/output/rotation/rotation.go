package rotation

import (
	"os"
	"time"
)

type Writer struct {
	f *os.File
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
	Trigger   TriggerStrategy
	Archive   ArchiveStrategy
	Retention RetentionStrategy
	Purge     PurgeStrategy
}

func NewWriter(cfg Config) *Writer {
	return &Writer{}
}

func (writer *Writer) Write(buf []byte) (int, error) {
	return writer.f.Write(buf)
}
