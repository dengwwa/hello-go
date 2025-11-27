package main

//
//import (
//	"fmt"
//	"io/ioutil"
//	"net/http"
//)
//
//func main() {
//
//	appId := []int64{
//		//1,
//		//2,
//		//3,
//		4,
//		5,
//		6,
//		7,
//		9,
//		10,
//		11,
//		21,
//		22,
//		31,
//		32,
//		45,
//		46,
//		49,
//		51,
//		55,
//		57,
//		58,
//		59,
//		61,
//		62,
//		69,
//		71,
//		75,
//		78,
//		91,
//		92,
//		98,
//		107,
//		109,
//		111,
//		112,
//		119,
//		120,
//		125,
//		126,
//		128,
//		130,
//		131,
//	}
//
//	client := &http.Client{}
//	for _, id := range appId {
//		url := fmt.Sprintf("https://api-funnydb.sg.xmfunny.com/api/v1/zone-sync/apps/apply-to-resource-manager/%d", id)
//		req, err := http.NewRequest("GET", url, nil)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		req.Header.Add("Authorization", "Basic ZnVubnlkYjo0bHpiczk4SXJpTjNVUmJTSFhHYQ==")
//
//		res, err := client.Do(req)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		defer res.Body.Close()
//
//		body, err := ioutil.ReadAll(res.Body)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		fmt.Println("The appId=", id, string(body))
//	}
//}
