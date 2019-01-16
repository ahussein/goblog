package service

import (
	"log"
	"net/http"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func StartWebServer(port string) {
	log.Println("Starting HTTP service at " + port)
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "accountservice",
	})
	if err != nil {
		log.Fatalf("Failed to create Promethues exporter: %v", err)
	}
	view.RegisterExporter(pe)
	r := NewRouter()
	r.Methods("GET", "POST").
		Path("/metrics").
		Name("Metrics").
		Handler(pe)

	http.Handle("/", r)
	h := &ochttp.Handler{Handler: r}
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		log.Fatal("Fail to register ochttp.DefaultServerViews")
	}
	// register client views
	if err = view.Register(
		// register few default views
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,

		// register custom view
		&view.View{
			Name:        "httpclient_latency_by_path",
			TagKeys:     []tag.Key{ochttp.KeyClientPath},
			Measure:     ochttp.ClientRoundtripLatency,
			Aggregation: ochttp.DefaultLatencyDistribution,
		},
	); err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":"+port, h)
	if err != nil {
		log.Println("An error occured start HTTP listener at port " + port)
		log.Println("Error: " + err.Error())
	}
}
