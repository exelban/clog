benchWithTraceAllocation:
	go test -c
	GODEBUG=allocfreetrace=1 ./logg.test -test.run=none -test.bench=BenchmarkLogg_Write -test.benchtime=10ms 2>trace.log

saveBenchOld:
	go test -bench=BenchmarkLogg_Write | tee old.txt

compare:
	go test -bench=BenchmarkLogg_Write > new.txt
	benchcmp old.txt new.txt
