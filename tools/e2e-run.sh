#!/bin/bash

set -e

# Resolve project root (parent of tools/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

echo '* clean up'
rm -rf /tmp/zt*

echo '* build binary'
go build -o /tmp/zt main.go

echo '* run tests (parallel via pytest-xdist)'
exec "$PROJECT_ROOT/.venv/bin/python" -m pytest "$@" -n auto --dist loadfile --junitxml e2e-result.xml
