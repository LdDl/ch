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
cpu: Intel(R) Core(TM) i5-10600K CPU @ 4.10GHz
BenchmarkShortestPathOneToMany/CH_shortest_path/4/vertices-4-edges-9-shortcuts-1-12         	  394494	      3470 ns/op	    2041 B/op	      68 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path/8/vertices-8-edges-62-shortcuts-6-12        	  131952	      9674 ns/op	    4276 B/op	     136 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path/16/vertices-16-edges-312-shortcuts-39-12    	   52059	     22458 ns/op	    8392 B/op	     226 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path/32/vertices-32-edges-1390-shortcuts-109-12  	   18764	     65199 ns/op	   18632 B/op	     427 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path/64/vertices-64-edges-5845-shortcuts-602-12  	    6642	    172839 ns/op	   40654 B/op	     834 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path/128/vertices-128-edges-24033-shortcuts-1257-12         	    2407	    446408 ns/op	   82683 B/op	    1584 allocs/op
BenchmarkShortestPathOneToMany/CH_shortest_path/256/vertices-256-edges-97037-shortcuts-5233-12         	     798	   1330968 ns/op	  184523 B/op	    3461 allocs/op
PASS
ok  	github.com/LdDl/ch	68.262s
```

If you want to make comparison between OneToMany in term of ShortestPathOneToMany() and OneToMany in term of looping:
```go
go test -benchmem -run=^$ github.com/LdDl/ch -bench BenchmarkOldWayShortestPathOneToMany > old.txt
go test -benchmem -run=^$ github.com/LdDl/ch -bench BenchmarkShortestPathOneToMany > new.txt
sed -i 's/BenchmarkOldWayShortestPathOneToMany/BenchmarkShortestPathOneToMany/g' old.txt
```
and then use [benchcmp](https://godoc.org/golang.org/x/tools/cmd/benchcmp):
```bash
go install golang.org/x/tools/cmd/benchcmp@latest
benchcmp old.txt new.txt
```

Output should be something like this:
```bash
benchmark                                                                                          old ns/op     new ns/op     delta
BenchmarkShortestPathOneToMany/CH_shortest_path/4/vertices-4-edges-9-shortcuts-1-12                3298          3388          +2.73%
BenchmarkShortestPathOneToMany/CH_shortest_path/8/vertices-8-edges-61-shortcuts-1-12               7029          7453          +6.03%
BenchmarkShortestPathOneToMany/CH_shortest_path/16/vertices-16-edges-316-shortcuts-31-12           24339         23665         -2.77%
BenchmarkShortestPathOneToMany/CH_shortest_path/32/vertices-32-edges-1404-shortcuts-123-12         60482         59596         -1.46%
BenchmarkShortestPathOneToMany/CH_shortest_path/64/vertices-64-edges-5894-shortcuts-322-12         159670        157954        -1.07%
BenchmarkShortestPathOneToMany/CH_shortest_path/128/vertices-128-edges-23977-shortcuts-1315-12     493555        529951        +7.37%
BenchmarkShortestPathOneToMany/CH_shortest_path/256/vertices-256-edges-97227-shortcuts-5276-12     1409281       1338109       -5.05%

benchmark                                                                                          old allocs     new allocs     delta
BenchmarkShortestPathOneToMany/CH_shortest_path/4/vertices-4-edges-9-shortcuts-1-12                72             68             -5.56%
BenchmarkShortestPathOneToMany/CH_shortest_path/8/vertices-8-edges-61-shortcuts-1-12               121            117            -3.31%
BenchmarkShortestPathOneToMany/CH_shortest_path/16/vertices-16-edges-316-shortcuts-31-12           239            232            -2.93%
BenchmarkShortestPathOneToMany/CH_shortest_path/32/vertices-32-edges-1404-shortcuts-123-12         432            407            -5.79%
BenchmarkShortestPathOneToMany/CH_shortest_path/64/vertices-64-edges-5894-shortcuts-322-12         844            792            -6.16%
BenchmarkShortestPathOneToMany/CH_shortest_path/128/vertices-128-edges-23977-shortcuts-1315-12     1725           1645           -4.64%
BenchmarkShortestPathOneToMany/CH_shortest_path/256/vertices-256-edges-97227-shortcuts-5276-12     3646           3494           -4.17%

