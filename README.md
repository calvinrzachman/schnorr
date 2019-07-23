# Schnorr Signatures
Go implementation of Schnorr Signatures

The following describes how the mathematical concepts of elliptic curves and finite fields can be used to construct
a digital signature scheme.

Consider the following: We would like to send Bitcoin from Party A to Party B in such a way that only Party B is able to authorize the next moving of funds. Ideally, Party A would authorize this transfer and, by doing so, relinquish all control/claim giving full discretion on how the coins will be spent next to Party B. In practice, this is achieved using a public/private key pair and digital signatures


## DIGITAL SIGNATURES

A digital signature is a mathematical structure which verifies the authenticity of a digital document or message. Just like a regular signature, a digital signature provides *authentication* . Digital signatures have several additional properties including that, after producing a signature, the signer cannot deny having produced that signature *(Non-repudiation)* and they provide that the contents of a message have not been altered in any way *(Data Integrity)*.

In Bitcoin, coins (UTXO) are locked to addresses. Addresses can encode different things, but most commonly they encode public keys or scripts. If coins are locked to an address representing a public key, the act of spending coins becomes a problem of proving knowledge of the associated private key (every transfer eventually comes down to proving this). The most simple way to do this would be to simply announce the private key to the network. However, this is undesirable and prompts search for a better method. Digital signatures are this alternative.

A digital signaure allows for a party to prove knowledge of a secret without divulging that secret. In the context of Bitcoin, this rather interesting property enables data to be encumbered or associated with a public key in such a way that only the party with knowledge of the associated private key can authorize the transfer of funds.
    

The following description assumes some basic familiarity with how mathematical operation of addition/multiplication
is defined on an elliptic curve over a finite field.


## SCHNORR SIGNATURE

Schnorr Signatures is a particular type of digital signature scheme with notable improvements over the more traditionally
used and current signature scheme in Bitcoin - Elliptic Curve Digital Signatures. These improvements include, most notably,

* Batch Verification
* Signature Aggregation 

Like ECDSA, a Schnorr Signature is characterized by the tuple:

    (R, s)

    R - x-coordinate an elliptic curve point (x,y) known as the public nonce. This appears to be used to help blind the secret
    s - Blinded secret key. The value that lets you demonstrate knowledge of a secret without revealing the secret
    
where s is given by:

    s = k + e*x   (1)

with

    k - a uniformally random 256 bit integer known as a nonce,
    e - a commitment to the associated public key P, the x-cordinate of the public nonce R, and the message m
    x - a a uniformally random 256 bit integer known in this context as the private key

and R given by:

    R = k*G

with

    G - generator point from an elliptic curve


### KEY PAIR
let (x, P) be a public/private key pair with the public key *P* given by:

    P = x*G

for private key *x* and elliptic curve generator *G*

### NONCE
let (k, R) represent a public nonce and its associated private nonce according to:

    R = k*G

for nonce *k*, and elliptic curve generator *G*

The private key *x* is analogous to nonce *k* and public key *P* is analogous to the public nonce *R*.

    x --> k
    P --> R

Mathematically the contructs are the same. They differ in both semantics and how they are used in signature construction.
*(k, R)* is essentially an ephemeral public/private key pair used only for the purposes of creating a digital signature.

NOTE: A single linear equation with one unknown variable is solvable for that variable.

    s = k + e*x  --> x = (s - k)/e

I like to think of *k* as a blinding factor which allows us to make the value s public. As result, just like the private key, *k* must be kept secret.

### COMMITMENT

Unlike handwritten signatures (if you have terrible penmanship like I do), a digital signature should be specific to a particular message. A signature construction 

    s = k + x

is enough to ensure uniqueness and protection of the private key, however this signature does not commit to any particular message. 

A signature *(R, s)* on a message *m1* should not be valid for a message *m2*. This is important as we need to show that not only did the owner of a private key produce the signature, but also that the signature was produced unforgeably for the document being signed. To do this our signature needs to be associated to the data it purports to sign.

This can be achieved by choosing a scaling factor e, such that it incorporates the data we would like to commit to. In this case the message and public key. In practice, this is done by taking the 32 byte hash of the data: 

    Hash(P, R, m)

This incoporates the message to be signed into the signature equation and earns us the desired property that if the message changes, the signature changes.

### SIGNING AND VERIFICATION

Given a signature *(R, s)*, validation can be performed by computing the following:

    s*G = R + e*P

this comes from the putting the signing equation on the elliptic curve by multiplying by G

    s*G = k*G + e*x*G

but from (1) and (2) we can rewrite this as:

    s*G = R + e*P

Note that the equation above contains only information which is safe to be made public. Thus, it can act as method
of signature validation.

 - The signature (R, s) is by necessity public
 - The public key P is public. lol
 - The public nonce R is associated with the nonce k, but k is protected by the assumed hardness of the Discrete Logarithm problem (the same security as the private key).

The validator cannot construct R directly, as they do not know *k*. They must solve

    R = s*G + e*P
    
They can then verify that the x-coordinate of *R* matches the 
    
    
<----------- Future Work ------------>
BATCH VERIFICATION
CROSS INPUT AGGREGATION
