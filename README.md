# ZobTube

_ZobTube, passion of the Zob, lube for the Tube!_


[![Go Report](https://goreportcard.com/badge/github.com/zobtube/zobtube)](https://goreportcard.com/report/github.com/zobtube/zobtube)
[![Go Coverage](https://github.com/zobtube/zobtube/wiki/coverage.svg)](https://raw.githack.com/wiki/zobtube/zobtube/coverage.html)

ZobTube is a library viewer and management tools for all your ~~porn~~ scientific videos.

## :information_source: Current status

> [!WARNING]
> ZobTube is under active development! Everything ***should*** work properly.
> Yet, as long as `1.0.0` is not released, there will be no promises.

**TL;DR**: Not stable yet but should work fine.

The remaining work towards the incoming stable release is [available below](#coming-developments).

## :cop: About piracy

Piracy is bad. The FBI will not blow up your door for some downloaded porn. Yet, please respect the hard work of the actors by paying for their content.

> [!CAUTION]
> This goes without saying but **neither store illegal content nor content you don't own**.

## :wave: Everyone welcomed

As healthy sex life goes in pair of welcoming everyone, ZobTube endorses **LGBTQIA+** community.

For now, the only reference to some sexual identity is through the definition of actors. Only `male`, `female` and `shemale` are supported for now. But if anything's missing, feel free to create a pull request or a feature request.

## :vertical_traffic_light: Getting started

ZobTube works as a single binary. It needs a database to work but can rely on a local sqlite database. Parameters to start the binary can either be passed as a configration file or as environmental variables, as described below.

### Quickstart

Just start the binary without any parameter

```sh
./zobtube
```

### Configuration file example

In `config.yml`, in the same folder as where ZobTube is started.

```yaml
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

## :white_check_mark: Coming developments

### Road to 1.0.0

**Features**

- [x] Manage external commands through async tasks
  - [x] Compute duration
  - [x] Generate thumbnail
  - [x] Generate mini thumbnail
  - [x] Move files
  - [x] View those tasks on the administration page
- [x] List all videos / actors / channels through an admin page
- [x] Implement channels for videos
- [ ] Add categories
- [x] Add file/folder deletion in triage
- [ ] Add folder creation in triage
- [ ] Add 'set as image' in triage towards videos / actors / channels
- [ ] Clips view
- [x] Edit actor aliases

**Bugs**

- [x] Fix bug with large uploaded files in triage
- [x] Create own favicon
- [x] Add linting validation
- [x] Add tests
  - [x] Validate routes requiring authentication
  - [x] Validate routes requiring admin rights
  - [x] Validate upload
  - [x] Validate video viewing
- [x] Manually add links of actors
- [x] Delete actors
- [x] Add actors on video page after selections
- [x] Fix mini thumb down-sizing scale issue
- [x] Fix suggestion videos

**Enhancement**

- [x] Add a welcome page when no configuration is specified to bootstrap the app
- [x] Write onboarding readme
- [ ] Add actor description
- [ ] Rename actors

### Enhancements not expected before 1.1.0

- [ ] Add pagination
- [ ] Movie scenes
- [ ] Collections
- [ ] Merge actors
- [ ] Redirect to requested page after login

## :page_facing_up: License

ZobTube © 2024 by sblablaha is licensed under CC BY-SA 4.0. To view a copy of this license, visit [CC website](https://creativecommons.org/licenses/by-sa/4.0/)
