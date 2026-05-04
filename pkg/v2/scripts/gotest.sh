#!/usr/bin/env bash

echo "==> Running go test and creating a coverage profile..."
i=0
failed=0
for testingpkg in $(go list ./.../testing .); do
  coverpkg=${testingpkg%/testing}
  go test -v -covermode count -coverprofile "./${i}.coverprofile" -coverpkg $coverpkg $testingpkg
  if [ $? -eq 1 ]; then
     failed+=1
  fi
  ((i++))
done
gocovmerge $(ls ./*.coverprofile) | grep -v '\.gen\.go' | grep -v 'mock_' > coverage.out
rm *.coverprofile

if ((failed>0)); then
	exit 1
fi

exit 0
