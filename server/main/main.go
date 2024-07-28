package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func receiveMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// 요청 본문에서 메트릭 데이터를 읽음
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	log.Printf("Received metrics: \n%s", string(body))

	// 메트릭 데이터를 저장하거나 처리하는 로직을 여기에 추가
	fmt.Fprintf(w, "Metrics received successfully")
}

func main() {
	http.HandleFunc("/receive_metrics", receiveMetricsHandler)

	log.Println("Agent server is running on port 9001")
	log.Fatal(http.ListenAndServe(":9001", nil))
}
