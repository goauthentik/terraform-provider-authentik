# This Provider is in alpha, the resources are not guaranteed to be stable!

<p align="center">
    <img src="https://goauthentik.io/img/icon_top_brand_colour.svg" height="150" alt="authentik logo">
</p>

---

[![](https://img.shields.io/discord/809154715984199690?label=Discord&style=for-the-badge)](https://discord.gg/jg33eMhnj6)
[![CI Build status](https://img.shields.io/github/checks-status/beryju/terraform-provider-authentik/master?style=for-the-badge)](https://github.com/BeryJu/terraform-provider-authentik/actions)
[![Code Coverage](https://img.shields.io/codecov/c/gh/beryju/terraform-provider-authentik?style=for-the-badge)](https://codecov.io/gh/BeryJu/terraform-provider-authentik)
![Latest version](https://img.shields.io/github/v/tag/beryju/terraform-provider-authentik?style=for-the-badge)


# Terraform Provider authentik

Run the following command to build the provider

```shell
go build -o terraform-provider-authentik
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

## Supported apps

- [ ] authentik.admin
- [ ] authentik.api
- [ ] authentik.core
- [x] authentik.crypto
- [ ] authentik.events
- [ ] authentik.flows
- [x] authentik.outposts
- [x] authentik.policies
- [x] authentik.policies.dummy
- [x] authentik.policies.event_matcher
- [x] authentik.policies.expiry
- [x] authentik.policies.expression
- [x] authentik.policies.hibp
- [x] authentik.policies.password
- [x] authentik.policies.reputation
- [ ] authentik.providers.ldap
- [x] authentik.providers.oauth2
- [x] authentik.providers.proxy
- [ ] authentik.providers.saml
- [ ] authentik.sources.ldap
- [ ] authentik.sources.oauth
- [ ] authentik.sources.plex
- [ ] authentik.sources.saml
- [x] authentik.stages.authenticator_duo
- [x] authentik.stages.authenticator_static
- [x] authentik.stages.authenticator_totp
- [x] authentik.stages.authenticator_validate
- [x] authentik.stages.authenticator_webauthn
- [x] authentik.stages.captcha
- [x] authentik.stages.consent
- [x] authentik.stages.deny
- [x] authentik.stages.dummy
- [x] authentik.stages.email
- [x] authentik.stages.identification
- [x] authentik.stages.invitation
- [x] authentik.stages.password
- [x] authentik.stages.prompt
- [x] authentik.stages.user_delete
- [x] authentik.stages.user_login
- [x] authentik.stages.user_logout
- [x] authentik.stages.user_write
- [x] authentik.tenants
