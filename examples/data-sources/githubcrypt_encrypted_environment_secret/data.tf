data "githubcrypt_encrypted_environment_secret" "reference-data-fetch-and-transform-job_prod_argocd_user_secret" {
  public_key_base64 = data.githubcrypt_environment_public_key.my_environment_public_key.public_key
  secret = var.my_secret
}
