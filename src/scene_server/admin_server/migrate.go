/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/pflag"

	"configdatabase/src/common"
	"configdatabase/src/common/blog"
	"configdatabase/src/common/types"
	"configdatabase/src/common/util"
	"configdatabase/src/scene_server/admin_server/app"
	"configdatabase/src/scene_server/admin_server/app/options"
	"configdatabase/src/scene_server/admin_server/command"
)

func main() {
	common.SetIdentification(types.CC_MODULE_MIGRATE)

	runtime.GOMAXPROCS(runtime.NumCPU())

	blog.InitLogs()
	defer blog.CloseLogs()

	op := options.NewServerOption()
	op.AddFlags(pflag.CommandLine)

	if err := command.Parse(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "parse arguments failed, %v\n", err)
		os.Exit(1)
	}
	util.InitFlags()

	ctx, cancel := context.WithCancel(context.Background())
	if err := app.Run(ctx, cancel, op); err != nil {
		fmt.Fprintf(os.Stderr, "run app failed, %v\n", err)
		blog.Errorf("process stopped by %v", err)
		blog.CloseLogs()
		os.Exit(1)
	}
}
