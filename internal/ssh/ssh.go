package ssh

import "strings"

// FilterSigingKeys will let us know if the ssh key
// used can be used to sign a git commit.
func FilterSigningKeys(keys []string) []string {
	valid := []string{}

	for _, key := range keys {
		if validKey(key) {
			valid = append(valid, key)
		}
	}

	return valid
}

var validPrefixes = [...]string{
	"ssh-rsa",
	"ecdsa-sha2-nistp256",
	"ecdsa-sha2-nistp384",
	"ecdsa-sha2-nistp521",
	"ssh-ed25519",
	"sk-ecdsa-sha2-nistp256@openssh.com",
	"sk-ssh-ed25519@openssh.com",
}

func validKey(key string) bool {
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}

	return false
}
