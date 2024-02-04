#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run ./internal/cmd/policygen completion "$sh" >"completions/policy-gen.$sh"
done
