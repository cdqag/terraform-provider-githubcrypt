---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "githubcrypt_environment_public_key Data Source - terraform-provider-githubcrypt"
subcategory: ""
description: |-
  GitHub Repository Environment Public Key Data Source
---

# githubcrypt_environment_public_key (Data Source)

GitHub Repository Environment Public Key Data Source



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment` (String) The name of the GitHub repository environment.
- `repo_id` (Number) The ID of the GitHub repository.

### Read-Only

- `public_key` (String) The public key of the GitHub repository environment.