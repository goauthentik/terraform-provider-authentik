resource "authentik_blueprint" "instance" {
  name = "blueprint-instance"
  path = "default/10-flow-default-authentication-flow.yaml"
  context = jsonencode(
    {
      foo = "bar"
    }
  )
}
