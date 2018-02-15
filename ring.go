package ring

import (
	"crypto/elliptic"
	crand "crypto/rand"
	"fmt"
	"io"
)

// Generate generates a new public-private key pair.
// If no random generator is provided, Generate will use
// go's default cryptographic random generator.
// The private key should be safely stored.
// The public key can be shared with anyone.
func Generate(rand io.Reader) (PublicKey, PrivateKey) {
	if rand == nil {
		rand = crand.Reader
	}

	curve := elliptic.P384()
	sk, x, y, err := elliptic.GenerateKey(curve, rand)
	if err != nil {
		panic(fmt.Sprintf("Could not generate keys: %s", err.Error()))
	}

	pk := elliptic.Marshal(curve, x, y)

	return PublicKey(pk), PrivateKey(sk)
}

// Signing algorithm (Schnorr Ring Signature):
//	* Let (P(0),...,P(R-1)) be all the public keys in the ring
//	* P(i)=x(i)*G (x(i) is the private key)
//	* Let H be the chosen hash function (probably SHA256)
//	* Let n be the number of bits of the prime number defining the curve (probably 256).
//	* Let r be the index of the actual signer in the ring
//	* Randomly choose k in {0,1}^n
//	* Compute e(r+1 % R) = H(m || k*G)
//	* for i := r+1 % R; i != r; i++:
//		* Randomly choose s(i) in {0,1}^n
//		* Compute e(i+1 % R) = H(m || s(i)*G + e(i)*P(i))
//	* Compute s(r) = k - e(r)*x(r)
//	* Output signature: (P(0),...,P(1),e(0),s(0),...,s(r))

// Verifying algorithm:
//	* Let (P(0),...P(R-1),e,s(0),...,s(r)) be the input signature of message m
//	* Let ee = e
//	* for i := 0; i < R; i++:
//		* ee = H(m || s(i)*G + ee*P(i))
//	* If ee = e, signature is valid. Otherwise it's invalid.
