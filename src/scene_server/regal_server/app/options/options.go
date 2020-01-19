// app options info

package options

import (
	"configdatabase/src/auth/authcenter"
	"configdatabase/src/common/auth"
	"configdatabase/src/common/core/cc/config"
	"configdatabase/src/storage/dal/mongo"
	"configdatabase/src/storage/dal/redis"

	"github.com/spf13/pflag"
)

// ServerOption
// 服务器信息地址信息
type ServerOption struct {
	ServConf *config.CCAPIConfig

	// TODO: 下面可以添加自定义命令行参数
}

// NewServerOption create a ServerOption object
func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf: config.NewCCAPIConfig(),
	}

	return &s
}

// AddFlags add flags
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServConf.AddrPort, "addrport", "127.0.0.1:65001", "The ip address and port for the serve on")
	// fs.UintVar(&s.ServConf.Port, "port", 60003, "The port for the serve on")
	fs.StringVar(&s.ServConf.RegDiscover, "regdiscv", "", "hosts of register and discover server. e.g: 127.0.0.1:2181")
	fs.StringVar(&s.ServConf.ExConfig, "config", "", "The config path. e.g conf/api.conf")
	fs.Var(auth.EnableAuthFlag, "enable-auth", "The auth center enable status, true for enabled, false for disabled")
}

// Config 主Server配置
// 用于接收 zk同步的配置文件信息
type Config struct {
	Redis redis.Config
	Mongo mongo.Config

	// 权限中心配置
	Auth  authcenter.AuthConfig
}
