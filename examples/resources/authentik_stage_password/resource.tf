# Create a password stage that tests against the interla database

resource "authentik_stage_password" "test" {
  name     = "test-stage"
  backends = ["django.contrib.auth.backends.ModelBackend"]
}
