package nuomi

import (
	"context"
	"fmt"
	"go-nuomi/dirbuster"
	"go-nuomi/lib"
	"go-nuomi/lib/http"
	"log"
	"os"
	"sync"
)

type NuoMiOption struct {
	LibOption  lib.Options
	DirOption  dirbuster.DirOptions
	HttpOption http.HTTPOptions
}

type NuoMi struct {
	Opts             *lib.Options
	context          context.Context
	requestsExpected int
	requestsIssued   int
	mu               *sync.RWMutex
	plugin           NuoMiPlugin
	resultChan       chan Result
	errorChan        chan error
	LogInfo          *log.Logger
	LogError         *log.Logger
}

func NewNuoMi(c context.Context, opts *lib.Options, plugin NuoMiPlugin) (*NuoMi, error) {
	var n NuoMi
	n.Opts = opts
	n.plugin = plugin
	n.mu = new(sync.RWMutex)
	n.context = c
	n.resultChan = make(chan Result)
	n.errorChan = make(chan error)
	//可能用不着
	n.LogInfo = log.New(os.Stdout, "", log.LstdFlags)
	n.LogInfo = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
	return &n, nil
}

func NuoMiRunner(prevCtx context.Context, opts *lib.Options, plugin NuoMiPlugin) error {
	if opts == nil {
		return fmt.Errorf("please provide valid options")
	}

	if plugin == nil {
		return fmt.Errorf("please provide a valid plugin")
	}

	ctx, cancel := context.WithCancel(prevCtx)
	defer cancel()

	nuomi, err := NewNuoMi(ctx, opts, plugin)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go errorWorker(nuomi, &wg)
	// 监听操作入库操作
	go resultWorker(nuomi, opts.OutputFilename, &wg)

	if !opts.Quiet && !opts.NoProgress {
		// if not quiet add a new workgroup entry and start the goroutine
		wg.Add(1)
		go progressWorker(ctx, nuomi, &wg)
	}

	err = nuomi.roll()
	//fixme:为什么要cancel?
	cancel()
	wg.Wait()
	if err != nil{
		return err
	}

	return nil
}
