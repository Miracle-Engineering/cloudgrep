# gon.hcl
#
# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source = ["./bin/cloudgrep_darwin_amd64"]
bundle_id = "dev.runx.cloudgrep"
zip {
  output_path = "./bin/cloudgrep_darwin_amd64.zip"
}

apple_id {
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "74Y9V676W7"
}
