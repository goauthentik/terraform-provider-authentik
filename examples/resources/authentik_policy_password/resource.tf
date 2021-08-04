# Create a password policy to require 8 chars

resource "authentik_policy_password" "name" {
  name          = "password"
  length_min    = 8
  error_message = "foo"
}
