package ssh

import "golang.org/x/crypto/ssh"

// signingKeyTypes is the set of SSH public key algorithms that can be
// used to sign a git commit.
var signingKeyTypes = map[string]struct{}{
	ssh.KeyAlgoRSA:        {},
	ssh.KeyAlgoECDSA256:   {},
	ssh.KeyAlgoECDSA384:   {},
	ssh.KeyAlgoECDSA521:   {},
	ssh.KeyAlgoED25519:    {},
	ssh.KeyAlgoSKECDSA256: {},
	ssh.KeyAlgoSKED25519:  {},
}

// FilterSigningKeys returns the subset of keys that can be used to
// sign a git commit. Keys that fail to parse are discarded.
func FilterSigningKeys(keys []string) []string {
	valid := []string{}

	for _, key := range keys {
		pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key))
		if err != nil {
			continue
		}
		if _, ok := signingKeyTypes[pk.Type()]; ok {
			valid = append(valid, key)
		}
	}

	return valid
}
