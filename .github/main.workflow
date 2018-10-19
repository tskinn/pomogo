workflow "main" {
  on = "push"
  resolves = "test"
}

action "test" {
  uses = "golang:alpine"
  runs = "go test ./..."
}
