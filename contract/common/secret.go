package common

import "github.com/ZoneCNH/testkitx/requirex"

const SecretNoLeakID = "common.secret.no_secret_leak"

func RunSecretNoLeak(t requirex.TestingT, value any, secrets ...string) {
	t.Helper()
	requirex.NoSecretLeak(t, value, secrets...)
}
