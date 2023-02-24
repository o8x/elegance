#!/bin/bash

set -ex

export GOOS=linux GOARCH=amd64 CGO_ENABLED=0

if [[ -f /.dockerenv ]]; then
    fpm -t rpm -v 1.0.0 --verbose -f -s dir -n elegance \
        --config-files /var/elegance/elegance.sqlite \
        -a amd64 \
        -p /app \
        -C /app/build

    fpm -t deb -v 1.0.0 --verbose -f -s dir -n elegance \
        --config-files /var/elegance/elegance.sqlite \
        -a amd64 \
        -p /app \
        -C /app/build
    exit 0
fi

[[ -d build ]] && rm -rf build
mkdir -p build/sbin build/var/elegance build/etc/systemd/system
cp elegance.sqlite build/var/elegance
cp elegance.service build/etc/systemd/system
cp ~/go/pkg/mod/github.com/huichen/sego@*/data/dictionary.txt build/var/elegance/dictionary.txt

cd view || exit
npm run build
mv build ../build/var/elegance/views
cd ..

go env -w "GOPROXY=https://goproxy.cn,direct"
go build -ldflags="
    -X main.StaticRoot=/var/elegance/views
    -X github.com/o8x/elegance/app/database.DataFile=/var/elegance/elegance.sqlite
    -X github.com/o8x/elegance/app/word.DictFile=/var/elegance/dictionary.txt
" -v -o build/sbin/elegance .

docker run --platform linux/amd64 -i -v $(pwd):/app -w /app unsafe/porter /app/build.sh

rm -rf build
