package argonhasher

import (
	"testing"
)

/***
*
*	Tests
*
***/

/*
	results is key value store to store hashes and the respective passwords used to create them.  We will fill this up in the first test, and then compare them in the second test
*/
var results = make(map[string]string)

func TestEncodeThenConfirm(t *testing.T) {
	/*
		test is the base type for our individual tests
	*/
	type test struct {
		/*
			desc is the written reason for conducting the test
		*/
		desc string

		/*
			pw is the message that will be hashed
			by the Encode function
		*/
		pw string

		/*
			cost is difficulty in performing the hash
		*/
		cost uint

		/*
			shouldFail is true for tests that we would expect to fail
			i.e. not hash correctly
			shouldFail is false for tests that we would expect to result in a valid hash
		*/
		shouldFail bool
	}

	/*
		tests is a slice of multiple tests, poor passwords and then rediculously long complex passwords, at different costs
	*/
	tests := []test{
		{desc: "invalid, zero length password", pw: "", cost: 0, shouldFail: true},
		{desc: "basic short password", pw: "password", cost: 0, shouldFail: false},
		{desc: "repeated basic short password", pw: "password", cost: 0, shouldFail: false},
		{desc: "repeated basic short password", pw: "password", cost: 0, shouldFail: false},
		{desc: "very complex, long password, cost 1", pw: "SxbLPNws8kT!if4P69wV9PDy@KcMHPRjsMU%7L*f&8jA&LRaNJVyz!zAU*gWrr9tfiYKxKcw#wM43^q7QrtHQyoj@%!EEtNLZPt#saxHrZu5#tGX#UgE5mwHNy$b!vm", cost: 1, shouldFail: false},
		{desc: "very complex, long password, cost 2", pw: "nJqQ7^aDVxUrjTnDxjFNZzF&CT!5BajSycDm8u#v8bxvx$gBtr7iL7WVpnkqLDqGqgRYdQN!guXZDLXvsW6DKhCce25UWbGpZnWUqdGQ88^KmtLSvySjz7FRXWyU9kC", cost: 2, shouldFail: false},
		{desc: "very complex, long password, cost 3", pw: "*SsXYqgVX8x7KBnHTgHkPbFSMFUa&9os$j2C!fRmU2r%BdpP6RXs6k^m2HnyyCQ3NehEZ!gRtVkGX@RhAE&$*ZbNV*!gmwMuVLm!WA^7ouSnJeFAAHGZs@KxAWW8H9S", cost: 3, shouldFail: false},
		{desc: "very complex, long password, cost 4", pw: "Pdid&4b2gjQA!K$qstUNtgv!8B!FsgHFe2hLFh9co#r7%FgS3jz$RZYUnz6r99ta#nkqC@Q#xPUj46xGufprKT@$SMFkud6SWdR7ecyurxPtXvfDqbDUXTeLGzEkb3H", cost: 4, shouldFail: false},
		{desc: "very complex, long password, cost 5", pw: "ywgDTeXuzijxaZY4zspnSRe5M^gCsJaNxosZrrRfayc7mBaU44!kPrJz6&*dG6oh6a!uGG5@3D##qCi&Yxk$gBrC@3hk^HzveXrW4NPjMUsvRa9umDG$cqSsrPWQ&Gx", cost: 5, shouldFail: false},
		{desc: "very complex, long password, cost 6", pw: "iq@Ni@iCPGeSQavNSbw3FxQ89^9SW^KCW5tCaEBLKSWdxKPXR3HrmzRU%&p&wyo8@cP!ZAqcdzV#JeSr4!fdoKbDjh7vdr6r*hHHvLU8GjQKgb&YwpuRfpp*8E^igAh", cost: 6, shouldFail: false},
		{desc: "very complex, long password, cost 7", pw: "hwq&65MEtQ86JJ$&JxWRUvA*wErV3NfgT!qKJ%Xgg6hfQ8PwgT%xt9P@uK4#hGub#JmQ@p666*r3PQU7p$ssAC2!7*g6@5dudbDCrC$pwXTsn7@o8T3GoFhX3ix*6U%", cost: 7, shouldFail: false},
		{desc: "very complex, long password, cost 8", pw: "KSky7@NBde3mYukPXt&Su5&w&8Uy7GaRib5e@57qADc^Rncm3TM$Rh*YWi9ZACm9@Exn8J^cUg&#Bt&p3riR$gN#cGxci725f@c!!aw%DGyofPQZS5kSF9*bDBS9Dk@", cost: 8, shouldFail: false},
		{desc: "very complex, long password, cost 9", pw: "PB#E&EPug%JAW@SAs62u3n4CJns$K#B*R7Q8ogkdszeFMmNz2UzvVe@ZQKxL%%B&m!%k34F#cNQdijewhYYDm6vd&GMoxaxqnXBtdM*4fiwCH*Pkds6yFMqa4hm^Joo", cost: 9, shouldFail: false},
		{desc: "very complex, long password, cost 10", pw: "#3$Vj26hkk2ihuNm7YhG@2HabxGdBHQ&S*XLpgMCsVd!SnMqw2PJ&eJ@jB#EFfznTD4@Jfk6u^&xzQj&HV@szhs7rTYxpJEVqoWn^RuaJ6Yh6k&bq2Ki$BfXojMR9#o", cost: 10, shouldFail: false},
	}

	for _, test := range tests {

		/*
			scope it at the test level so it can
			be used in the comparison tests
		*/
		var hash string

		t.Run("Hashing :"+test.desc, func(t *testing.T) {
			/*
				Hash the password and check whether the hash is already present in the map.  If it is, our encoder is producing identical keys, and
				something has gone wrong
			*/
			hash = Encode(test.pw, test.cost)
			if _, ok := results[hash]; ok {
				t.Errorf("Multiple identical hashes produced: \n%v", hash)
			}

			/*
				Empty passwords should return empty strings,
				as they will be more easily brute forced
			*/
			if !ValidHash(hash) && !test.shouldFail {
				t.Error("Invalid hash produced")
			}

			/*
				Empty passwords should return empty strings,
				as they will be more easily brute forced
			*/
			if hash != "" && test.shouldFail {
				t.Errorf("Failed test has produced a valid hash.  \nGot: \n%v\nWanted:\n\"\"", hash)
			}

			/*
				Assign the password to the validated hash so we can
				decrypt it in the next stage
			*/
			if !test.shouldFail {
				results[hash] = test.pw
			}
		})

		/*
			Use the Compare function to check the plaintext password
			against the created hash.
		*/
		t.Run("Comparing: "+test.desc, func(t *testing.T) {
			valid := Confirm(test.pw, hash)
			if !valid && !test.shouldFail {
				t.Errorf("failed to decode hash \n%v\nwith the password\n%v", hash, test.pw)
			}
		})
	}
}
