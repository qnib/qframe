# filter-grok
Filter to apply GROK statements against a message.

## Hello World

As a standalone program like this:

```go
func main() {
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"filter.test.pattern": "%{INT:number}",
		"filter.test.inputs": "test",
	}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	p, err := qframe_filter_grok.New(qChan, *cfg, "test")
	if err != nil {
		log.Printf("[EE] Failed to create filter: %v", err)
		return
	}
	go p.Run()
	time.Sleep(2*time.Second)
	bg := qChan.Data.Join()
	qm := qtypes.NewQMsg("test", "test")
	qm.Msg = "1"
	log.Println("Send message")
	qChan.Data.Send(qm)
	for {
		qm = bg.Recv().(qtypes.QMsg)
		if qm.Source == "test" {
			continue
		}
		fmt.Printf("#### Received result from grok (pattern:%s) filter for input: %s\n", p.GetPattern(), qm.Msg)
		for k, v := range qm.KV {
			fmt.Printf("%+15s: %s\n", k, v)
		}
		break

	}
}
```

The plugin produces the following outcome:

```bash
$ go run main.go
2017/04/14 07:28:16 [II] Dispatch broadcast for Data and Tick
2017/04/14 07:28:16 [II] Start grok filter 'test' v0.0.0
2017/04/14 07:28:18 Send message
#### Received result from grok (pattern:%{INT:number}) filter for input: 1
         number: 1
```

## Benchmark

```bash
$ go test -bench=Grok  -benchtime=5s
BenchmarkGrok-4   	  200000	     55903 ns/op
PASS
ok  	github.com/qframe/filter-grok	12.040s
```
