package test

import (
	"fmt"
	"github.com/spf13/pflag"
	"strings"
	"testing"
)

func init()  {

}

// WordSepNormalizeFunc changes all flags that contain "_" separators
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

// usage of pflag:
func TestPflag(t *testing.T) {
	// 和 flag使用方式类似
	// 默认初始化
	var name string = "default"
	var username string = "default"
	var age int = -1
	// pflag默认初始化
	pflag.StringVar(&name, "name", "vincent", "help message for name")
	pflag.StringVar(&username, "username", "vincent", "help message for username")
	pflag.IntVar(&age, "age", 0, "help message for age")

	// 设置规范参数名称函数 WordSepNormalizeFunc
	pflag.CommandLine.SetNormalizeFunc(WordSepNormalizeFunc)
	//pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	// 为 age 参数设置 NoOptDefVal
	// 如果命令行只传了 --age，而没有具体 value时的默认值
	pflag.Lookup("age").NoOptDefVal = "25"
	// 将 --name标记为即将要废弃的，请用户使用 --username，标记废弃的参数不会出现在用户提示中
	pflag.CommandLine.MarkDeprecated("name", "please use --username instead")

	// 从命令行解析
	//pflag.Parse()
	//pflag.CommandLine.Parse([]string{
	//	"--name", "tes", "--age", "100",
	//})
	_ = pflag.CommandLine.Parse([]string{
		"--name", "tes", "--username", "tes", "--age", "100",
	})
	fmt.Println(name)
	fmt.Println(username)
	fmt.Println(age)

}