benchmark                                                                                          old bytes     new bytes     delta
BenchmarkShortestPathOneToMany/CH_shortest_path/4/vertices-4-edges-9-shortcuts-1-12                1725          2041          +18.32%
BenchmarkShortestPathOneToMany/CH_shortest_path/8/vertices-8-edges-61-shortcuts-1-12               3409          3607          +5.81%
BenchmarkShortestPathOneToMany/CH_shortest_path/16/vertices-16-edges-316-shortcuts-31-12           9912          8745          -11.77%
BenchmarkShortestPathOneToMany/CH_shortest_path/32/vertices-32-edges-1404-shortcuts-123-12         24182         18083         -25.22%
BenchmarkShortestPathOneToMany/CH_shortest_path/64/vertices-64-edges-5894-shortcuts-322-12         55777         38899         -30.26%
BenchmarkShortestPathOneToMany/CH_shortest_path/128/vertices-128-edges-23977-shortcuts-1315-12     124457        85822         -31.04%
BenchmarkShortestPathOneToMany/CH_shortest_path/256/vertices-256-edges-97227-shortcuts-5276-12     271506        186313        -31.38
```


If you want to make comparison between OneToMany in term of ShortestPathOneToMany() and OneToMany in term of looping:
```go
go test -benchmem -run=^$ github.com/LdDl/ch -bench BenchmarkOldWayShortestPathManyToMany > old_m_n.txt
go test -benchmem -run=^$ github.com/LdDl/ch -bench BenchmarkShortestPathManyToMany > new_m_n.txt
sed -i 's/BenchmarkOldWayShortestPathManyToMany/BenchmarkShortestPathManyToMany/g' old_m_n.txt
```
and then use [benchcmp](https://godoc.org/golang.org/x/tools/cmd/benchcmp):
```bash
go install golang.org/x/tools/cmd/benchcmp@latest
benchcmp old_m_n.txt new_m_n.txt
```

Output should be something like this:
```bash
benchmark                                                                                           old ns/op     new ns/op     delta
BenchmarkShortestPathManyToMany/CH_shortest_path/4/vertices-4-edges-9-shortcuts-1-12                3087          4593          +48.79%
BenchmarkShortestPathManyToMany/CH_shortest_path/8/vertices-8-edges-61-shortcuts-1-12               6907          7927          +14.77%
BenchmarkShortestPathManyToMany/CH_shortest_path/16/vertices-16-edges-316-shortcuts-31-12           23256         23335         +0.34%
BenchmarkShortestPathManyToMany/CH_shortest_path/32/vertices-32-edges-1404-shortcuts-123-12         67078         72903         +8.68%
BenchmarkShortestPathManyToMany/CH_shortest_path/64/vertices-64-edges-5894-shortcuts-322-12         175228        281475        +60.63%
BenchmarkShortestPathManyToMany/CH_shortest_path/128/vertices-128-edges-23977-shortcuts-1315-12     494103        1277792       +158.61%
BenchmarkShortestPathManyToMany/CH_shortest_path/256/vertices-256-edges-97227-shortcuts-5276-12     1268879       5410451       +326.40%

benchmark                                                                                           old allocs     new allocs     delta
BenchmarkShortestPathManyToMany/CH_shortest_path/4/vertices-4-edges-9-shortcuts-1-12                72             87             +20.83%
BenchmarkShortestPathManyToMany/CH_shortest_path/8/vertices-8-edges-61-shortcuts-1-12               121            120            -0.83%
BenchmarkShortestPathManyToMany/CH_shortest_path/16/vertices-16-edges-316-shortcuts-31-12           239            198            -17.15%
BenchmarkShortestPathManyToMany/CH_shortest_path/32/vertices-32-edges-1404-shortcuts-123-12         432            349            -19.21%
BenchmarkShortestPathManyToMany/CH_shortest_path/64/vertices-64-edges-5894-shortcuts-322-12         846            670            -20.80%
BenchmarkShortestPathManyToMany/CH_shortest_path/128/vertices-128-edges-23977-shortcuts-1315-12     1727           1372           -20.56%
BenchmarkShortestPathManyToMany/CH_shortest_path/256/vertices-256-edges-97227-shortcuts-5276-12     3633           2867           -21.08%

benchmark                                                                                           old bytes     new bytes     delta
BenchmarkShortestPathManyToMany/CH_shortest_path/4/vertices-4-edges-9-shortcuts-1-12                1588          3163          +99.18%
BenchmarkShortestPathManyToMany/CH_shortest_path/8/vertices-8-edges-61-shortcuts-1-12               3061          4242          +38.58%
BenchmarkShortestPathManyToMany/CH_shortest_path/16/vertices-16-edges-316-shortcuts-31-12           8876          8211          -7.49%
BenchmarkShortestPathManyToMany/CH_shortest_path/32/vertices-32-edges-1404-shortcuts-123-12         21948         17150         -21.86%
BenchmarkShortestPathManyToMany/CH_shortest_path/64/vertices-64-edges-5894-shortcuts-322-12         50805         35714         -29.70%
BenchmarkShortestPathManyToMany/CH_shortest_path/128/vertices-128-edges-23977-shortcuts-1315-12     112964        77437         -31.45%
BenchmarkShortestPathManyToMany/CH_shortest_path/256/vertices-256-edges-97227-shortcuts-5276-12     243966        164183        -32.70%
```