package util

import (
	"configdatabase/src/common/json"
	"encoding/base64"
	"fmt"
)

func GetTokenSuper(token string) string {
	super := "0"
	t := SplitStrField(token, ".")
	if len(t) != 3 {
		return super
	}

	decodeBytes, err := base64.RawURLEncoding.DecodeString(t[1])
	if err != nil {
		return super
	}

	userInfo := struct {
		Username 		string		`json:"username"`
		IsSuperuser		bool		`json:"is_superuser"`
		TrueName 		string 		`json:"true_name"`
	}{}

	if err = json.Unmarshal(decodeBytes, &userInfo); err != nil {
		return super
	}

	if userInfo.IsSuperuser {
		super = "1"
	}

	fmt.Println(userInfo)
	return super
}