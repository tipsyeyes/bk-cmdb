package operation

import (
	"context"
	"net/http"

	"configdatabase/src/apimachinery/rest"
	"configdatabase/src/common/metadata"
)

type OperationClientInterface interface {
	SearchChartDataCommon(ctx context.Context, h http.Header, data metadata.ChartConfig) (resp *metadata.Response, err error)
	SearchInstCount(ctx context.Context, h http.Header, data interface{}) (resp *metadata.CoreUint64Response, err error)
	CreateOperationChart(ctx context.Context, h http.Header, data interface{}) (resp *metadata.CoreUint64Response, err error)
	SearchOperationChart(ctx context.Context, h http.Header, data interface{}) (resp *metadata.SearchChartResponse, err error)
	DeleteOperationChart(ctx context.Context, h http.Header, data string) (resp *metadata.Response, err error)
	UpdateOperationChart(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	SearchTimerChartData(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	UpdateChartPosition(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)
	SearchChartCommon(ctx context.Context, h http.Header, data interface{}) (resp *metadata.SearchChartCommon, err error)
	TimerFreshData(ctx context.Context, h http.Header, data interface{}) (resp *metadata.BoolResponse, err error)
}

func NewOperationClientInterface(client rest.ClientInterface) OperationClientInterface {
	return &operation{client: client}
}

type operation struct {
	client rest.ClientInterface
}
