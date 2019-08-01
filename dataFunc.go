package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)


func PushIt(L []*MetricValue) error {
	postThing , err:= json.Marshal(L)
	if err != nil {
		return err
	}
	fmt.Println(string(postThing))
	url := "http://127.0.0.1:1988/v1/push"
	done := make(chan [2] interface{}, 1)
	go func() {
		var result  [2] interface{}
		resp, err := http.Post(url,
			"application/x-www-form-urlencoded",
			strings.NewReader(string(postThing)))

		result[0] = resp
		result[1] = err
		done <- result

	}()

	select {
	case <-time.After(time.Duration(5 * time.Second)):
		fmt.Println("http call timeout")
		return errors.New(" http call timeout")
	case result := <-done:
		if result[1] != nil{
			fmt.Println(result[1])
		}else{
			fmt.Println(result[0])
		}
	}


	return nil
}


func pushIt(value, timestamp, metric, tags, containerId, counterType, endpoint string) error {

	postThing := `[{"metric":"` + metric + `","endpoint":"` + endpoint + `", "timestamp":` + timestamp + `,"step":` + "60" + `,"value":` + value + `,"counterType":"` + counterType + `","tags":"` + tags + `"}]`
	fmt.Println(postThing)
	LogRun(postThing)
	//push data to falcon-agent
	url := "http://127.0.0.1:1988/v1/push"
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(postThing))
	if err != nil {
		LogErr(err, "Post err in pushIt")
		return err
	}
	defer resp.Body.Close()
	_, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		LogErr(err1, "ReadAll err in pushIt")
		return err1
	}
	return nil
}
