# Copyright (c) 2014 SameGoal LLC. All Rights Reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

lint:
	ls -1d *.go | xargs -I {} go fmt {}
	ls -1d *.go | xargs -I {} go vet {}
	! golint *.go | grep .

test:
	go test -v -x
