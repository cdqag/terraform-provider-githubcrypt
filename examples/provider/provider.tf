terraform {
  required_providers {
    githubcrypt = {
      source  = "cdqag/githubcrypt"
      version = "~> 1.0.0"
    }
  }
}

provider "githubcrypt" {
    owner = var.owner
    app_id = var.github_app_id
    app_installation_id = var.github_app_installation_id
    pem_file = var.github_app_pem_file
}
