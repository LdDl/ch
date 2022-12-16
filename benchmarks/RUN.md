# This is just tech MD for comparing different development branches

## M-N search
```shell
git checkout redudant_pointers && \
go test -benchmem -run=^$ -bench ^BenchmarkShortestPathManyToMany$ github.com/LdDl/ch -v -count=1 > benchmarks/new_mn_ptr.txt && \
git checkout be2e2d4c1a059db481a6d0a8d250225ee3e68597 && \
go test -benchmem -run=^$ -bench ^BenchmarkShortestPathManyToMany$ github.com/LdDl/ch -v -count=1 > benchmarks/old_mn_ptr.txt && \
benchcmp benchmarks/old_mn_ptr.txt benchmarks/new_mn_ptr.txt && \
git checkout redudant_pointers
```

## 1-N search
```shell
git checkout redudant_pointers && \
go test -benchmem -run=^$ -bench ^BenchmarkShortestPathOneToMany$ github.com/LdDl/ch -v -count=1 > benchmarks/new_1n_ptr.txt && \
git checkout be2e2d4c1a059db481a6d0a8d250225ee3e68597 && \
go test -benchmem -run=^$ -bench ^BenchmarkShortestPathOneToMany$ github.com/LdDl/ch -v -count=1 > benchmarks/old_1n_ptr.txt && \
benchcmp benchmarks/old_1n_ptr.txt benchmarks/new_1n_ptr.txt && \
git checkout redudant_pointers
```

## 1-1 search
```shell
git checkout redudant_pointers && \
go test -benchmem -run=^$ -bench ^BenchmarkShortestPath$ github.com/LdDl/ch -v -count=1 > benchmarks/new_11_ptr.txt && \
git checkout be2e2d4c1a059db481a6d0a8d250225ee3e68597 && \
go test -benchmem -run=^$ -bench ^BenchmarkShortestPath$ github.com/LdDl/ch -v -count=1 > benchmarks/old_11_ptr.txt && \
benchcmp benchmarks/old_11_ptr.txt benchmarks/new_11_ptr.txt && \
git checkout redudant_pointers
```

## 1-1 search (single b.Run(...))
```shell
git checkout redudant_pointers && \
go test -benchmem -run=^$ -bench ^BenchmarkStaticCaseShortestPath$ github.com/LdDl/ch -v -count=1 > benchmarks/new_11static_ptr.txt && \
git checkout be2e2d4c1a059db481a6d0a8d250225ee3e68597 && \
go test -benchmem -run=^$ -bench ^BenchmarkStaticCaseShortestPath$ github.com/LdDl/ch -v -count=1 > benchmarks/old_11static_ptr.txt && \
benchcmp benchmarks/old_11static_ptr.txt benchmarks/new_11static_ptr.txt && \
git checkout redudant_pointers
```

## CH Prepare
```shell
git checkout redudant_pointers && \
go test -benchmem -run=^$ -bench ^BenchmarkPrepareContracts$ github.com/LdDl/ch -v -count=1 > benchmarks/new_ch_prepare_ptr.txt && \
git checkout be2e2d4c1a059db481a6d0a8d250225ee3e68597 && \
go test -benchmem -run=^$ -bench ^BenchmarkPrepareContracts$ github.com/LdDl/ch -v -count=1 > benchmarks/old_ch_prepare_ptr.txt && \
benchcmp benchmarks/old_ch_prepare_ptr.txt benchmarks/new_ch_prepare_ptr.txt && \
git checkout redudant_pointers
```


