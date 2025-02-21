# polkadot_exporter
Prometheus exporter for polkadot nodes.

### Usage:
```
## Local usage
$ make build/expoter
$ bin/polkadot-exporter-local-compatible
```
```
## Using as a docker image
$ make build/exporter/dockerImage DOCKER_IMAGENAME=<reponame>

# Image tag will be created based on the git commit, tag or version
$ docker image ls
$ docker run -d -p 9100:9100 --name exporter <imageid>:<tag>


# additional options of the commandline 
docker run -d -p 9100:9100 --name exporter <imageid>:<tag> "--log-level=debug --listen-addr=http://0.0.0.0:9100"

```

### Makefile
```
$ make help 
Usage:
  help                         print the help message
  build/exporter               building the exporter
  build/exporter/dockerImage   building the docker image of the exporter
  run/exporter                 run the exporter
  audit                        executing quality check, linting and unit tests
  vendor                       vendoring the packages in case necessary
```

### Commandline
```
// Golang version go1.23.6

This exporter is used to get connected to the polkadot binary api for collection
        of data and converting it to prometheus metrics

Usage:
  polkadot_exporter [flags]

Flags:
      --RpcBackend string    rpc backend to expose metrics from (default "ws://localhost:9944")
      --crt string           HTTPs certificate .pem file path
      --crt-key string       HTTPs key .pem file path
  -h, --help                 help for polkadot_exporter
      --listen-addr string   listen address for the exporter (default "http://0.0.0.0:9100")
      --log-level string     loglevel of the exporter. possible values are debug, info, warn, error, fatal, panic, and trace (default "info")
  -t, --toggle               Help message for toggle
      --version              show the version and build time of the exporter

```

### Github actions and workflows
To avoid failure on github actions consider add the below secrets and environment variables in github "secrets and variables section"
<pre>
secrets:
DOCKER_REGISTRY_USERNAME
DOCKER_REGISTRY_PASS

envs:
DOCKER_IMAGE_NAME
</pre>
