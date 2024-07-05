#!/bin/bash
set -o nounset -o errexit -o pipefail

sqlite3def --file=./test.sql ./test.db

