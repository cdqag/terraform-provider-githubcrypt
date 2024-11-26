data "githubcrypt_environment_public_key" "my_environment_public_key" {
  repo_id = var.repo_id
  environment = github_repository_environment.my_environment.name
}
