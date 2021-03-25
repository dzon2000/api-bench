package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Api struct {
	Url string
	AppKey string
}

type Response struct {
	Body []byte
	Time time.Duration
}

func (api Api) Call() Response {
	client := &http.Client{}
	var start time.Time
	var total time.Duration
	trace := &httptrace.ClientTrace {
		GotFirstResponseByte: func() {
            total = time.Since(start)
        },
	}

	req, _ := http.NewRequest("POST", api.Url, nil)
	req.Header.Add("AppKey", api.AppKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Println(
			err,
			resp.StatusCode,
		)
		return Response{}
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return Response{body, total}
}
