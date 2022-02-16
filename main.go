package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	kitmetrics "github.com/go-kit/kit/metrics"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	namespace, subsystem := "iam", "user_sync_kafka"
	//	fieldKeys := []string{"method", "error", "code"}
	fieldKeys2 := []string{"error"}

	crmPushCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "crm_push_total",
		Help:      "Total number of requests to CRM sent.",
	}, fieldKeys2)
	crmPushLatency := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "crm_push_time_seconds",
		Help:      "Response from CRM time in seconds.",
	}, fieldKeys2)

	totalProcessTime := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "total_process_time_seconds",
		Help:      "Time difference between message write to kafka and read from kafka.",
	}, fieldKeys2)

	kafkaReadWriteLatency := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "kafka_read_write_time_seconds",
		Help:      "Time difference between message write to kafka and read from kafka.",
	}, fieldKeys2)

	userEventsCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "user_events_total",
		Help:      "Total number of user events.",
	}, fieldKeys2)

	go createFakeMetrics(crmPushLatency, kafkaReadWriteLatency, totalProcessTime, crmPushCount, userEventsCount)

	startPromServer()
}

/*
	kafkaReadWriteLatency kitmetrics.Histogram,
	totalProcessTime kitmetrics.Histogram,
	crmPushCount kitmetrics.Counter,
	crmPushLatency kitmetrics.Histogram,
	userEventsCount kitmetrics.Counter
*/
func createFakeMetrics(crmpl kitmetrics.Histogram, krwl kitmetrics.Histogram, tpt kitmetrics.Histogram, crmpc kitmetrics.Counter, uec kitmetrics.Counter) {

	for {
		startTime := time.Now()

		max1 := 300
		min1 := 10
		time.Sleep(time.Duration(rand.Intn(max1-min1)+min1) * time.Millisecond)

		readFromKafkaTime := time.Now()

		max2 := 500
		min2 := 50
		time.Sleep(time.Duration(rand.Intn(max2-min2)+min2) * time.Millisecond)

		startWriteToMangoTime := time.Now()

		max3 := 3000
		min3 := 1000
		time.Sleep(time.Duration(rand.Intn(max3-min3)+min3) * time.Millisecond)

		// num1 := rand.Int()
		// num2 := rand.Int()

		labels := []string{
			"error", "none",
		}

		// curTime := time.Now()

		// past1 := curTime.Add(-time.Duration(num1) * time.Second) //twk = time.Now().Add(-5 * time.Second)
		// past2 := curTime.Add(-time.Duration(num2) * time.Second)

		// past3 := curTime.Add(-time.Duration(num1+num2) * time.Second)

		crmpl.With(labels...).Observe(time.Since(startWriteToMangoTime).Seconds())
		krwl.With(labels...).Observe(float64(readFromKafkaTime.Sub(startTime).Seconds())) //).Since().Seconds())
		tpt.With(labels...).Observe(time.Since(startTime).Seconds())

		crmpc.With(labels...).Add(1)
		uec.With(labels...).Add(1)

		time.Sleep(1 * time.Second)
		// c.crmPushCount.With(lvs...).Add(1)
		// c.crmPushLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		// c.totalProcessTime.With(lvs2...).Observe(time.Since(job.TimeWrittenToKafka).Seconds())

	}

}

func startPromServer() {

	metricsAddr := ":8081"
	// Serve prometheus "/metrics" endpoint
	//	go func() {
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("logging metrics")
	//		logger.Log("endpoint", "/metrics", "address", metricsAddr, "msg", "listening")
	err := http.ListenAndServe(metricsAddr, nil)
	if err != nil {
		fmt.Printf("failed to log metrics: %v\n", err)
		//			logger.Log("endpoint", "/metrics", "during", "listen", "err", err)
		os.Exit(1)
	}
	//	}()

}
