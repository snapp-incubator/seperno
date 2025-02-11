#!/bin/bash

go build -o seperno/seperno.so -buildmode=c-shared export.go

rm -rf dist

python setup.py sdist bdist_wheel

#twine upload dist/*