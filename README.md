# AllowedSignatures

This is a little script that allows you to use GitHub's API to download the public keys of GitHub collaborators and save them to a file for use with Git's SSH signing features.

## Installation

You can download the binary from the [Releases](https://github.com/frankywahl/allowedSignersFile/releases) or use it directly with go, as described below

## Example Usage

```bash
go run *.go --owner frankywahl --repository allowedSignersFile > .git/allowedSignersFile
git config gpg.ssh.allowedSignersFile .git/allowedSignersFile
```

Note: we can use the `--use-contributors` as a means to get all the contributors to a repo. However, this is much more expensive on GitHub requests.

## Requirements

* [GITHUB_API_TOKEN](https://github.com/settings/tokens) To use the Github API
* [go](https://go.dev/) (if you want to run if from source)

## Limitations

There is an assumption that users do not have more that 100 SSH keys attached to their profile.

## SSH Signing

```
# .git/config - can also be global configuration
[user]
        signingKey = $(cat ~/.ssh/id_ed25519.pub) # the output of the public key
[gpg]
        format = ssh
[gpg "ssh"]
        allowedSignersFile = .git/allowedSignatures
[commit]
        gpgsign = true
[tag]
        gpgsign = true
```

https://calebhearth.com/sign-git-with-ssh

# SSH Signing Github Support

Github has supported SSH Signing since [August 2022](https://github.blog/changelog/2022-08-23-ssh-commit-verification-now-supported/)
That being said SSH Commit signing was part of [Git](https://git-scm.com/) [beforehand](https://lore.kernel.org/git/xmqq8rxpgwki.fsf@gitster.g/)

If commits were signed before the release on Github, they will still appear as verified provided the SSH public key was still uploaded as a [Signing Key](https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key#telling-git-about-your-ssh-key)

https://github.com/community/community/discussions/7744#discussioncomment-2564268
