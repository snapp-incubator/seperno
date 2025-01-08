#!/bin/bash

go build -o seperno.so -buildmode=c-shared export.go

mv seperno.so ./seperno/.
mv seperno.h ./seperno/.

rm -rf dist

python setup.py sdist bdist_wheel

#twine upload dist/*