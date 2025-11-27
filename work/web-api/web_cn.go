package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	appId := []int64{
		12,
		13,
		14,
		15,
		16,
		17,
		30,
		39,
		42,
		44,
		52,
		53,
		54,
		60,
		64,
		66,
		70,
		72,
		73,
		74,
		76,
		77,
		80,
		81,
		82,
		83,
		84,
		85,
		86,
		87,
		88,
		89,
		90,
		93,
		94,
		95,
		97,
		99,
		100,
		106,
		110,
		114,
		115,
		116,
		117,
		118,
		121,
		122,
		123,
		124,
		129,
		132,
		133,
		134,
	}

	client := &http.Client{}
	for _, id := range appId {
		url := fmt.Sprintf("https://api-funnydb.zh-cn.xmfunny.com/api/v1/zone-sync/apps/apply-to-resource-manager/%d", id)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Authorization", "Basic ZnVubnlkYjozTjIyR1U4OGZkemF2bEtmeFVMaQ==")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("The appId=", id, string(body))
	}
}
