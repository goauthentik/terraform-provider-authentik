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
- [x] authentik.outposts
- [ ] authentik.policies
- [ ] authentik.policies.dummy
- [ ] authentik.policies.event_matcher
- [ ] authentik.policies.expiry
- [ ] authentik.policies.expression
- [ ] authentik.policies.hibp
- [x] authentik.policies.password
- [x] authentik.policies.reputation
- [ ] authentik.providers.ldap
- [ ] authentik.providers.oauth2
- [ ] authentik.providers.proxy
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
