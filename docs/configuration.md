# Configuration

ZobTube configuration is pretty straightforward. It can either be passed through environmental variables, command line parameters or through a configuration file.

## Reference

Name | CLI parameter | Environmental parameter | YAML key | Default
--|---|--|--|--
Configuration file | `--config-file` | `$ZT_CONFIG_FILE` | — | `config.yml`
Gin debug mode | `--gin-debug` | `$ZT_GIN_DEBUG` | `log.gin.debug` | `false`
Log level | `--log-level` | `$ZT_LOG_LEVEL` | `log.level` | `1` (info)
Server bind address | `--server-bind` | `$ZT_SERVER_BIND` | `server.bind` | `0.0.0.0:8069`
Database driver | `--db-driver` | `$ZT_DB_DRIVER` | `db.driver` | `sqlite`
Database connection string | `--db-connstring` | `$ZT_DB_CONNSTRING` | `db.connstring` | `zobtube.sqlite`
Media path | `--media-path` | `$ZT_MEDIA_PATH` | `media.path` | `data`
Metadata storage type | `--metadata-type` | `$ZT_METADATA_TYPE` | `metadata.type` | `filesystem`
Metadata path (filesystem) | `--metadata-path` | `$ZT_METADATA_PATH` | `metadata.path` | same as `media-path`
Metadata S3 bucket | `--metadata-s3-bucket` | `$ZT_METADATA_S3_BUCKET` | `metadata.s3.bucket` | —
Metadata S3 region | `--metadata-s3-region` | `$ZT_METADATA_S3_REGION` | `metadata.s3.region` | `us-east-1`
Metadata S3 prefix | `--metadata-s3-prefix` | `$ZT_METADATA_S3_PREFIX` | `metadata.s3.prefix` | —
Metadata S3 endpoint | `--metadata-s3-endpoint` | `$ZT_METADATA_S3_ENDPOINT` | `metadata.s3.endpoint` | —
Metadata S3 access key ID | `--metadata-s3-access-key-id` | `$ZT_METADATA_S3_ACCESS_KEY_ID` | `metadata.s3.access_key_id` | —
Metadata S3 secret access key | `--metadata-s3-secret-access-key` | `$ZT_METADATA_S3_SECRET_ACCESS_KEY` | `metadata.s3.secret_access_key` | —

Log level values: `5` panic, `4` fatal, `3` error, `2` warn, `1` info, `0` debug, `-1` trace.

`db-driver` must be `sqlite` or `postgresql`. `metadata-type` must be `filesystem` or `s3`; when using S3, `metadata-s3-bucket` is required.

## Example with CLI parameters and local SQLite database

```sh
./zobtube server \
  --server-bind 127.0.0.1:8069 \
  --log-level 1 \
  --db-driver sqlite \
  --db-connstring ./zobtube.sqlite \
  --media-path ./data \
  --metadata-type filesystem \
  --metadata-path ./metadata
```

Only non-default values are required; the same setup with defaults omitted:

```sh
./zobtube server \
  --db-connstring ./zobtube.sqlite \
  --media-path ./data
```

Equivalent `config.yml` (use `--config-file` or place it as `config.yml` in the working directory):

```yaml
server:
  bind: "127.0.0.1:8069"
log:
  level: 1
db:
  driver: sqlite
  connstring: ./zobtube.sqlite
media:
  path: ./data
metadata:
  type: filesystem
  path: ./metadata
```

## Example with environmental variables and PostgreSQL database

```sh
export ZT_SERVER_BIND="0.0.0.0:8069"
export ZT_LOG_LEVEL=1
export ZT_DB_DRIVER=postgresql
export ZT_DB_CONNSTRING="host=localhost user=zobtube password=secret dbname=zobtube port=5432 sslmode=disable"
export ZT_MEDIA_PATH=/var/lib/zobtube/media
export ZT_METADATA_TYPE=filesystem
export ZT_METADATA_PATH=/var/lib/zobtube/metadata

./zobtube
```

Docker Compose fragment:

```yaml
environment:
  - ZT_SERVER_BIND=0.0.0.0:8069
  - ZT_DB_DRIVER=postgresql
  - ZT_DB_CONNSTRING=host=postgres user=zobtube password=secret dbname=zobtube port=5432 sslmode=disable
  - ZT_MEDIA_PATH=/var/lib/zobtube/media
  - ZT_METADATA_TYPE=filesystem
  - ZT_METADATA_PATH=/var/lib/zobtube/metadata
```

Equivalent `config.yml`:

```yaml
server:
  bind: "0.0.0.0:8069"
log:
  level: 1
db:
  driver: postgresql
  connstring: "host=postgres user=zobtube password=secret dbname=zobtube port=5432 sslmode=disable"
media:
  path: /var/lib/zobtube/media
metadata:
  type: filesystem
  path: /var/lib/zobtube/metadata
```