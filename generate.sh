#!/bin/bash

go build -o seperno.so -buildmode=c-shared export.go

python setup.py sdist bdist_wheel

twine upload dist/*