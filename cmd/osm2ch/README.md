# osm2ch
Convert *.osm.pbf files to CSV

- [About](#about)
- [Installation](#installation)
- [Usage](#usage)
- [Example](#example)
- [Dependencies](#dependencies)
- [License](#license)

## About
With this CLI tool you can convert *.osm.pbf (Compressed Open Street Map) file to CSV (Comma-Separated Values) file, which is used in [our contraction hierarchies library](https://github.com/LdDl/ch#ch---contraction-hierarchies)

## Installation
```shell
go install github.com/LdDl/ch/...
```
After installation step is complete you can call 'osm2ch' from any place in your system.

## Usage
```shell
osm2ch -h
```
Output:
```sh
Usage of osm2ch:
  -file string
        Filename of *.osm.pbf file (it has to be compressed) (default "my_graph.osm.pbf")
  -out string
        Filename of 'Comma-Separated Values' (CSV) formatted file (default "my_graph.csv")
  -tags string
        Set of needed tags (separated by commas) (default "motorway,primary,primary_link,road,secondary,secondary_link,residential,tertiary,tertiary_link,unclassified,trunk,trunk_link")
```

## Example
You can find example file of *.osm.pbf file in nested child [/example_data](/example_data)

## Dependencies
Thanks to [paulmach](https://github.com/paulmach) for his [OSM-parser](https://github.com/paulmach/osm) written in Go.

Paulmach's license is [here](https://github.com/paulmach/osm/blob/master/LICENSE.md) (it's MIT)

## License
- Please see [here](https://github.com/LdDl/ch/blob/master/LICENSE)