#!/bin/bash -e

cd /tmp
curl -LO https://raw.githubusercontent.com/bitnami/bitnami-docker-harbor-portal/master/docker-compose.yml
curl -L https://github.com/bitnami/bitnami-docker-harbor-portal/archive/master.tar.gz -o master.tar.gz
tar xzf master.tar.gz
mv bitnami-docker-harbor-portal-master/2/debian-10/config/ .
# vim config/proxy/nginx.conf
# add
# server > location
#   proxy_set_header Authorization $http_authorization;
#   proxy_pass_header  Authorization;
#
# add to /etc/hosts
# 127.0.0.1       localhost reg.mydomain.com
docker-compose up -d

echo "login to http://localhost using admin and bitnami"

for t in $(seq 1 5); do
    docker tag utrecht/n3dr:6."${t}".0 localhost/library/n3dr:6."${t}".0
    docker push localhost/library/n3dr:6."${t}".0
done
