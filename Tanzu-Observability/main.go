//package main
//
//import (
//	"fmt"
//	"github.com/wavefronthq/go-metrics-wavefront/reporting"
//	"github.com/wavefronthq/wavefront-sdk-go/application"
//	"github.com/wavefronthq/wavefront-sdk-go/senders"
//)
//
//func main() {
//	directCfg := &senders.DirectConfiguration{
//		Server: "https://surf.wavefront.com/api/v2",
//		Token:  "2c827e40-32ba-4aab-9e2f-3db734d425e6",
//	}
//
//	sender, err := senders.NewDirectSender(directCfg)
//	if err != nil {
//		panic(err)
//	}
//
//	var testmap = map[string]string{
//		"a": "1",
//		"b": "2",
//		"c": "3",
//	}
//
//	err = sender.SendMetric("test-metric", 0.001111, 24, "test string", testmap)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	app := application.New("test-app", "test-service")
//	reporter := reporting.NewReporter(
//		sender,
//		app,
//		reporting.Source("go-metrics-test"),
//		reporting.Prefix("test.prefix"),
//		reporting.LogErrors(true),
//	)
//
//	reporter2 := reporting.NewReporter(
//		sender,
//		application.New("test-app2", "test-service"),
//		reporting.Source("go-metrics-test"),
//		reporting.Prefix("test.prefix"),
//		reporting.LogErrors(true),
//		reporting.RuntimeMetric(true),
//	)
//
//	reporting.RuntimeMetric(true)
//	reporter3 := reporting.NewReporter(
//		sender,
//		application.New("test-app", "test-service"),
//		reporting.Source("go-metrics-test"),
//		reporting.Prefix("test.prefix"),
//		reporting.LogErrors(true),
//		reporting.RuntimeMetric(true),
//	)
//
//	fmt.Println(reporter, reporter2, reporter3)
//}
//
////
////tekton: read-writepush : c94bad9e-7ede-4d62-94b1-411009dade3b
////tekton: read-only : d1dc39ba-d269-461f-b374-abc36b4fbf58
////
////docker login -u lnanjangud653
////
////https://jira.eng.vmware.com/browse/MAPBUA-967

////////////////////

package main

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/wavefronthq/go-metrics-wavefront/reporting"
	"log"
	"math/rand"
	"time"

	"github.com/wavefronthq/wavefront-sdk-go/application"
	"github.com/wavefronthq/wavefront-sdk-go/event"

	"github.com/wavefronthq/wavefront-sdk-go/senders"
)

func main() {
	var wfSenders []senders.Sender

	//urls := strings.Split(os.Getenv("WF_URL"), "|")
	//for _, url := range urls {
	//	sender, err := senders.NewSender(url)
	//	if err != nil {
	//		panic(err)
	//	}
	//	wfSenders = append(wfSenders, sender)
	//}

	//// OLD PROXY way
	//proxyCfg := &senders.ProxyConfiguration{
	//	Host:                 "localhost",
	//	MetricsPort:          2878,
	//	DistributionPort:     2878,
	//	TracingPort:          2878,
	//	EventsPort:           2878,
	//	FlushIntervalSeconds: 10,
	//}
	//
	//sender, err := senders.NewProxySender(proxyCfg)
	//if err != nil {
	//	panic(err)
	//}
	//wfSenders = append(wfSenders, sender)
	//
	//// OLD DIRECT way
	//directCfg := &senders.DirectConfiguration{
	//	Server:               "https://-----.wavefront.com",
	//	Token:                "--------------",
	//	BatchSize:            10000,
	//	MaxBufferSize:        500000,
	//	FlushIntervalSeconds: 1,
	//}

	directCfg := &senders.DirectConfiguration{
		Server: "https://surf.wavefront.com",
		Token:  "2c827e40-32ba-4aab-9e2f-3db734d425e6",
	}

	sender, err := senders.NewDirectSender(directCfg)
	if err != nil {
		panic(err)
	}
	wfSenders = append(wfSenders, sender)

	wf := senders.NewMultiSender(wfSenders...)
	log.Print("senders ready")

	source := "go_sdk_example"

	reporter := reporting.NewMetricsReporter(
		sender,
		reporting.ApplicationTag(application.New("metric-app", "srv")),
		reporting.Source("go-metrics-test"),
		reporting.Prefix("some.prefix"),
		reporting.RuntimeMetric(true),
	)

	tags := map[string]string{
		"key2": "val2",
		"key1": "val1",
		"key0": "val0",
		"key4": "val4",
		"key3": "val3",
	}

	// Create a counter metric and register with tags
	counter := metrics.NewCounter()
	reporter.RegisterMetric("foo", counter, tags)
	counter.Inc(47)

	// Create a histogram and register with tags
	histogram := reporting.NewHistogram()
	reporter.RegisterMetric("duration", histogram, tags)

	// Create a histogram and register without tags
	histogram2 := reporting.NewHistogram()
	// reporter.Register("duration2", histogram2)

	deltaCounter := metrics.NewCounter()
	reporter.RegisterMetric(reporting.DeltaCounterName("delta.metric"), deltaCounter, tags)
	deltaCounter.Inc(10)

	fmt.Println("Search wavefront: ts(\"test.prefix.foo.count\")")
	fmt.Println("Entering loop to simulate metrics flushing. Hit ctrl+c to cancel")

	for {
		counter.Inc(rand.Int63())
		histogram.Update(rand.Int63())
		histogram2.Update(rand.Int63())
		deltaCounter.Inc(10)
		time.Sleep(time.Second * 2)
	}

	app := application.New("sample app", "main.go")
	application.StartHeartbeatService(wf, app, source)

	tags = make(map[string]string)
	tags["namespace"] = "default"
	tags["Kind"] = "Deployment"

	options := []event.Option{event.Details("Details"), event.Type("type"), event.Severity("severity")}

	for i := 0; i < 10; i++ {
		err := wf.SendMetric("sample.metric", float64(i), time.Now().UnixNano(), source, map[string]string{"env": "test"})
		if err != nil {
			println("error:", err.Error())
		}

		txt := fmt.Sprintf("test event %d", i)
		sendEvent(wf, txt, time.Now().Unix(), 0, source, tags, options...)

		time.Sleep(10 * time.Second)
	}

	wf.Flush()
	wf.Close()
}

func sendEvent(sender senders.Sender, name string, startMillis, endMillis int64, source string, tags map[string]string, setters ...event.Option) {
	err := sender.SendEvent(name, startMillis, endMillis, source, tags, setters...)
	if err != nil {
		println("error:", err)
	}
}
