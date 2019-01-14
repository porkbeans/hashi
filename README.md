# Hashi

Hashi is a tool for downloading HashiCorp tools.

[![Build Status](https://travis-ci.org/porkbeans/hashi.svg?branch=master)](https://travis-ci.org/porkbeans/hashi)
[![Coverage Status](https://coveralls.io/repos/github/porkbeans/hashi/badge.svg?branch=master)](https://coveralls.io/github/porkbeans/hashi?branch=master)

# Installation

`go get -u github.com/porkbeans/hashi/cmd/hashi`

# Examples

```bash
# Show help
hashi help

# List all products
hashi list

# List versions of consul
hashi list consul

# List zip files of vault 1.0.1
hashi list vault 1.0.1

# Install packer 1.3.3 for your environment
hashi install packer 1.3.3 /usr/local/bin/packer

# Install terraform 0.11.11 for darwin amd64
hashi install terraform 0.11.11 /usr/local/bin/terraform --os darwin --arch amd64
```
