package token

func jwtioSecret() func(KeyID string) (string, error) {
	return func(key string) (string, error) {
		return "secretcode", nil
	}
}

const jwtioToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tL21hcmtzdGFuZGVuIiwiYXVkIjoiZ2l0aHViLmNvbS9tYXJrc3RhbmRlbi9hdXRoZW50aWNhdGlvbiIsInN1YiI6IjEyMzQ1Njc4OTAiLCJqdGkiOiJuVjFNMkIyWmwtU0MwNEdhWkpwN3FEcVA0M0duQzFaZ3R0VDBFOGR2aC1qc2VQRjBsNXAwRUVrS01IOHdJejVNMnpsenI1R0wzUi1UODltSy1OUndBUT09Iiwia2lkIjoiTVdJeFNZbl9RZFgybVBGRml3ZnUyTHVzT2lYaWRNUGpEX2lzMEtyNEJLdnZzYmdBQUUyM0xuVmRqSThVQUZXMUZ6LTlMSlBPcUs5TEFueldwWHBRcHc9PSIsImlhdCI6MTYwMDAwMDAwMCwibmJmIjoxNjAwMDAwMDAwLCJleHAiOjE2NTAwMDAwMDB9.HFX2e6yxMOgeife_-EAKr3Mgv03pFXB7TWb5M6aTwCpcT_oz3zBb67e-jIwKmd141JwAxuxanYG4eaErh014NA"

/*  The JSON payload output from the jwt.io website for out token, with our secret
const jwtioJSON = `
{
  "iss": "github.com/markstanden",
  "aud": "github.com/markstanden/authentication",
  "sub": "1234567890",
  "jti": "nV1M2B2Zl-SC04GaZJp7qDqP43GnC1ZgttT0E8dvh-jsePF0l5p0EEkKMH8wIz5M2zlzr5GL3R-T89mK-NRwAQ==",
  "kid": "MWIxSYn_QdX2mPFFiwfu2LusOiXidMPjD_is0Kr4BKvvsbgAAE23LnVdjI8UAFW1Fz-9LJPOqK9LAnzWpXpQpw==",
  "iat": 1600000000,
  "nbf": 1600000000,
  "exp": 1650000000
}
`*/

var jwtioStruct = Payload{
	Issuer:   "github.com/markstanden",
	Audience: "github.com/markstanden/authentication",

	UserID: "1234567890",
	JwtID:  "nV1M2B2Zl-SC04GaZJp7qDqP43GnC1ZgttT0E8dvh-jsePF0l5p0EEkKMH8wIz5M2zlzr5GL3R-T89mK-NRwAQ==",

	KeyID: "MWIxSYn_QdX2mPFFiwfu2LusOiXidMPjD_is0Kr4BKvvsbgAAE23LnVdjI8UAFW1Fz-9LJPOqK9LAnzWpXpQpw==",

	IssuedAtTime:   1600000000,
	NotBeforeTime:  1600000000,
	ExpirationTime: 1650000000,
}
