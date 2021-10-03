## ROADMAP
New ideas, thought about needed features will be store in this file.

### Done
* Initial core
    * Import file of specific format **Done as CSV**
    * Export file of specific format **Done as CSV**
    * Turn Restricted Shortest Path extension for CH-algorithm **Propably not modify algorithm, but graph. Excluded: done with osm2ch - https://github.com/LdDl/osm2ch#osm2ch**
    * Thoughts and discussions about OSM graph and extensions **Need some ideas about parsing and preparing. Excluded: osm2ch - https://github.com/LdDl/osm2ch#osm2ch**
    * Map matcher as another project **Here is the one https://github.com/LdDl/horizon**
    * Bring OSM parser **Excluded**
    * Bring OSM restrictions **Excluded**
    * OneTwoMany function (contraction hierarchies) **Done, may be some bench comparisons**
    * Replace int with int64 (OSM purposes) **Done**
    * Separate benchmarks to BENCHMARK.md **Done**
    * Better CSV format or another format (JSON / binary). **W.I.P. Splitting single file to multiple**
    * Separate export functions

### WIP
* Parallel version as optional feature (See branch [optional-parallelism](https://github.com/LdDl/ch/tree/)). Status update: 22.08.2021
* Refactor code completely. May be add some simple variation of algorithm (or add docs/wiki, the best I've found so far: https://jlazarsfeld.github.io/ch.150.project/sections/1-intro/)

### Planned
* Better heuristics for calculationg importance of each vertex.
* Max-cost path finder.
* N-best shortest pathes.
