#!/usr/bin/env bash

export APP_ENV="dev"

SERVICES_ROOT="services"

echo $SERVICES_ROOT/$1_service

cd $SERVICES_ROOT/$1_service
go run .

