package secretariat

import (
	"context"
	"fmt"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Folder struct {
	ctx     context.Context
	lock    sync.Mutex
	key     string
	workDir string

	fileName   string
	filePath   string
	fileStream chan []byte

	startAt  time.Time
	stopChan chan struct{}
}

func NewFolder(ctx context.Context, bufferSize int, key string) *Folder {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath, err := filepath.Abs(filepath.Join(workdir, fileDir, fmt.Sprintf("%d", os.Getpid()), key))
	if err != nil {
		panic(err)
	}
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		_ = os.MkdirAll(filePath, 0744)
	}
	fmt.Printf("\nfilePath: %s\n", filePath)

	return &Folder{
		ctx:        ctx,
		lock:       sync.Mutex{},
		key:        key,
		workDir:    workdir,
		filePath:   filePath,
		fileStream: make(chan []byte, bufferSize),
		stopChan:   make(chan struct{}, 1),
	}
}

func (folder *Folder) typewriterWorking() {
	folder.startAt = time.Now()
	for {
		select {
		case stream := <-folder.fileStream:
			folder.write(stream)
		case <-folder.stopChan:
			fmt.Printf("[%s] typewriter[%s] Stopped...[duration: %s]\n", time.Now().Format("2006-01-02T15:04:05.999999"), folder.key, time.Since(folder.startAt).String())
		}
	}
}

func (folder *Folder) write(stream []byte) {
	defer folder.lock.Unlock()
	var filePath = folder.getCurrentFilePath(time.Now())

	folder.lock.Lock()
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
	if err != nil {
		logger.Helper().Errorf("[Write] os.OpenFile Error. err[%v], filePath[%s]", err, filePath)
		panic("can not open file")
	}
	defer file.Close()

	if _, err := file.Write(stream); err != nil {
		logger.Helper().Errorf("[Write] file.Write Error. err[%v], filePath[%s]", err, filePath)
		panic("can not write file")
	}
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}
}

func (folder *Folder) read([]byte) {

}

func (folder *Folder) getCurrentFilePath(now time.Time) string {
	var (
		_filepath string
		_now      = now.Format("150405")
	)
	for i := 1; i <= 9999; i++ {
		_filepath = filepath.Join(folder.filePath, fmt.Sprintf("%s_%04d.zfs", _now, i))
		if _, err := os.Stat(_filepath); os.IsNotExist(err) {
			return _filepath
		}
		continue
	}
	panic("file index exhausted...")
}
