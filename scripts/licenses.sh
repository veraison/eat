#!/bin/bash

set -e

type go-licenses &> /dev/null || go get github.com/google/go-licenses

MODULES+=("github.com/veraison/eat")

for module in ${MODULES[@]}
do
  echo ">> retrieving licenses [ ${module} ]"
  go-licenses csv ${module}
done
