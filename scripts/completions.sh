#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run ./cmd/policy-gen completion "$sh" >"completions/policy-gen.$sh"
done
