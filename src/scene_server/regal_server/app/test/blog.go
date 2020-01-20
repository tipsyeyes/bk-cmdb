package test

import (
	"configdatabase/src/common/blog"
	"fmt"
)

// blog usage:
// fmt.Print是直接打到标准输出的
// blog的日志分四个等级:info、warning、error、fatal
// blog.Debug与 Info、Infof相同
// 日志根据等级分类一般会存储在 ./logs目录下相应的文件中，info文件中会冗余存储所有等级日志信息，同理，fatal文件中只保存最高等级日志
// 日志默认格式：I0120 15:08:02.710501   27200 test/blog.go:18] info...
// I[W/E/F]月日 时:分:秒.纳秒 进程id 文件:行数] content
// --logtostderr=true的时候，所有等级的 blog日志都会默认打到标准错误中
// 还可以自定义 level日志，通过 --v=3定义当前系统的日志等级，然后 blog.V(5).Info的方式来决定是否记录
func logTest()  {

	blog.Debug("debug...")
	blog.Info("info...")
	blog.Infof("infof...%d", 100)
	blog.Warn("warn...")
	blog.Warnf("warnf...%s", "100")
	blog.Error("error...")
	blog.Errorf("errorf...%+v", fmt.Errorf("100"))
	// blog.Fatal("fatal...")

	// 只有日志等级小于等于 --v参数指定的等级时，日志才会被记录
	blog.V(3).Info("the 3 level log...")
	blog.V(5).Info("the 5 level log...")
	if blog.V(9) {
		blog.Info("the 9 level log...")
	}

}