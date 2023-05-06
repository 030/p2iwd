# Nexus3

[Start a Nexus3 server, create a docker repository and upload images](https://github.com/030/n3dr/blob/main/docs/quickstarts/DOCKER.md).

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
docker_ref="repository/docker-images/utrecht/n3dr"
for t in {1..4}; do
  tag="6.${t}.0"
  docker load -i ${docker_ref}/${tag}/image.tar
  docker run localhost:8082/${docker_ref}:${tag}
done
```

Push the images:

```bash
./p2iwd push \
  --host http://localhost:9001 \
  -u admin \
  -p $(docker exec -it nexus3-n3dr-dest cat /nexus-data/admin.password) \
  --dir $PWD
```

Cleanup:

```bash
docker stop nexus3-n3dr-dest nexus3-n3dr-src
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
