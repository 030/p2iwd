#!/bin/bash -e

GITHUB_TAG="${GITHUB_TAG:-local}"
SHA512_CMD="${SHA512_CMD:-sha512sum}"
export p2iwd_DELIVERABLE="${p2iwd_DELIVERABLE:-p2iwd}"

echo "GITHUB_TAG: '$GITHUB_TAG' p2iwd_DELIVERABLE: '$p2iwd_DELIVERABLE'"
cd cmd/p2iwd
go build -ldflags "-X main.Version=${GITHUB_TAG}" -o "${p2iwd_DELIVERABLE}"
$SHA512_CMD "${p2iwd_DELIVERABLE}" >"${p2iwd_DELIVERABLE}.sha512.txt"
chmod +x "${p2iwd_DELIVERABLE}"
cd ../..
