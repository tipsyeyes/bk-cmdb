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


func TestHostSearch(t *testing.T) {
	// 查找主机
	condition := map[string]interface{}{
		"$or": []map[string]interface{}{
			{"bk_sn": map[string]interface{}{"$regex": "nidetian"}},
			{"bk_host_name": map[string]interface{}{"$regex": "hostname-aa"}},
		},
	}
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
	fmt.Println("host count: ", gResult.Data.Count)

}
