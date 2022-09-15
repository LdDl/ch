My PC is:

    Processor: Intel(R) Core(TM) i5-10600K CPU @ 4.10GHz x 12
    Memory: 16.3GiB
    Linux Kernel: 5.15.0-47-lowlatency
    OS: Ubuntu 22.04 LTS

I have used graph with ~187k vertices for benchmark.

For one-to-one query ([ShortestPath](bidirectional_ch.go#L16)):
```bash
goos: linux
goarch: amd64
pkg: github.com/LdDl/ch
cpu: Intel(R) Core(TM) i5-10600K CPU @ 4.10GHz
BenchmarkShortestPath
    bidirectional_ch_test.go:71: BenchmarkShortestPath is starting...
BenchmarkShortestPath/CH_shortest_path/4/vertices-4-edges-9-shortcuts-1
BenchmarkShortestPath/CH_shortest_path/4/vertices-4-edges-9-shortcuts-1-12               1957400             649.5 ns/op           345 B/op         14 allocs/op
BenchmarkShortestPath/CH_shortest_path/8/vertices-8-edges-63-shortcuts-8
BenchmarkShortestPath/CH_shortest_path/8/vertices-8-edges-63-shortcuts-8-12               741454              1592 ns/op           735 B/op         25 allocs/op
BenchmarkShortestPath/CH_shortest_path/16/vertices-16-edges-313-shortcuts-25
BenchmarkShortestPath/CH_shortest_path/16/vertices-16-edges-313-shortcuts-25-12           269764              4262 ns/op           1721 B/op        43 allocs/op
BenchmarkShortestPath/CH_shortest_path/32/vertices-32-edges-1396-shortcuts-104
BenchmarkShortestPath/CH_shortest_path/32/vertices-32-edges-1396-shortcuts-104-12          83739             13438 ns/op           4981 B/op        89 allocs/op
BenchmarkShortestPath/CH_shortest_path/64/vertices-64-edges-5882-shortcuts-382
BenchmarkShortestPath/CH_shortest_path/64/vertices-64-edges-5882-shortcuts-382-12          33704             36579 ns/op          11221 B/op       171 allocs/op
BenchmarkShortestPath/CH_shortest_path/128/vertices-128-edges-24050-shortcuts-1456
BenchmarkShortestPath/CH_shortest_path/128/vertices-128-edges-24050-shortcuts-1456-12      10416            109526 ns/op          25132 B/op       348 allocs/op
BenchmarkShortestPath/CH_shortest_path/256/vertices-256-edges-97234-shortcuts-5719
BenchmarkShortestPath/CH_shortest_path/256/vertices-256-edges-97234-shortcuts-5719-12       4414            286111 ns/op          53513 B/op       722 allocs/op
PASS
ok      github.com/LdDl/ch      72.896s
```

For one-to-many query ([ShortestPathOneToMany](bidirectional_ch_one_to_n.go#L15)):
```bash
goos: linux
goarch: amd64
pkg: github.com/LdDl/ch
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/1/vertices-187853-20         	     160	   7175512 ns/op	 6773392 B/op	   14463 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/2/vertices-187853-20         	     176	   6847237 ns/op	 6773416 B/op	   14463 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/4/vertices-187853-20         	     172	   7101499 ns/op	 6773230 B/op	   14461 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/8/vertices-187853-20         	     181	   6706642 ns/op	 6773435 B/op	   14463 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/16/vertices-187853-20        	     170	   6915546 ns/op	 6773325 B/op	   14462 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/32/vertices-187853-20        	     174	   6887815 ns/op	 6773307 B/op	   14462 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/64/vertices-187853-20        	     174	   6964305 ns/op	 6773370 B/op	   14462 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/128/vertices-187853-20       	     168	   6916208 ns/op	 6773333 B/op	   14463 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/256/vertices-187853-20       	     170	   7161520 ns/op	 6773373 B/op	   14463 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/512/vertices-187853-20       	     172	   6710753 ns/op	 6773492 B/op	   14464 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/1024/vertices-187853-20      	     181	   6680762 ns/op	 6773273 B/op	   14462 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/2048/vertices-187853-20      	     171	   6695043 ns/op	 6773313 B/op	   14462 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/4096/vertices-187853-20      	     176	   6674091 ns/op	 6773373 B/op	   14462 allocs/op
PASS
ok  	github.com/LdDl/ch	33.806s
```

Also if you want to make comparison between OneToMany in term of ShortestPathOneToMany() and OneToMany in term of looping:
```go
go test -benchmem -run=^$ github.com/LdDl/ch -bench BenchmarkOldWayShortestPathOneToMany > old.txt
go test -benchmem -run=^$ github.com/LdDl/ch -bench BenchmarkShortestPathOneToMany > new.txt
sed -i 's/BenchmarkOldWayShortestPathOneToMany/BenchmarkShortestPathOneToMany/g' old.txt
```
and then use [benchcmp](https://godoc.org/golang.org/x/tools/cmd/benchcmp):
```bash
benchcmp old.txt new.txt
```
Output should be something like this:
```bash
benchmark                                                                                 old ns/op     new ns/op     delta
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/1/vertices-187853-20        10608955      7175512       -32.36%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/2/vertices-187853-20        10813368      6847237       -36.68%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/4/vertices-187853-20        10583636      7101499       -32.90%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/8/vertices-187853-20        10500989      6706642       -36.13%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/16/vertices-187853-20       10470206      6915546       -33.95%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/32/vertices-187853-20       10421460      6887815       -33.91%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/64/vertices-187853-20       10499903      6964305       -33.67%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/128/vertices-187853-20      10735268      6916208       -35.57%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/256/vertices-187853-20      10836504      7161520       -33.91%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/512/vertices-187853-20      10544817      6710753       -36.36%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/1024/vertices-187853-20     10619897      6680762       -37.09%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/2048/vertices-187853-20     10772554      6695043       -37.85%
BenchmarkShortestPathOneToMany/CH_shortest_path_(one_to_many)/4096/vertices-187853-20     10257450      6674091       -34.93%
```