# Create a captcha stage

resource "authentik_stage_captcha" "name" {
  name        = "captcha"
  private_key = "foo"
  public_key  = "bar"
}
