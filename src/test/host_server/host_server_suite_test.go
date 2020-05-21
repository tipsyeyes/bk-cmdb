package host_server_test

import (
	"testing"

	"configdatabase/src/test"
	"configdatabase/src/test/reporter"
	"configdatabase/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var hostServerClient = test.GetClientSet().HostServer()
var apiServerClient = test.GetClientSet().ApiServer()
var instClient = test.GetClientSet().TopoServer().Instance()

func TestHostServer(t *testing.T) {
	RegisterFailHandler(util.Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"hostserver.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "HostServer Suite", reporters)
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()
})
