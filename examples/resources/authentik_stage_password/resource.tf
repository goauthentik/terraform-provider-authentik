# Create a password stage that tests against the interla database

resource "authentik_stage_password" "test" {
  name     = "test-stage"
  backends = ["authentik.core.auth.InbuiltBackend"]
}
