resource "authentik_blueprint" "instance" {
  name = "blueprint-instance"
  path = "default/flow-default-authentication-flow.yaml"
  context = jsonencode(
    {
      foo = "bar"
    }
  )
}
