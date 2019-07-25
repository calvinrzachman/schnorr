package schnorr_test

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/calvinrzachman/schnorr"
)

/*
	Schnorr Signature Unit Tests

	Unit tests verifying expected functionality of schnorr.go

*/

// TestSign() verifies the correct function of Sign()
func TestSign(t *testing.T) {
	Curve := btcec.S256()
	// Curve := elliptic.P256()
	// N := Curve.Params().N
	testCases := []struct {
		name              string
		publicKey         string
		privateKey        string
		messageDigest     string
		expectedSignature string
	}{
		{ // TEST CASE #1: Test that schnorr signatures are produced for 32-byte messages
			name:              "Produces a signature for a 32-byte message",
			publicKey:         "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
			privateKey:        "0000000000000000000000000000000000000000000000000000000000000001",
			messageDigest:     "0000000000000000000000000000000000000000000000000000000000000000",
			expectedSignature: "787a848e71043d280c50470e8e1532b2dd5d20ee912a45dbdd2bd1dfbf187ef67031a98831859dc34dffeedda86831842ccd0079e1f92af177f7f22cc1dced05",
		},
		{ // TEST CASE #2: Test that schnorr signatures are produced for 32-byte messages
			name:              "Produces a signature for a 32-byte message",
			publicKey:         "02dff1d77f2a671c5f36183726db2341be58feae1da2deced843240f7b502ba659",
			privateKey:        "b7e151628aed2a6abf7158809cf4f3c762e7160f38b4da56a784d9045190cfef",
			messageDigest:     "243f6a8885a308d313198a2e03707344a4093822299f31d0082efa98ec4e6c89",
			expectedSignature: "2a298dacae57395a15d0795ddbfd1dcb564da82b0f269bc70a74f8220429ba1d1e51a22ccec35599b8f266912281f8365ffc2d035a230434a1a64dc59f7013fd",
		},
		{ // TEST CASE #3: Test that schnorr signatures are produced for 32-byte messages
			name:              "Produces a signature for a 32-byte message",
			publicKey:         "03fac2114c2fbb091527eb7c64ecb11f8021cb45e8e7809d3c0938e4b8c0e5f84b",
			privateKey:        "c90fdaa22168c234c4c6628b80dc1cd129024e088a67cc74020bbea63b14e5c7",
			messageDigest:     "5e2d58d8b3bcdf1abadec7829054f90dda9805aab56c77333024b9d0a508b75c",
			expectedSignature: "00da9b08172a9b6f0466a2defd817f2d7ab437e0d253cb5395a963866b3574be00880371d01766935b92d2ab4cd5c8a2a5837ec57fed7660773a05f0de142380",
		},
		// Populate with the remainder of the test cases
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var messageDigest [32]byte
			// var secKey [33]byte

			// Decode Test Parameters
			private, _ := hex.DecodeString(testCase.privateKey)
			// publicBytes, _ := hex.DecodeString(testCase.publicKey)
			msgBytes, _ := hex.DecodeString(testCase.messageDigest)
			copy(messageDigest[:], msgBytes)
			// fmt.Printf("Message Digest: %x\n", messageDigest)
			// fmt.Println("Public: ", publicBytes, len(publicBytes))

			privateKey, _ := btcec.PrivKeyFromBytes(Curve, private)
			// compressedKey := publicKey.SerializeCompressed()
			// fmt.Printf("PrivateKey: %x\n", privateKey.D)
			// fmt.Println("PublicKey: ", publicKey.SerializeCompressed())
			// fmt.Printf("SEC Compressed Key: %x\n", compressedKey)

			// Produce a Schnorr signature on the message
			signature, err := schnorr.Sign(messageDigest, privateKey.ToECDSA())

			// t.Log("Private Key: ", privateKey.ToECDSA().D, len(intToBytes(privateKey.ToECDSA().D)), N, len(intToBytes(N)))
			// t.Log("Hex to String: ", hex.EncodeToString(signature[:]))
			// t.Log("Expected Sign: ", testCase.expectedSignature)

			if err != nil || hex.EncodeToString(signature[:]) != testCase.expectedSignature {
				t.Fatalf("Unable to correctly sign message. %s", err)
			} else {
				t.Logf("SUCCESS: Produced Schnorr signature %x for message", signature)
			}

			// Verify we generate the expected signature
			// copy(secKey[:], publicBytes)
			// result, err := schnorr.VerifySignature(secKey, messageDigest, signature)
			// if err != nil { // || err != nil
			// 	t.Fatalf("Did not confirm/deny validity of signature as expected - Error: %s", err)
			// } else {
			// 	t.Logf("SUCCESS: Expected verify result. Schnorr Signature is valid: %t    Error: %s", result, err)
			// }
		})
	}
}

