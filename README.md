# git-credential-op

A `git-credential` helper for 1Password.

Only username and password authentication (Basic Auth) is supported at this time (no [`capability[]`](https://git-scm.com/docs/git-credential#Documentation/git-credential.txt-codecapabilitycode) codes supported). 

## Installation

```
go install github.com/brettbuddin/git-credential-op@latest
```

## Setup

```
# .gitconfig or .git/config

[credential]
    helper = op
```

### Customize

#### Account and Vault

```
# .gitconfig or .git/config

[credential "https://github.com"]
    helper = "op --account personalaccount.1password.com --vault Private"

[credential "https://githubenterprise.companyname.com"]
    helper = "op --account companyaccount.1password.com --vault Private"
```

We've included `--vault` above to illustrate it can be set, but the tool will use whatever Vault 1Password considers the
default for the account; usually "Private".

#### Locator Tag

Every 1Password item managed by `git-credential-op` is tagged with a locator tag. By default this is
`git-credential-op`, but you can change it if you don't care for it. Once you change it in your configuration, you'll
need to make sure any items in 1Password with the old tag are updated so the helper can find them.

```
# .gitconfig or .git/config

[credential]
    helper = "op --locator-tag my-cool-tag"
```

#### Title Template

Every 1Password item managed by `git-credential-op` is named in accordance to a template that you can customize.
Changing this won't affect the helper's ability to locate the item. 

```
# .gitconfig or .git/config

[credential]
    helper = "op --title 'git: {{.Host}}'"
```

#### Additional Tags

You can add additional tags to items managed by the helper. All created items will have these tags included alongside
the locator tag.

```
# .gitconfig or .git/config

[credential]
    helper = "op --additional-tags one,two,three"
```
