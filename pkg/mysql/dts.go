package mysql

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type DTS struct {
	binlogSyncer *replication.BinlogSyncer
	binlogEvent  *BinlogEvent
}

type BinlogEvent struct {
}

func NewDts(c DB) *DTS {
	var (
		host   = strings.Split(c.Addr, ":")
		config = replication.BinlogSyncerConfig{
			ServerID: uint32(os.Geteuid()),
			Host:     host[0],
			Port:     uint16(helpers.StringToInt(host[1])),
			User:     c.Username,
			Password: c.Password,
		}
	)
	return &DTS{
		binlogSyncer: replication.NewBinlogSyncer(config),
		binlogEvent:  new(BinlogEvent),
	}
}

func (e *BinlogEvent) onRow(event *replication.BinlogEvent) (err error) {
	rowsEvent := event.Event.(*replication.RowsEvent)
	table := rowsEvent.Table
	for _, row := range rowsEvent.Rows {
		fmt.Printf("table: %s, scheme: %v\n", table.Table, table.Schema)
		fmt.Printf("row: %v\n", row)
	}
	return nil
}

func (e *BinlogEvent) onDDL(event *replication.BinlogEvent) (err error) {
	event.Dump(os.Stdout)
	return
}

func (d *DTS) Run() {
	streamer, err := d.binlogSyncer.StartSync(mysql.Position{})
	if err != nil {
		panic(err)
	}
	for {
		ctx, cancel := context.WithTimeout(helpers.GetContextWithRequestId(), 2*time.Second)
		event, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			zap.S().Errorf("contex timeout")
			continue
		}
		switch event.Event.(type) {
		case *replication.RowsEvent:
			// 处理行事件
			err := d.binlogEvent.onRow(event)
			if err != nil {
				zap.S().Errorf("d.binlogEvent.OnRow Error. Err: %v", err)
				continue
			}
		case *replication.QueryEvent:
			// 处理DDL事件
			err := d.binlogEvent.onDDL(event)
			if err != nil {
				zap.S().Errorf("d.binlogEvent.OnDDL Error. Err: %v", err)
				continue
			}
		}
	}
}
