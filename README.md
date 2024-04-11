<p align="center">
    <img src="https://goauthentik.io/img/icon_top_brand_colour.svg" height="150" alt="authentik logo">
</p>

---

[![](https://img.shields.io/discord/809154715984199690?label=Discord&style=for-the-badge)](https://discord.gg/jg33eMhnj6)
[![Code Coverage](https://img.shields.io/codecov/c/gh/goauthentik/terraform-provider-authentik?style=for-the-badge)](https://codecov.io/gh/goauthentik/terraform-provider-authentik)
[![Latest version](https://img.shields.io/github/v/tag/goauthentik/terraform-provider-authentik?style=for-the-badge)](https://registry.terraform.io/providers/goauthentik/authentik/latest)
[![CI Build status](https://img.shields.io/github/actions/workflow/status/goauthentik/terraform-provider-authentik/test.yml?branch=main&style=for-the-badge)](https://github.com/goauthentik/terraform-provider-authentik/actions)

# Terraform Provider authentik

Tested against authentik main and stable, on terraform 1.2.1

Run the following command to build the provider

```shell
make build
```

### Generate Documentation

Run `make` from the project root to regenerate the latest provider documentation

## Run tests using a Dev Container

Running the included Dev Container will create a full authentik development environment automatically as if following the instructions found here: https://goauthentik.io/docs/installation/docker-compose

Once the Dev Container is running, simply use the VS Code Command Palette or Test UI to run any tests as needed.

Note: If running all tests, this is very CPU-intensive on your local authentik environment, so depending on your hardware they can take several minutes to complete.

## Run tests using a local environment

Start a local authentik instance by following https://goauthentik.io/docs/installation/docker-compose

Before starting the instance, add this static token to the `.env` file:

```
AUTHENTIK_TOKEN=this-token-is-for-testing-dont-use
```

Afterwards, tests can be run from VS Code with the Command Palette or Test UI, or via CLI like so:

```
export TF_ACC=1
export AUTHENTIK_URL=http://localhost:9000
export AUTHENTIK_TOKEN=this-token-is-for-testing-dont-use
go test -timeout 30m ./... -count=1
```

If you're trying to run tests with VS Code in your local environment, be sure to change `AUTHENTIK_URL` in `.vscode/settings.json` to: `"AUTHENTIK_URL": "http://localhost:9000"`

Note: If running all tests, this is very CPU-intensive on your local authentik environment, so depending on your hardware they can take several minutes to complete.

## Versioning

This provider's version is based on the authentik version it's tested against.

Provider version 2021.8.1 is tested against 2021.8.x.

Provider version 2021.8.2 is tested against 2021.8.x but has some bugfixes.
