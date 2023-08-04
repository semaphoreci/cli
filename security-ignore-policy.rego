package trivy

default ignore = false

ignore {
  input.PkgName == "gopkg.in/yaml.v2"
}
