package argonhasher

import (
	"testing"
)

func TestValidHash(t *testing.T) {
	tests := []struct {
		desc  string
		hash  string
		valid bool
	}{
		{
			desc:  "Empty Hash",
			hash:  "",
			valid: false,
		},
		{
			desc:  "Missing Version",
			hash:  "$argon2id$v=$t=2,m=65537,p=2$RfA+5nFy7ASo+3pk0A2X3ANgYCTt/LvT15n2m9Ctj76ok0+AO9xUKhev4YGAb2c6ne48DKGFaErxzOTbjn2qcg$ReE+hbfju1q1AsV/YEimBA",
			valid: false,
		},
		{
			desc:  "Missing Argument, t",
			hash:  "$argon2id$v=19$t=,m=65537,p=2$RfA+5nFy7ASo+3pk0A2X3ANgYCTt/LvT15n2m9Ctj76ok0+AO9xUKhev4YGAb2c6ne48DKGFaErxzOTbjn2qcg$ReE+hbfju1q1AsV/YEimBA",
			valid: false,
		},
		{
			desc:  "Missing Argument, m",
			hash:  "$argon2id$v=19$t=2,m=,p=2$p3+NKdkMbCvwUZsekr5OHD/MdYEuwIoOLt2Fw1luKkqlAkwBamwbvN1ffRhw/l1LXnFb/KpvlXvnWsUXdiRIxQ$/mNt5y3Ee8eOmxf4vWHVMQ",
			valid: false,
		},
		{
			desc:  "Missing Argument, p",
			hash:  "$argon2id$v=19$t=2,m=65537,p=$g51GtGoeQJhHq5qr2UpQ1rRfhmz0BIMo9Lqtmo6QHmLurzm6NqqGuA43V90PrNmx7RGtDhMRmLAZjronQ9eqow$/eKGyISmy0t4f1O2dy8Z7A",
			valid: false,
		},
		{
			desc:  "Missing Salt",
			hash:  "$argon2id$v=19$t=4,m=131073,p=4$$UkQv7pi3iOzuyvJ0ji3IMQ",
			valid: false,
		},
		{
			desc:  "Missing Key",
			hash:  "$argon2id$v=19$t=6,m=196609,p=5$tQbqPI/StC+CZqtPacUD4J3uLu72oVYVCkFr2OBuCIKG9y8s/MTArL9WvWkPNJUICLkbhgRgl3StYSG4rRs2jQ$",
			valid: false,
		},
		{
			desc:  "Low Cost",
			hash:  "$argon2id$v=19$t=2,m=65537,p=2$6wstRR9j1/f+xDwjOyRbnvXYdIa1OQJdHQx6wqNaW/8NpUJ+ouE7B1DGNqs+oVZtcLHcQmbe8mypqswzu4xXrA$NSWN6OP1VrRaNTeEkih68A",
			valid: true,
		},
		{
			desc:  "High Cost",
			hash:  "$argon2id$v=19$t=20,m=655361,p=16$JCPbFW5sAW1wFtBaD7WfDFnU/opWa4fiuBfygndnlO3m0ydLU8lE6M7RK1fvcl6SjniTC9iEpNc+AmqP7Z62FQ$i0bWevgcv9tb2DWCV8FKJw",
			valid: true,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			if ValidHash(test.hash) != test.valid {
				t.Error("invalid hash")
			}
		})
	}
}
