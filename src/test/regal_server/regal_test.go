package regal_server_test

import (
	"context"
	"fmt"
	"testing"

	"configdatabase/src/common/metadata"
	"configdatabase/src/test"
)

/*
db.getCollection("cc_HostBase").find({
    "$or": [
        {"bk_sn": {"$regex": "nidetian"}}
    ],
	"$or": [
        {"bk_host_name": {"$regex": "hostname-aa"}}
    ]
})
db.getCollection("cc_HostBase").find({
    "$or": [
        {"bk_sn": {"$regex": "nidetian"}},
        {"bk_host_name": {"$regex": "hostname-aa"}}
    ]
})
*/

// 测试 core service主机查询
func TestCoreServiceHostSearch(t *testing.T) {
	// 查找主机
	condition := map[string]interface{}{
		"$or": []map[string]interface{}{
			{"bk_sn": map[string]interface{}{"$regex": "nidetian"}},
			{"bk_host_name": map[string]interface{}{"$regex": "hostname-aa"}},
		},
	}

	// 问题 3.6.1的版本，当多个包含查询时，只保存了最后的 $or条件，导致查询结果异常
	query := &metadata.QueryInput{
		Condition: condition,
		Start:     0,
		Limit:     10000,
		Sort:      "bk_host_id",
		Fields:    "",
	}

	gResult, err := test.GetClientSet().CoreService().Host().GetHosts(context.Background(), test.GetHeader(), query)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("\nhost count: ", gResult.Data.Count)

}