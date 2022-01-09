#!/bin/bash

# build the provider binary
go build -o terraform-provider-hashicups

# get system's OS and ARCH
export OS_ARCH="$(go env GOHOSTOS)_$(go env GOHOSTARCH)"

# ensure provider path within the user plugins directory exists
mkdir -p ~/.terraform.d/plugins/hashicorp.com/edu/hashicups/0.2/$OS_ARCH

# move binary to the appropriate subdirectory
mv terraform-provider-hashicups ~/.terraform.d/plugins/hashicorp.com/edu/hashicups/0.2/$OS_ARCH
