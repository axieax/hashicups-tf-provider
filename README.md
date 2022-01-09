# Terraform Provider HashiCups

Run the following command to build the provider

```shell
go build -o terraform-provider-hashicups
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```

References:

- [Custom Terraform Provider Tutorial](https://learn.hashicorp.com/collections/terraform/providers)
- [Official HashiCups Provider](https://github.com/hashicorp/terraform-provider-hashicups)
