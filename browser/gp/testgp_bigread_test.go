//go:build e2e
// +build e2e

package gp

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaybinks/immich-go/helpers/fshelper"
	"github.com/jaybinks/immich-go/logger"
)

func TestReadBigTakeout(t *testing.T) {
	f, err := os.Create("bigread.log")
	if err != nil {
		panic(err)
	}

	j := logger.NewJournal(logger.NewLogger(logger.Info, true, false).SetWriter(f))
	m, err := filepath.Glob("../../../test-data/full_takeout/*.zip")
	if err != nil {
		t.Error(err)
		return
	}
	cnt := 0
	fsyss, err := fshelper.ParsePath(m, true)
	to, err := NewTakeout(context.Background(), j, fsyss...)
	if err != nil {
		t.Error(err)
		return
	}

	for range to.Browse(context.Background()) {
		cnt++
	}
	to.jnl.Report()
	t.Logf("seen %d files", cnt)
}
