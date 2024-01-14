package argonhasher

/*
	KDFconfig is the base struct for argonhasher, the Argon2id wrapper.
	It uses the standard library's argon2 IDKey function:
	func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
	$argon2id$v=19$t=10,m=65536,p=8$SALT$HASH
*/
type KDFconfig struct {

	/*
		Salt is the base64 string used to salt our derived keys.
	*/
	Salt []byte

	/*
		SaltLength
		length of random-generated salt
		min 16 bytes recommended for password hashing)
	*/
	SaltLength uint

	/*
		Time (i.e. iterations) - t
		number of iterations or pass throughs to perform
	*/
	Time uint32

	/*
		Memory - m
		amount of memory (in kilobytes) to use
	*/
	Memory uint32

	/*
		Threads (parallelism) p: degree of parallelism (i.e. number of threads)
	*/
	Threads uint8

	/*
		KeyLen T: desired number of returned bytes
		128 bit (16 bytes) sufficient for most applications
	*/
	KeyLen uint32
}
