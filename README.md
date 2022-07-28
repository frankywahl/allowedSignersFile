# AllowedSignatures

This is a little script that allows you to use GitHub's API to download the public keys of GitHub collaborators and save them to a file for use with Git's SSH signing features.

## Example Usage

```bash
go run *.go --owner frankywahl --repository allowedSignersFile > .git/allowedSignersFile
git config gpg.ssh.allowedSignersFile .git/allowedSignersFile
```

Note: we can use the `--use-contributors` as a means to get all the contributors to a repo. However, this is much more expensive on GitHub requests.

## Requirements

* [go](https://go.dev/)

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

https://github.com/community/community/discussions/7744#discussioncomment-2564268
