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
	"math/rand"
	"github.com/v2pro/plz/countlog/loglog"
)

// normal => triggered => opened new => normal
const statusNormal = 0
const statusTriggered = 1
const statusArchived = 2

type Writer struct {
	cfg *Config
	// file is owned by the write goroutine
	file *os.File
	// newFile and status shared between write and rotate goroutine
	newFile     unsafe.Pointer
	status      int32
	stat        interface{}
	executor    *concurrent.UnboundedExecutor
	archiveChan chan struct{}
}

type NamingStrategy interface {
	ListFiles() ([]Archive, error)
	NextFile() (string, error)
}

type ArchiveStrategy interface {
	Archive(path string) ([]Archive, error)
}

type Compressor interface {
}

type ArchiveByCompression struct {
	RawArchive ArchiveStrategy
	Retention  RetainStrategy
	Naming     NamingStrategy
	Compressor Compressor
}

type Archive struct {
	Path       string
	ArchivedAt time.Time
	Size       int64
}

type RetainStrategy interface {
	PurgeSet(archives []Archive) []Archive
}

type TriggerStrategy interface {
	UpdateStat(stat interface{}, file *os.File, buf []byte) (interface{}, bool, error)
	TimeToTrigger() time.Duration
}

type PurgeStrategy interface {
	Purge(purgeSet []Archive) error
}

type Config struct {
	WritePath       string
	FileMode        os.FileMode
	DirectoryMode   os.FileMode
	TriggerStrategy TriggerStrategy
	ArchiveStrategy ArchiveStrategy
	RetainStrategy  RetainStrategy
	PurgeStrategy   PurgeStrategy
}

func NewWriter(cfg Config) (*Writer, error) {
	if cfg.FileMode == 0 {
		cfg.FileMode = 0644
	}
	if cfg.DirectoryMode == 0 {
		cfg.DirectoryMode = 0755
	}
	executor := concurrent.NewUnboundedExecutor()
	writer := &Writer{executor: executor, cfg: &cfg}
	err := writer.reopen()
	if err != nil {
		return nil, err
	}
	executor.Go(writer.rotateInBackground)
	return writer, nil
}

func (writer *Writer) reopen() error {
	cfg := writer.cfg
	file, err := os.OpenFile(cfg.WritePath, os.O_WRONLY|os.O_APPEND, cfg.FileMode)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		os.MkdirAll(path.Dir(cfg.WritePath), cfg.DirectoryMode)
		file, err = os.OpenFile(cfg.WritePath,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, cfg.FileMode)
		if err != nil {
			return err
		}
	}
	writer.file = file
	return nil
}

func (writer *Writer) Write(buf []byte) (int, error) {
	if atomic.LoadInt32(&writer.status) == statusArchived {
		err := writer.file.Close()
		if err != nil {
			loglog.Error(err)
		}
		writer.reopen()
	}
	file := writer.file
	triggerStrategy := writer.cfg.TriggerStrategy
	n, err := file.Write(buf)
	if atomic.LoadInt32(&writer.status) == statusNormal {
		var triggered bool
		var err error
		writer.stat, triggered, err = triggerStrategy.UpdateStat(writer.stat, file, buf[:n])
		if err != nil {
			loglog.Error(err)
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
	triggerStrategy := writer.cfg.TriggerStrategy
	archiveStrategy := writer.cfg.ArchiveStrategy
	retainStrategy := writer.cfg.RetainStrategy
	purgeStrategy := writer.cfg.PurgeStrategy
	var timer <-chan time.Time
	for {
		duration := triggerStrategy.TimeToTrigger()
		if duration > 0 {
			duration += time.Duration(rand.Int63n(int64(duration)))
			timer = time.NewTimer(duration).C
		}
		select {
		case <-ctx.Done():
			return
		case <-writer.archiveChan:
		case <-timer:
		}
		archives, err := archiveStrategy.Archive(writer.cfg.WritePath)
		if err != nil {
			loglog.Error(err)
			// retry after one minute
			timer = time.NewTimer(time.Minute).C
			continue
		}
		atomic.StoreInt32(&writer.status, statusArchived)
		purgeSet := retainStrategy.PurgeSet(archives)
		err = purgeStrategy.Purge(purgeSet)
		if err != nil {
			loglog.Error(err)
		}
	}
}
