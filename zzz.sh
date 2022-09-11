#!/bin/sh
curPath=`pwd`


export GIT_COMMIT=$(git rev-parse HEAD)
export BUILD_TIME=$(date -u '+%Y-%m-%d %I:%M:%S %Z')

# zzz run -ldflags "-w -s -X \"github.com/phper95/mail-server/internal/conf.BuildCommit=$(git rev-parse HEAD)\" -X \"github.com/phper95/mail-server/internal/conf.BuildTime=$(date -u '+%Y-%m-%d %I:%M:%S %Z')\""

echo zzz run -ldflags "-X \"github.com/phper95/mail-server/internal/conf.BuildTime=${BUILD_TIME}\" -X \"github.com/phper95/mail-server/internal/conf.BuildCommit=${GIT_COMMIT}\""
zzz run --ld "-X \"github.com/phper95/mail-server/internal/conf.BuildTime=${BUILD_TIME}\" -X \"github.com/phper95/mail-server/internal/conf.BuildCommit=${GIT_COMMIT}\""
