package token

const jwtioSecret = "secretcode"

const jwtioToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tL21hcmtzdGFuZGVuIiwic3ViIjoiMTIzNDU2Nzg5MCIsImF1ZCI6ImdpdGh1Yi5jb20vbWFya3N0YW5kZW4vYXV0aGVudGljYXRpb24iLCJleHAiOjE2NTAwMDAwMDAsIm5iZiI6MTYwMDAwMDAwMCwiaWF0IjoxNjAwMDAwMDAwLCJqdGkiOiJuVjFNMkIyWmwtU0MwNEdhWkpwN3FEcVA0M0duQzFaZ3R0VDBFOGR2aC1qc2VQRjBsNXAwRUVrS01IOHdJejVNMnpsenI1R0wzUi1UODltSy1OUndBUT09Iiwia2lkIjoiTVdJeFNZbl9RZFgybVBGRml3ZnUyTHVzT2lYaWRNUGpEX2lzMEtyNEJLdnZzYmdBQUUyM0xuVmRqSThVQUZXMUZ6LTlMSlBPcUs5TEFueldwWHBRcHc9PSJ9.tbQ5tU9f6TdKPwiftAAwbgst1fpqzT1kBQ2TU2d7ADt9AE632AhXVqSnAxFzET2wt6Nz47MJERCjvPVj_Pe2uQ"

const jwtioJSON = `
{	
	"iss":"github.com/markstanden",
	"sub":"1234567890",
	"aud":"github.com/markstanden/authentication",
	"exp":1650000000,
	"nbf":1600000000,
	"iat":1600000000,
	"jti":"nV1M2B2Zl-SC04GaZJp7qDqP43GnC1ZgttT0E8dvh-jsePF0l5p0EEkKMH8wIz5M2zlzr5GL3R-T89mK-NRwAQ==",
	"kid":"MWIxSYn_QdX2mPFFiwfu2LusOiXidMPjD_is0Kr4BKvvsbgAAE23LnVdjI8UAFW1Fz-9LJPOqK9LAnzWpXpQpw=="
}`

var jwtioStruct = Payload{
	Issuer:         "github.com/markstanden",
	Subject:        "1234567890",
	Audience:       "github.com/markstanden/authentication",
	ExpirationTime: 1650000000,
	NotBeforeTime:  1600000000,
	IssuedAtTime:   1600000000,
	TokenID:        "nV1M2B2Zl-SC04GaZJp7qDqP43GnC1ZgttT0E8dvh-jsePF0l5p0EEkKMH8wIz5M2zlzr5GL3R-T89mK-NRwAQ==",
	KeyID:          "MWIxSYn_QdX2mPFFiwfu2LusOiXidMPjD_is0Kr4BKvvsbgAAE23LnVdjI8UAFW1Fz-9LJPOqK9LAnzWpXpQpw==",
}