func TestVerify(t *testing.T) {
	// Test that we mark a signature as valid when it is and invalid when it is not valid
	Curve := btcec.S256()
	testCases := []struct {
		name           string
		publicKey      string
		messageDigest  string
		signature      string
		expectedResult bool
	}{
		{ // TEST CASE #1:
			name:           "Confirms a Valid Signature",
			publicKey:      "03defdea4cdb677750a420fee807eacf21eb9898ae79b9768766e4faa04a2d4a34",
			messageDigest:  "4df3c3f68fcc83b27e9d42c90431a72499f17875c81a599b566c9889b9696703",
			signature:      "00000000000000000000003b78ce563f89a0ed9414f5aa28ad0d96d6795f9c6302a8dc32e64e86a333f20ef56eac9ba30b7246d6d25e22adb8c6be1aeb08d49d",
			expectedResult: true,
		},
		{ // TEST CASE #2:
			name:           "Confirms a Valid Signature",
			publicKey:      "031b84c5567b126440995d3ed5aaba0565d71e1834604819ff9c17f5e9d5dd078f",
			messageDigest:  "0000000000000000000000000000000000000000000000000000000000000000",
			signature:      "52818579aca59767e3291d91b76b637bef062083284992f2d95f564ca6cb4e3530b1da849c8e8304adc0cfe870660334b3cfc18e825ef1db34cfae3dfc5d8187",
			expectedResult: true,
		},
		// { // TEST CASE #3:
		// 	name:           "Fails signature with public key not on the curve",
		// 	publicKey:      "03eefdea4cdb677750a420fee807eacf21eb9898ae79b9768766e4faa04a2d4a34",
		// 	messageDigest:  "4df3c3f68fcc83b27e9d42c90431a72499f17875c81a599b566c9889b9696703",
		// 	signature:      "00000000000000000000003b78ce563f89a0ed9414f5aa28ad0d96d6795f9c6302a8dc32e64e86a333f20ef56eac9ba30b7246d6d25e22adb8c6be1aeb08d49d",
		// 	expectedResult: false,
		// },
		{ // TEST CASE #4: FAILING
			name:           "Fails signature with incorrect R residuosity",
			publicKey:      "02dff1d77f2a671c5f36183726db2341be58feae1da2deced843240f7b502ba659",
			messageDigest:  "243f6a8885a308d313198a2e03707344a4093822299f31d0082efa98ec4e6c89",
			signature:      "2a298dacae57395a15d0795ddbfd1dcb564da82b0f269bc70a74f8220429ba1dfa16aee06609280a19b67a24e1977e4697712b5fd2943914ecd5f730901b4ab7",
			expectedResult: false,
		},
		{ // TEST CASE #5:
			name:           "Fails signature with negated message",
			publicKey:      "03fac2114c2fbb091527eb7c64ecb11f8021cb45e8e7809d3c0938e4b8c0e5f84b",
			messageDigest:  "5e2d58d8b3bcdf1abadec7829054f90dda9805aab56c77333024b9d0a508b75c",
			signature:      "00da9b08172a9b6f0466a2defd817f2d7ab437e0d253cb5395a963866b3574bed092f9d860f1776a1f7412ad8a1eb50daccc222bc8c0e26b2056df2f273efdec",
			expectedResult: false,
		},
		{ // TEST CASE #6:
			name:           "Fails signature with negated s value",
			publicKey:      "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",
			messageDigest:  "0000000000000000000000000000000000000000000000000000000000000000",
			signature:      "787a848e71043d280c50470e8e1532b2dd5d20ee912a45dbdd2bd1dfbf187ef68fce5677ce7a623cb20011225797ce7a8de1dc6ccd4f754a47da6c600e59543c",
			expectedResult: false,
		},
		{ // TEST CASE #7: FAILING
			name:           "Fails signature with sG - eP is infinite",
			publicKey:      "02dff1d77f2a671c5f36183726db2341be58feae1da2deced843240f7b502ba659",
			messageDigest:  "243f6a8885a308d313198a2e03707344a4093822299f31d0082efa98ec4e6c89",
			signature:      "00000000000000000000000000000000000000000000000000000000000000009e9d01af988b5cedce47221bfa9b222721f3fa408915444a4b489021db55775f",
			expectedResult: false,
		},
		{ // TEST CASE #8:
			name:           "Fails signature with sG - eP is infinite",
			publicKey:      "02dff1d77f2a671c5f36183726db2341be58feae1da2deced843240f7b502ba659",
			messageDigest:  "243f6a8885a308d313198a2e03707344a4093822299f31d0082efa98ec4e6c89",
			signature:      "0000000000000000000000000000000000000000000000000000000000000001d37ddf0254351836d84b1bd6a795fd5d523048f298c4214d187fe4892947f728",
			expectedResult: false,
		},
		{ // TEST CASE #9:
			name:           "Fails signature with x-coordinate of R not on curve",
			publicKey:      "02dff1d77f2a671c5f36183726db2341be58feae1da2deced843240f7b502ba659",
			messageDigest:  "243f6a8885a308d313198a2e03707344a4093822299f31d0082efa98ec4e6c89",
			signature:      "4a298dacae57395a15d0795ddbfd1dcb564da82b0f269bc70a74f8220429ba1d1e51a22ccec35599b8f266912281f8365ffc2d035a230434a1a64dc59f7013fd",
			expectedResult: false,
		},
		{ // TEST CASE #10:
			name:           "Fails signature with x-coordinate equal to finite field size",
			publicKey:      "02dff1d77f2a671c5f36183726db2341be58feae1da2deced843240f7b502ba659",
			messageDigest:  "243f6a8885a308d313198a2e03707344a4093822299f31d0082efa98ec4e6c89",
			signature:      "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc2f1e51a22ccec35599b8f266912281f8365ffc2d035a230434a1a64dc59f7013fd",
			expectedResult: false,
		},
		{ // TEST CASE #11:
			name:           "Fails signature with s equal to curve order",
			publicKey:      "02dff1d77f2a671c5f36183726db2341be58feae1da2deced843240f7b502ba659",
			messageDigest:  "243f6a8885a308d313198a2e03707344a4093822299f31d0082efa98ec4e6c89",
			signature:      "2a298dacae57395a15d0795ddbfd1dcb564da82b0f269bc70a74f8220429ba1dfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141",
			expectedResult: false,
		},
	}

	var messageDigest [32]byte
	// var secKey [33]byte
	var signature [64]byte
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			// Decode Test Parameters
			publicBytes, _ := hex.DecodeString(testCase.publicKey)
			msgBytes, _ := hex.DecodeString(testCase.messageDigest)
			sig, _ := hex.DecodeString(testCase.signature)
			copy(messageDigest[:], msgBytes)
			copy(signature[:], sig)
			public, _ := btcec.ParsePubKey(publicBytes, Curve)
			// fmt.Printf("Message Digest: %x\n", messageDigest)
			// fmt.Printf("Public: %x    %d\n", publicBytes, len(publicBytes))

			// privateKey, publicKey := btcec.PrivKeyFromBytes(Curve, private)
			// fmt.Printf("PrivateKey: %x\n", privateKey.D)
			// fmt.Println("PublicKey: ", publicKey.SerializeCompressed())
			// compressedKey := publicKey.SerializeCompressed()
			// fmt.Printf("SEC Compressed Key: %x\n", compressedKey)

			// t.Log("Hex to String: ", hex.EncodeToString(signature[:]))
			// t.Log("Expected Sign: ", testCase.expectedSignature)

			// Verify we generate the expected signature
			// copy(secKey[:], publicBytes)
			publicKey := public.ToECDSA()
			result, err := schnorr.VerifySignature(publicKey, messageDigest, signature)
			if result != testCase.expectedResult { // || err != nil
				t.Fatalf("Did not confirm/deny validity of signature as expected: Want: %t    Got: %t   Error: %s", testCase.expectedResult, result, err)
			} else {
				t.Logf("SUCCESS: Expected verify result. Schnorr Signature is valid: %t    Error: %s", result, err)
			}
		})
	}
}
