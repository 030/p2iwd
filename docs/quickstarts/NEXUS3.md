# Nexus3

Start a Nexus3 server:

```bash
docker run --rm -d -p 9000:8081 -p 9001:8082 --name nexus3-p2iwd \
sonatype/nexus3:3.42.0
```

Create a docker repository once Nexus3 has been started after a couple of
minutes:

```bash
n3dr configRepository -u admin \
-p $(docker exec -it nexus3-p2iwd cat /nexus-data/admin.password) \
-n localhost:9000 --https=false --configRepoName some-name \
--configRepoType docker
```

Push several docker images:

```bash
docker login localhost:9001 \
-u admin \
-p $(docker exec -it nexus3-p2iwd cat /nexus-data/admin.password) && \
for t in {0..2}; do
docker pull utrecht/n3dr:6.8.${t} && \
docker tag utrecht/n3dr:6.8.${t} \
localhost:9001/repository/some-name/utrecht/n3dr:6.8.${t} && \
docker push localhost:9001/repository/some-name/utrecht/n3dr:6.8.${t}
done
```

Pull the images:

```bash
p2iwd pull --host http://localhost:9001 -u admin \
-p $(docker exec -it nexus3-p2iwd cat /nexus-data/admin.password) \
--dir $PWD
```

Run the images:

```bash
for t in {0..2}; do
docker load -i repository/some-name/utrecht/n3dr/6.8.${t}/image.tar
docker run localhost:9001/repository/some-name/utrecht/n3dr:6.8.${t} --version
done
```

Push the images:

```bash
p2iwd push --host http://localhost:9001 -u admin \
-p $(docker exec -it nexus3-p2iwd cat /nexus-data/admin.password) \
--dir $PWD
```

Cleanup:

```bash
docker stop nexus3-p2iwd
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
