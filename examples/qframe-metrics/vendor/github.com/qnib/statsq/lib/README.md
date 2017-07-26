## Benchmark

```
$ go test -bench=Bench .
BenchmarkOneBigTimer-4            	       1	1390692436 ns/op
BenchmarkLotsOfTimers-4           	       1	1528838842 ns/op
BenchmarkParseLineCounter-4       	 1000000	      1401 ns/op
BenchmarkParseLineGauge-4         	 1000000	      1303 ns/op
BenchmarkParseLineTimer-4         	 1000000	      1373 ns/op
BenchmarkParseLineSet-4           	 1000000	      1256 ns/op
BenchmarkPacketHandlerCounter-4   	10000000	       137 ns/op
BenchmarkPacketHandlerGauge-4     	10000000	       145 ns/op
BenchmarkPacketHandlerTimer-4     	10000000	       183 ns/op
BenchmarkPacketHandlerSet-4       	10000000	       197 ns/op
PASS
ok  	github.com/ChristianKniep/statsq/lib	17.375s
```

## Testcases

```
$ go test -cover .
ok  	github.com/ChristianKniep/statsq/lib	0.011s	coverage: 67.2% of statements 
```
