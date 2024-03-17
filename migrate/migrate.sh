#!/bin/bash

psql -w -f migrate/create_db.sql
psql -w -f migrate/insert.sql