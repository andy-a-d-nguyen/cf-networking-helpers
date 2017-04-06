#!/bin/bash

set -eu
set -o pipefail

cd $(dirname $0)/..

function bootPostgres {
	echo -n "booting postgres"
	(/docker-entrypoint.sh postgres &> /var/log/postgres-boot.log) &
	trycount=0
	for i in `seq 1 30`; do
		set +e
		psql -h localhost -U postgres -c '\conninfo' &>/dev/null
		exitcode=$?
		set -e
		if [ $exitcode -eq 0 ]; then
			echo "connection established to postgres"
			return 0
		fi
		echo -n "."
		sleep 1
	done
	echo "unable to connect to postgres"
	exit 1
}

function bootMysql {
	echo -n "booting mysql"
	(MYSQL_ROOT_PASSWORD=password  /entrypoint.sh mysqld &> /var/log/mysql-boot.log) &
	trycount=0
	for i in `seq 1 30`; do
		set +e
		echo '\s;' | mysql -h 127.0.0.1 -u root --password="password" &>/dev/null
		exitcode=$?
		set -e
		if [ $exitcode -eq 0 ]; then
			echo "connection established to mysql"
			return 0
		fi
		echo -n "."
		sleep 1
	done
	echo "unable to connect to mysql"
	exit 1
}

go get github.com/onsi/ginkgo
go install github.com/onsi/ginkgo/ginkgo

# mutualtls test builds the binary
# TODO(gabe): test should just use certstrap library calls, rather than shelling out to the binary
# that way it is properly detected by go get -t ./...
go get github.com/square/certstrap


if [ "${1:-""}" = "" ]; then
  extraArgs=""
else
  extraArgs="${@}"
fi


if [ ${DB:-"none"} = "mysql" ]; then
  bootMysql
elif [ ${DB:-"none"} = "postgres" ]; then
  bootPostgres
else
  echo "skipping database"
  extraArgs="-skipPackage=db ${extraArgs}"
fi

go get -t ./...
ginkgo -r -p --race -randomizeAllSpecs -randomizeSuites ${extraArgs}
