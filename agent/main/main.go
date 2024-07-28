package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// 타겟 서버와 에이전트 서버의 URL 설정
const (
	targetMetricsURL = "http://localhost:9000/metrics"
	agentServerURL   = "http://localhost:9001/receive_metrics"
)

// 주기적으로 메트릭을 수집하는 함수
func collectMetrics() {
	for {
		// 타겟 서버에서 메트릭 데이터를 수집
		resp, err := http.Get(targetMetricsURL)
		if err != nil {
			log.Printf("Error getting metrics: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Received non-200 response: %d", resp.StatusCode)
			resp.Body.Close()
			time.Sleep(10 * time.Second)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		log.Printf("Metrics collected: \n%s", string(body))

		// 수집한 메트릭 데이터를 에이전트 서버로 전송
		postMetricsToAgent(body)

		// 10초마다 메트릭을 수집
		time.Sleep(10 * time.Second)
	}
}

// 에이전트 서버로 메트릭 데이터를 전송하는 함수
func postMetricsToAgent(metrics []byte) {
	resp, err := http.Post(agentServerURL, "text/plain", bytes.NewBuffer(metrics))
	if err != nil {
		log.Printf("Error posting metrics to agent: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response from agent: %d", resp.StatusCode)
	}

	resp.Body.Close()
}

func main() {
	// 메트릭 수집 시작
	go collectMetrics()

	// 간단한 HTTP 서버를 실행하여 수집이 제대로 되는지 확인
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Metrics collector is running")
	})

	log.Println("Metrics collector server is running on port 9002")
	log.Fatal(http.ListenAndServe(":9002", nil))
}
