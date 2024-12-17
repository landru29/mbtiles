# MbTiles

Create a MbTiles file from OACI map from Geoportail

# Usage

```manage MbTiles from OACI

Usage:
  mbtiles [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    generate MbTiles from OACI
  help        Help about any command
  metadata    manage metadata on MbTiles
  tile        manage tiles on MbTiles

Flags:
  -d, --database string   database filename (default "oaci.mbtiles")
  -h, --help              help for mbtiles

Use "mbtiles [command] --help" for more information about a command.
```

# Build

```bash
go build -o mbtiles ./cmd
```