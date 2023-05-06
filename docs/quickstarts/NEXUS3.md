# Nexus3

[Start a Nexus3 server](https://github.com/030/n3dr/blob/359-maven-quickstart/docs/quickstarts/snippets/nexus3/SERVER.md).

Create a docker repository once Nexus3 has been started after a couple of
minutes:

```bash
n3dr configRepository \
  -u admin \
  -p $(docker exec -it nexus3-n3dr-src cat /nexus-data/admin.password) \
  -n localhost:8081 \
  --https=false \
  --configRepoName some-name \
  --configRepoType docker
```

Push several docker images:

```bash
docker login localhost:8082 \
  -u admin \
  -p $(docker exec -it nexus3-n3dr-src cat /nexus-data/admin.password) && \
  for t in {0..2}; do
    docker pull utrecht/n3dr:6.8.${t} && \
    docker tag utrecht/n3dr:6.8.${t} \
    localhost:8082/repository/some-name/utrecht/n3dr:6.8.${t} && \
    docker push localhost:8082/repository/some-name/utrecht/n3dr:6.8.${t}
  done
```

Pull the images:

```bash
./p2iwd pull \
  --host http://localhost:8082 \
  -u admin \
  -p $(docker exec -it nexus3-n3dr-src cat /nexus-data/admin.password) \
  --dir $PWD
```

Run the images:

```bash
for t in {0..2}; do
  docker load -i repository/some-name/utrecht/n3dr/6.8.${t}/image.tar
  docker run localhost:8082/repository/some-name/utrecht/n3dr:6.8.${t} --version
done
```

Push the images:

```bash
./p2iwd push \
  --host http://localhost:8082 \
  -u admin \
  -p $(docker exec -it nexus3-n3dr-src cat /nexus-data/admin.password) \
  --dir $PWD
```

Cleanup:

```bash
docker stop nexus3-n3dr-src
```

Note:

- `p2iwd pull` and `p2iwd push` also work without arguments. Create a
  `~/.p2iwd/config.yml` file with the following content:

```bash
---
dir: some-dir
host: http://localhost:9001
logLevel: trace
pass: some-pass
syslog: false
user: admin
```

and try this as well.

- pull an individual image by specifying the `--repo` and `--tag` parameters:

```bash
p2iwd pull --repo repository/some-name/utrecht/n3dr --tag 6.8.2
docker load -i repository/some-name/utrecht/n3dr/6.8.2/image.tar
docker run localhost:9001/repository/some-name/utrecht/n3dr:6.8.2 --version
```
