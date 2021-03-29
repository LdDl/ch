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

### Planned
* Better heuristics for calculationg importance of each vertex.
* Max-cost path finder.
* N-best shortest pathes.