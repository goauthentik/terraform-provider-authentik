<p align="center">
    <img src="https://goauthentik.io/img/icon_top_brand_colour.svg" height="150" alt="authentik logo">
</p>

---

[![](https://img.shields.io/discord/809154715984199690?label=Discord&style=for-the-badge)](https://discord.gg/jg33eMhnj6)
[![Code Coverage](https://img.shields.io/codecov/c/gh/goauthentik/terraform-provider-authentik?style=for-the-badge)](https://codecov.io/gh/goauthentik/terraform-provider-authentik)
![Latest version](https://img.shields.io/github/v/tag/goauthentik/terraform-provider-authentik?style=for-the-badge)

|                   | Status        |
| ----------------- | ------------- |
| authentik master  | ![CI Build status](https://img.shields.io/github/workflow/status/goauthentik/terraform-provider-authentik/test-acc-authentik-master?style=for-the-badge) |
| authentik stable | ![CI Build status](https://img.shields.io/github/workflow/status/goauthentik/terraform-provider-authentik/test-acc-authentik-stable?style=for-the-badge) |

# Terraform Provider authentik

Tested against authentik master and stable, on terraform 0.15 and 1.0

Run the following command to build the provider

```shell
make build
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

## Versioning

This providers version is based ont he authentik version it's tested against.

Provider version 2021.8.1 is tested against 2021.8.x.

Provider version 2021.8.2 is tested against 2021.8.x but has some bugfixes.
