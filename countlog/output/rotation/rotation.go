package rotation

import (
	"os"
	"time"
	"path"
	"github.com/v2pro/plz/concurrent"
	"context"
	"unsafe"
	"sync/atomic"
	"github.com/v2pro/plz/countlog/spi"
)

// normal => triggered => opened new => normal
const statusNormal = 0
const statusTriggered = 1
const statusOpenedNew = 2

type Writer struct {
	cfg         *Config
	// file is owned by the write goroutine
	file        *os.File
	// newFile and status shared between write and rotate goroutine
	newFile     unsafe.Pointer
	status      int32
	stat        interface{}
	executor    *concurrent.UnboundedExecutor
	archiveChan chan struct{}
}

type NamingStrategy interface {
	ListFiles() ([]string, error)
	NextFile() (string, error)
}

type NameByTime struct {
	Directory string
	Pattern   string
}

type ArchiveStrategy interface {
	Archive(oldFile *os.File)
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
	UpdateStat(stat interface{}, file *os.File, buf []byte) (interface{}, bool, error)
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
	WritePath     string
	FileMode      os.FileMode
	DirectoryMode os.FileMode
	Trigger       TriggerStrategy
	Archive       ArchiveStrategy
	Retention     RetentionStrategy
	Purge         PurgeStrategy
}

func NewWriter(cfg Config) (*Writer, error) {
	fileMode := cfg.FileMode
	if fileMode == 0 {
		fileMode = 0644
	}
	dirMode := cfg.DirectoryMode
	if dirMode == 0 {
		dirMode = 0755
	}
	file, err := os.OpenFile(cfg.WritePath, os.O_WRONLY|os.O_APPEND, fileMode)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		os.MkdirAll(path.Dir(cfg.WritePath), dirMode)
		file, err = os.OpenFile(cfg.WritePath,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileMode)
		if err != nil {
			return nil, err
		}
	}
	executor := concurrent.NewUnboundedExecutor()
	writer := &Writer{executor: executor, file:file}
	executor.Go(writer.rotateInBackground)
	return writer, nil
}

func (writer *Writer) Write(buf []byte) (int, error) {
	file := writer.file
	n, err := file.Write(buf)
	if atomic.LoadInt32(&writer.status) == statusNormal {
		var triggered bool
		var err error
		trigger := writer.cfg.Trigger
		writer.stat, triggered, err = trigger.UpdateStat(writer.stat, file, buf[:n])
		if err != nil {
			spi.OnError(err)
			return n, err
		}
		if triggered {
			atomic.StoreInt32(&writer.status, statusTriggered)
		}
	}
	return n, err
}

func (writer *Writer) Close() error {
	writer.executor.StopAndWaitForever()
	return writer.file.Close()
}

func (writer *Writer) rotateInBackground(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-writer.archiveChan:
			return
		}
	}
}
