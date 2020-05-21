// the app main entrance

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"configdatabase/src/common"
	"configdatabase/src/common/blog"
	"configdatabase/src/common/types"
	"configdatabase/src/common/util"
	"configdatabase/src/scene_server/regal_server/app"
	"configdatabase/src/scene_server/regal_server/app/options"
	"configdatabase/src/scene_server/regal_server/app/test"

	"github.com/spf13/pflag"
)

func main() {
	common.SetIdentification(types.CC_MODULE_REGAL)
	runtime.GOMAXPROCS(runtime.NumCPU())

	op := options.NewServerOption()
	op.AddFlags(pflag.CommandLine)
	util.InitFlags()

	blog.InitLogs()
	defer blog.CloseLogs()
	if err := common.SavePid(); err != nil {
		blog.Error("fail to save pid. err: %s", err.Error())
	}

	test.TestRegal()
	ctx, cancel := context.WithCancel(context.Background())
	if err := app.Run(ctx, cancel, op); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		blog.Errorf("process stopped by %v", err)
		blog.CloseLogs()
		os.Exit(1)
	}
}
