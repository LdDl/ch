# osm2ch
## Convert *.osm.pbf files to CSV

- [About](#about)
- [Installation](#installation)
- [Usage](#usage)
- [Example](#example)
- [Dependencies](#dependencies)
- [License](#license)

## About
With this CLI tool you can convert *.osm.pbf (Compressed Open Street Map) file to CSV (Comma-Separated Values) file, which is used in our [contraction hierarchies library].
What it does:
- Edge expansion (single edge == single vertex);
- Handles some kind and types of restrictions:
    - Supported kind of restrictions:
        - EdgeFrom - NodeVia - EdgeTo.
    - Supported types of restrictions:
        - only_left_turn;
        - only_right_turn;
        - only_straight_on;
        - no_left_turn;
        - no_right_turn;
        - no_straight_on.
- Saves CSV file with geom in WKT format;
- Currently supports tags for 'highway' OSM entity only.

PRs are welcome!

## Installation
* Via 'go get':
    ```shell
    go install github.com/LdDl/ch/...
    ```
    After installation step is complete you can call 'osm2ch' from any place in your system.

* Or download prebuilt binary and make updates in yours PATH environment varibale (both Linux and Windows):
    * Windows - https://github.com/LdDl/ch/releases/download/v1.3.4/windows-osm2ch.zip
    * Linux - https://github.com/LdDl/ch/releases/download/v1.3.4/linux-osm2ch.tar.xz
## Usage
```shell
osm2ch -h
```
Output:
```shell
Usage of osm2ch:
  -file string
        Filename of *.osm.pbf file (it has to be compressed) (default "my_graph.osm.pbf")
  -geomf string
        Format of output geometry. Expected values: wkt / geojson (default "wkt")
  -out string
        Filename of 'Comma-Separated Values' (CSV) formatted file (default "my_graph.csv")
  -tags string
        Set of needed tags (separated by commas) (default "motorway,primary,primary_link,road,secondary,secondary_link,residential,tertiary,tertiary_link,unclassified,trunk,trunk_link")
  -units string
        Units of output weights. Expected values: km for kilometers / m for meters (default "km")
```
The default list of tags is this, since usually these tags are used for routing for personal cars.


## Example
You can find example file of *.osm.pbf file in nested child [/example_data](/example_data).

If you want WKT format for output geometry:
```shell
osm2ch --file example_data/moscow_center_reduced.osm.pbf --out graph.csv --geomf wkt --units m --tags motorway,primary,primary_link,road,secondary,secondary_link,residential,tertiary,tertiary_link,unclassified,trunk,trunk_link
```
If you want GeoJSON format for output geometry:
```shell
osm2ch --file example_data/moscow_center_reduced.osm.pbf --out graph.csv --geomf geojson --units m --tags motorway,primary,primary_link,road,secondary,secondary_link,residential,tertiary,tertiary_link,unclassified,trunk,trunk_link
```

After that file 'graph.csv' will be created.
Header of CSV-file is: from_vertex_id;to_vertex_id;weights;geom
- from_vertex_id Source vertex;
- to_vertex_id Target vertex;
- weight Traveling cost from source to target (actually length of edge in kilometers);
- geom Geometry of edge in WKT or GeoJSON format.

Now you can use this graph in [contraction hierarchies library].


## Dependencies
Thanks to [paulmach](https://github.com/paulmach) for his [OSM-parser](https://github.com/paulmach/osm) written in Go.

Paulmach's license is [here](https://github.com/paulmach/osm/blob/master/LICENSE.md) (it's MIT)

## License
- Please see [here](https://github.com/LdDl/ch/blob/master/LICENSE)


[contraction hierarchies library]: (https://github.com/LdDl/ch#ch---contraction-hierarchies)