package main

import (
	"bytes"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
	"go.opentelemetry.io/otel/sdk/resource"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

//https://github.com/rs/cors,https://github.com/zeyadkhaled/openversion/blob/874ebfca1df5f5a5937db082b079544ac8412a98/api/api.go

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

	// Create a meter
	exporter, err := prometheus.NewExportPipeline(
		prometheus.Config{},
		pull.WithResource(resource.New(label.String("R", "V"))),
	)
	if err != nil {
		panic(err)
	}
	meter := exporter.Provider().Meter("example")
	ctx := context.Background()


	// Use two instruments
	counter := metric.Must(meter).NewInt64Counter(
		"a.counter",
		metric.WithDescription("Counts things"),
	)

	recorder := metric.Must(meter).NewInt64ValueRecorder(
		"a.valuerecorder",
		metric.WithDescription("Records values"),
	)

	updown := metric.Must(meter).NewInt64UpDownCounter(
			"a.updown",
			metric.WithDescription("Updown values"),
		)

	_ = metric.Must(meter).NewInt64SumObserver("int.sumobserver.sum", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(-1, label.String("A", "B"))
	})

	counter.Add(ctx, 100, label.String("key", "value"))
	counter.Add(ctx, 100, label.String("key", "value"))

	recorder.Record(ctx, 100, label.String("key", "value"))

	updown.Add(ctx, 120, label.String("key", "value"))
	updown.Add(ctx, 10, label.String("key", "value"))
	updown.Add(ctx, -10, label.String("key", "value"))


	// GET the HTTP endpoint
	var input bytes.Buffer
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", &input)
	if err != nil {
		panic(err)
	}
	exporter.ServeHTTP(resp, req)
	data, err := ioutil.ReadAll(resp.Result().Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(data))

	// Output:
	// # HELP a_counter Counts things
	// # TYPE a_counter counter
	// a_counter{R="V",key="value"} 100
	// # HELP a_valuerecorder Records values
	// # TYPE a_valuerecorder histogram
	// a_valuerecorder_bucket{R="V",key="value",le="+Inf"} 1
	// a_valuerecorder_sum{R="V",key="value"} 100
	// a_valuerecorder_count{R="V",key="value"} 1

	//meter := global.MeterProvider().Meter("ex.com/basic")
	//
	//transaction := meter.NewInt64Counter(
	//	"transaction.volume",
	//	metric.WithKeys(key.New("status")),
	//)
	//transactionLabels := meter.Labels(key.String("status", "pending"))
	//
	//ctx := context.Background()
	//
	////rest
	//r := chi.NewRouter()
	//r.Get("/transaction", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("OK"))
	//	// add count
	//	meter.RecordBatch(
	//		ctx,
	//		transactionLabels,
	//		transaction.Measurement(1),
	//	)
	//})
	//
	//http.ListenAndServe(":3000", r)
}
