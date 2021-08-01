# This Provider is in alpha, the resources are not guaranteed to be stable!

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
- [ ] authentik.lib
- [ ] authentik.managed
- [x] authentik.outposts
- [ ] authentik.policies
- [ ] authentik.policies.dummy
- [ ] authentik.policies.event_matcher
- [ ] authentik.policies.expiry
- [ ] authentik.policies.expression
- [ ] authentik.policies.hibp
- [ ] authentik.policies.password
- [ ] authentik.policies.reputation
- [ ] authentik.providers.ldap
- [ ] authentik.providers.oauth2
- [ ] authentik.providers.proxy
- [ ] authentik.providers.saml
- [ ] authentik.recovery
- [ ] authentik.sources.ldap
- [ ] authentik.sources.oauth
- [ ] authentik.sources.plex
- [ ] authentik.sources.saml
- [ ] authentik.stages.authenticator_duo
- [ ] authentik.stages.authenticator_static
- [ ] authentik.stages.authenticator_totp
- [ ] authentik.stages.authenticator_validate
- [ ] authentik.stages.authenticator_webauthn
- [ ] authentik.stages.captcha
- [ ] authentik.stages.consent
- [ ] authentik.stages.deny
- [ ] authentik.stages.dummy
- [ ] authentik.stages.email
- [ ] authentik.stages.identification
- [ ] authentik.stages.invitation
- [x] authentik.stages.password
- [ ] authentik.stages.prompt
- [x] authentik.stages.user_delete
- [x] authentik.stages.user_login
- [x] authentik.stages.user_logout
- [x] authentik.stages.user_write
- [x] authentik.tenants
