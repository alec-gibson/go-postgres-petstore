#!/bin/bash

case $1 in
	up)
		migrate -source file://sql/migrations -database 'postgres://postgres:example@localhost:5432/postgres?sslmode=disable' up
		;;

	down)
		migrate -source file://sql/migrations -database 'postgres://postgres:example@localhost:5432/postgres?sslmode=disable' down
		;;

	create)
		if [ -z "$2" ]; then
			echo "Usage: ./migrate.sh create <migration_name>"
			exit 1
		fi
		migrate create -dir sql/migrations -ext sql $2
		;;

	*)
		echo "Valid commands: up, down, create"
		;;
esac
