# ZobTube

_ZobTube, passion of the Zob, lube for the Tube!_


[![Go Report](https://goreportcard.com/badge/github.com/zobtube/zobtube)](https://goreportcard.com/report/github.com/zobtube/zobtube)
[![Go Coverage](https://github.com/zobtube/zobtube/wiki/coverage.svg)](https://raw.githack.com/wiki/zobtube/zobtube/coverage.html)

ZobTube is a library viewer and management tools for all your ~~porn~~ scientific videos.

![demo viewer](docs/readme_assets/demo_viewer.png)

To view more screenshots :camera:, see the [docs/screenshots.md](docs/screenshots.md) file.

## :information_source: Current status

:tada: Finally, all features, enhancements and bug fixed expected for 1.0.0 are released.
ZobTube will now enter a testing / QA phase to ensure acceptable quality for the 1.0.0 release.

All future developments will be followed on the [Kanban project view](https://github.com/orgs/zobtube/projects/1).

## :cop: About piracy

Piracy is bad. The FBI will not blow up your door for some downloaded porn. Yet, please respect the hard work of the actors by paying for their content.

> [!CAUTION]
> This goes without saying but **neither store illegal content nor content you don't own**.

## :wave: Everyone welcomed

As healthy sex life goes in pair of welcoming everyone, ZobTube endorses **LGBTQIA+** community.

For now, the only reference to some sexual identity is through the definition of actors. Only `male`, `female` and `trans women` are supported for now. But if anything's missing, feel free to create a pull request or a feature request.

## :vertical_traffic_light: Getting started

ZobTube works as a single binary. It needs a database to work but can rely on a local sqlite database. Parameters to start the binary can either be passed as a configration file or as environmental variables, as described below.

### Docker quickstart

The easiest way to start ZobTube is through its docker image, with the following command.

```
docker run -v ./zt-config:/config -v ./zt-data:/data -e ZT_DB_DRIVER=sqlite -e ZT_DB_CONNSTRING=/config/db.sqlite -e ZT_MEDIA_PATH=/data -e ZT_SERVER_BIND="0.0.0.0:8080" -p 8080:8080 ghcr.io/zobtube/zobtube
```

Then, ZobTube will be reachable on [http://127.0.0.1:8080](http://127.0.0.1:8080).

### Binary quickstart

Just start the binary without any parameter

```sh
./zobtube
```

If no configuration is provided, a default one will be created.
Then, ZobTube will be reachable on [http://127.0.0.1:8080](http://127.0.0.1:8080).

### Configuration file example

In `config.yml`, in the same folder as where ZobTube is started.

```yaml
server:
    bind: "127.0.0.1:8080"
db:
  driver: "sqlite"
  connstring: "zt.db"
media:
  path: "my_library_folder"
```

### Environmental variables example

```sh
ZT_SERVER_BIND="0.0.0.0:8080"
ZT_MEDIA_PATH="/mnt/zobtube"
ZT_DB_DRIVER="postgresql"
ZT_DB_CONNSTRING="host=pg user=zt password=topsecret dbname=zobtube port=5432 sslmode=disable"
```

### Configuration reference

Environmental variable name | Configuration variable name | Example values | Description
-|-|-|-
`ZT_SERVER_BIND` | `bind` | `127.0.0.1:8080` - `0.0.0.0:8080` | IP and port to lisen to.
`ZT_DB_DRIVER` | `db.driver` | `postgresql` - `sqlite` | Driver used for the database
`ZT_DB_CONNSTRING` | `db.connstring` | `zt.db` - `host=pg user=zt password=topsecret dbname=zobtube port=5432 sslmode=disable` | Connection string to pass to the database driver
`ZT_MEDIA_PATH` | `media.path` | `/mnt/zobtube` - `./my_library` - `C:\Users\zt\videos` | Library base path, where all content will be stored.

## :page_facing_up: License

ZobTube © 2025 by sblablaha is licensed under the MIT license.
