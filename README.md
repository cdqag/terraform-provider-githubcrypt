# terraform-provider-githubcrypt

According to the GitHub API documentation, to set a secret in GitHub (on any level: organization, repository, codespace, repo's environment, etc.), you need to encrypt the secret locally, and send to GitHub the encrypted value (see [details](https://docs.github.com/en/rest/guides/encrypting-secrets-for-the-rest-api?apiVersion=2022-11-28)).

The currently available (as of 2024-08) [terraform GitHub provider](https://github.com/integrations/terraform-provider-github) does not support automatic encryption of secrets sent to GitHub by the provider if you want the secret to rest in terraform state as encrypted.

While the terraform GitHub provider can encrypt the secret for you, if you pass it to the provider as plain text, it will store the plan text value in the terraform state, which is not recommended.

On the other hand, the provider does not offer any logic to encrypt the secret locally, and pass it as the encrypted value - the crypt logic is left to the user.

What is more, crypting the secret requires the user to have the `public_key` generated by GitHub for each organization, repository, codespace, environment, etc. But unfotunatelly, while the provider provides data source to fetch the `public_key` from GitHub on all the levels, the repository environments level is not curently supported. And yes - each repository environment has its own `public_key`...

Having this in mind, the `terraform-provider-githubcrypt` was created to provide the missing functionality:

* to fetch the `public_key` from GitHub for the repository environments.
* to encrypt the secret locally, and pass the encrypted value to the official GitHub provider.

This means that the `terraform-provider-githubcrypt` is rather not a standalone provider, but a helper provider, which should be used together with the official GitHub provider.
