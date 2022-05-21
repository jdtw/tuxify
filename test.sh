#!/bin/bash
set -euo pipefail

PORT=9090
ADDR="http://localhost:${PORT}"
TESTDIR="./testtmp"
# 3 just because I liked how that one looked...
KEY="00000000000000000000000000000003"

cleanup() {
    exit_status=$?
    killall -q -u "${USER}" tuxify-server || true
    rm -rf ${TESTDIR}
    exit "$exit_status"
}
trap cleanup EXIT

mkdir "${TESTDIR}"
go build -o "${TESTDIR}" ./...

"${TESTDIR}/tuxify-server" --port "${PORT}" &
until curl -s -X POST "${ADDR}" -o /dev/null; do
    echo "Waiting for server to start..."
    sleep 1
done

echo "Testing tuxify-server..."
result=$(curl -s -F 'img=@./testdata/tux.png' \
        -F "key=${KEY}" \
        -o "${TESTDIR}/out.png" \
        -w '%{http_code} %{content_type}' \
        "${ADDR}")
echo "${result}"
test "${result}" = "200 image/png"
cmp -s "${TESTDIR}/out.png" ./testdata/expected.png || (echo "Images differ!"; exit 1)

echo "Testing tuxify-server invalid key..."
result=$(curl -s -F 'img=@./testdata/tux.png' \
        -F 'key=thisisnothex' \
        -o /dev/null \
        -w '%{http_code}' \
        "${ADDR}")
echo "${result}"
test "${result}" = "400"
cmp -s "${TESTDIR}/out.png" ./testdata/expected.png || (echo "Images differ!"; exit 1)

echo "Testing tuxify..."
"${TESTDIR}/tuxify" --in ./testdata/tux.png --out ./testout.png --key "${KEY}"
cmp -s "${TESTDIR}/out.png" ./testdata/expected.png || (echo "Images differ!"; exit 1)

echo "Tests pass!"