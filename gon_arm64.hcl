# gon.hcl
#
# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source = ["./bin/darwin_arm64/cloudgrep"]
bundle_id = "dev.runx.cloudgrep"
zip {
  output_path = "./bin/cloudgrep_darwin_arm64.zip"
}

apple_id {
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "74Y9V676W7"
}
