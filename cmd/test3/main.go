package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
)

var (
	lemonsKey = label.Key("ex.com/lemons")
)

func initMeter() {
	exporter, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	http.HandleFunc("/", exporter.ServeHTTP)
	go func() {
		_ = http.ListenAndServe(":2222", nil)
	}()


	fmt.Println("Prometheus server running on :2222")
}

func main() {
	initMeter()

	meter := global.Meter("ex.com/basic")
	observerLock := new(sync.RWMutex)
	observerValueToReport := new(float64)
	observerLabelsToReport := new([]label.KeyValue)
	cb := func(_ context.Context, result metric.Float64ObserverResult) {
		(*observerLock).RLock()
		value := *observerValueToReport
		labels := *observerLabelsToReport
		(*observerLock).RUnlock()
		result.Observe(value, labels...)
	}
	_ = metric.Must(meter).NewFloat64ValueObserver("ex.com.one", cb,
		metric.WithDescription("A ValueObserver set to 1.0"),
	)

	valuerecorder := metric.Must(meter).NewFloat64ValueRecorder("ex.com.two")
	counter := metric.Must(meter).NewFloat64Counter("ex.com.three")

	commonLabels := []label.KeyValue{lemonsKey.Int(10), label.String("A", "1"), label.String("B", "2"), label.String("C", "3")}
	notSoCommonLabels := []label.KeyValue{lemonsKey.Int(13)}

	ctx := context.Background()

	(*observerLock).Lock()
	*observerValueToReport = 1.0
	*observerLabelsToReport = commonLabels
	(*observerLock).Unlock()
	meter.RecordBatch(
		ctx,
		commonLabels,
		valuerecorder.Measurement(2.0),
		counter.Measurement(12.0),
	)

	time.Sleep(5 * time.Second)

	(*observerLock).Lock()
	*observerValueToReport = 1.0
	*observerLabelsToReport = notSoCommonLabels
	(*observerLock).Unlock()
	meter.RecordBatch(
		ctx,
		notSoCommonLabels,
		valuerecorder.Measurement(2.0),
		counter.Measurement(22.0),
	)

	time.Sleep(5 * time.Second)

	(*observerLock).Lock()
	*observerValueToReport = 13.0
	*observerLabelsToReport = commonLabels
	(*observerLock).Unlock()
	meter.RecordBatch(
		ctx,
		commonLabels,
		valuerecorder.Measurement(12.0),
		counter.Measurement(13.0),
	)



	fmt.Println("Example finished updating, please visit :2222")

	select {}
}