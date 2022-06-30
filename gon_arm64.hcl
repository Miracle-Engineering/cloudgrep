# gon.hcl
#
# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source = ["./bin/darwin_arm64/cloudgrep"]
bundle_id = "dev.runx.cloudgrep"
dmg {
  output_path = "./dist/cloudgrep_darwin_arm64.dmg"
  volume_name = "cloudgrep"
}
apple_id {
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "74Y9V676W7"
}
