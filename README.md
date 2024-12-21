# MbTiles

Create a MbTiles file from OACI map from Geoportail

# Usage

```
manage MbTiles from OACI

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

# Development

Open the project as a devContainer with VSCode. Please refer to:
* VSCode install: https://code.visualstudio.com/download
* DevContainers: https://code.visualstudio.com/docs/devcontainers/containers

# Build

```bash
go build -o mbtiles ./cmd
```