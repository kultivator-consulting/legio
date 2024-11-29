#!/bin/bash

migrate -path database/migration/ -database "postgres://root:secret@localhost:5432/cortexdb?sslmode=disable" down