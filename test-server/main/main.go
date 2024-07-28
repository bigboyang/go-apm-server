package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Metrics 구조체 정의
type Metrics struct {
	sync.Mutex
	RequestCount     int
	TotalDuration    time.Duration
	RequestDurations []time.Duration
}

// 메트릭 인스턴스 생성
var metrics = Metrics{
	RequestDurations: make([]time.Duration, 0),
}

// 메트릭을 업데이트하는 함수
func (m *Metrics) Record(duration time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.RequestCount++
	m.TotalDuration += duration
	m.RequestDurations = append(m.RequestDurations, duration)
}

// 메트릭 데이터를 문자열로 반환하는 함수
func (m *Metrics) String() string {
	m.Lock()
	defer m.Unlock()
	avgDuration := time.Duration(0)
	if m.RequestCount > 0 {
		avgDuration = m.TotalDuration / time.Duration(m.RequestCount)
	}
	return fmt.Sprintf("Request Count: %d\nTotal Duration: %s\nAverage Duration: %s\n", m.RequestCount, m.TotalDuration, avgDuration)
}

// HTTP 요청을 처리하는 핸들러
func helloHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Write([]byte("Hello, World!"))
	duration := time.Since(start)
	metrics.Record(duration)
}

// 메트릭을 노출하는 핸들러
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(metrics.String()))
}

// 서버 설정 및 실행
func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/metrics", metricsHandler)

	fmt.Println("Server is running on port 9000")
	http.ListenAndServe(":9000", nil)
}
