package secretariat

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"sync"
	"time"
)

type Something interface{}

type Workflow struct {
	ctx         context.Context
	wid         string
	lock        sync.Mutex
	withCompass bool
	conveyor    chan interface{}
	bufferPool  []byte
	bufferSize  int
	scannerRate time.Duration
	folder      *Folder
	startAt     time.Time
	endAt       time.Time
	stopChan    chan struct{}
}

func NewWorkflow(ctx context.Context, bufferSize int, scannerRate time.Duration, withCompass bool) *Workflow {
	return &Workflow{
		ctx:         ctx,
		wid:         helpers.GenNanoId(10),
		conveyor:    make(chan interface{}, bufferSize),
		bufferPool:  make([]byte, 0, bufferSize),
		bufferSize:  bufferSize,
		scannerRate: scannerRate,
		folder:      NewFolder(ctx, bufferSize*10, helpers.GenNanoId(10)),
		lock:        sync.Mutex{},
		stopChan:    make(chan struct{}),
		withCompass: withCompass,
	}
}

func (workflow *Workflow) Put(sth Something) {
	workflow.conveyor <- sth
}

func (workflow *Workflow) running() {
	workflow.startAt = time.Now()
	go workflow.poolWatcher()
	go workflow.poolCleaner()
	go workflow.folder.typewriterWorking()
}

func (workflow *Workflow) stop() {
	workflow.endAt = time.Now()
	workflow.stopChan <- struct{}{}
	fmt.Printf("[%s] Workflow[%s] Done...[duration: %s]\n", time.Now().Format("2006-01-02T15:04:05.999999"), workflow.wid, time.Since(workflow.startAt).String())
}

func (workflow *Workflow) poolCleaner() {
	var ticker = time.NewTicker(workflow.scannerRate)
	for {
		select {
		case <-ticker.C:
			workflow.lock.Lock()
			workflow.flush()
			workflow.lock.Unlock()
		case <-workflow.stopChan:
			workflow.flush()
			fmt.Printf("[%s] Workflow.poolCleaner Done...[duration: %s]\n", time.Now().Format("2006-01-02T15:04:05.999999"), time.Since(workflow.startAt).String())
		}
	}
}

func (workflow *Workflow) poolWatcher() {
	for {
		select {
		case sth := <-workflow.conveyor:
			stream, err := sonic.Marshal(sth)
			if err != nil {
				logger.Helper().Errorf("sonic.Marshal error: [%v], sth: [%#v]", err, sth)
				continue
			}
			workflow.pooling(stream)
		case <-workflow.stopChan:
			workflow.flush()
			fmt.Printf("[%s] Workflow.poolWatcher Done...[duration: %s]\n", time.Now().Format("2006-01-02T15:04:05.999999"), time.Since(workflow.startAt).String())
		}
	}
}

func (workflow *Workflow) pooling(stream []byte) {
	defer workflow.lock.Unlock()

	if workflow.withCompass {
		var (
			err error
		)
		if stream, err = compressData(stream); err != nil {
			panic("compress data err...")
		}
	}

	workflow.lock.Lock()
	if (len(workflow.bufferPool) + len(stream)) > workflow.bufferSize {
		workflow.flush()
		workflow.bufferPool = append(workflow.bufferPool, stream...)
		return
	}
	workflow.bufferPool = append(workflow.bufferPool, stream...)
	return
}

func (workflow *Workflow) flush() {
	workflow.folder.fileStream <- workflow.bufferPool
	workflow.bufferPool = make([]byte, 0, workflow.bufferSize)
	return
}
