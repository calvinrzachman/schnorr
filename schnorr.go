package schnorr

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

/*
	SCHNORR SIGNATURE METHODS

	Define the various methods for implementing Schnorr Signatures

*/

var (
	// Curve represents the secp256k1 elliptic curve
	Curve = btcec.S256() //elliptic.P256()

	// N represents the size of the finite cyclic group defined over the given elliptic curve
	N = Curve.Params().N

	// P represents the prime order of the finite field over which we take the elliptic curve...
	P = Curve.Params().P

	// One is a uint64 representation of the number 1
	One = new(big.Int).SetUint64(1)

	// Two is a uint64 representation of the number 1
	Two = new(big.Int).SetUint64(2)

	// Four is a uint64 representation of the number 4
	Four = new(big.Int).SetUint64(4)

	// Three is a uint64 representation of the number 3
	Three = new(big.Int).SetUint64(3)

	// Seven is a uint64 representation of the number 7
	Seven = new(big.Int).SetUint64(7)
)

// Sign creates a Schnorr digital signature (R, s) on the
// message digest using the given private key
func Sign(message [32]byte, privateKey *ecdsa.PrivateKey) ([64]byte, error) {
	signature := [64]byte{}
	publicKey := &privateKey.PublicKey

	if privateKey.D.Cmp(N) >= 0 || privateKey.D.Cmp(One) < 0 {
		return signature, errors.New("private key must be a positive integer no greater than the group order N-1")
	}

	// R = k*G
	var R Point
	D := intToBytes(privateKey.D)
	k0, err := deterministicK(D, message)
	if err != nil {
		return signature, err
	}
	R.X, R.Y = Curve.ScalarBaseMult(k0.Bytes())

	k := new(big.Int)
	if big.Jacobi(R.Y, Curve.P) != 1 {
		k = k.Sub(N, k0)
	} else {
		k = k0
	}

	// Calculate s = k + e*x
	e := calculateE(publicKey, R.X, message)
	e.Mul(e, privateKey.D)
	k.Add(k, e)
	k.Mod(k, N)

	// Construct signature (R, s)
	rX := intToBytes(R.X)
	s := intToBytes(k)
	copy(signature[:32], rX)
	copy(signature[32:], s)
	return signature, nil
}

// VerifySignature verifies that the signature (s, R)
// on the message is valid for the given public key
func VerifySignature(publicKey *ecdsa.PublicKey, message [32]byte, signature [64]byte) (bool, error) {

	// Obtain (R, s) from signature
	rX := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	if rX.Cmp(Curve.P) >= 0 || s.Cmp(N) >= 0 {
		return false, errors.New("parameters exceed field/group order")
	}

	e := calculateE(publicKey, rX, message)

	// Verify that s*G = R + e*P
	var sG, eP, R Point
	eP.X, eP.Y = Curve.ScalarMult(publicKey.X, publicKey.Y, intToBytes(new(big.Int).Sub(N, e)))
	sG.X, sG.Y = Curve.ScalarBaseMult(intToBytes(s))

	// Cannot construct R directly, as "k" unkown.
	// Solve for the point R using R = s*G + e*P
	R.X, R.Y = Curve.Add(sG.X, sG.Y, eP.X, eP.Y)

	// Note: To determine signature validity, we must verify that the x-coordinate of the public nonce, R,
	// that we calculate during verification matches the x-coordinate for R provided by the signature
	if R.X.Cmp(rX) != 0 || big.Jacobi(R.Y, Curve.P) != 1 {
		return false, errors.New("Invalid Signature")
	}

	return true, nil
}

// Point represents a point (X,Y) on the elliptic curve
type Point struct {
	X *big.Int
	Y *big.Int
}

// calculateE computes a commitment, e, to a particular public key,
// public nonce, and a message. The commitment is a SHA256
// hash operation of the form e = H(P, R, m)
func calculateE(publicKey *ecdsa.PublicKey, rX *big.Int, m [32]byte) *big.Int {
	secKey := SerializeCompressed(publicKey)
	data := append(intToBytes(rX), secKey[:]...)
	data = append(data, m[:]...)

	// Compute the commitment according to e = H(P, R, m)
	hash := sha256.Sum256(data)
	e := new(big.Int).SetBytes(hash[:])
	e.Mod(e, N)
	return e
}

// deterministicK accepts the integer representation of a private key
// and a 32 byte message digest and returns an integer k according to k = H(x, m)
func deterministicK(privateKey []byte, message [32]byte) (*big.Int, error) {
	data := append(privateKey, message[:]...)
	digest := sha256.Sum256(data)
	k := new(big.Int).SetBytes(digest[:])
	k.Mod(k, N)
	return k, nil
}

// CONSIDER REWORKING THIS
// intToBytes returns the big endian byte
// representation of a large integer
func intToBytes(bigInt *big.Int) []byte {
	var bytes [32]byte
	b := bigInt.Bytes()
	copy(bytes[32-len(b):], b)
	return bytes[:]
}

// SerializeCompressed converts a public key into compressed SEC format
func SerializeCompressed(publicKey *ecdsa.PublicKey) [33]byte {
	var compressedKey [33]byte
	b := make([]byte, 1, 33)
	// if y is even prefix with 0x2, if odd prefix with 0x3
	b[0] = 0x2
	b[0] += byte(publicKey.Y.Bit(0))
	b = append(b, publicKey.X.Bytes()...)
	copy(compressedKey[:], b)
	return compressedKey
}
