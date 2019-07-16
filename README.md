# schnorr-signatures
Go implementation of Schnorr signatures

# schnorr-signatures
Go implementation of Schnorr Signatures

Define the structure of a Schnorr Signature
The following describes how the mathematical concepts of elliptic curves and finite fields can be used to construct
a digital signature scheme.


## DIGITAL SIGNATURE

In Bitcoin, coins (UTXO) are locked to addresses. Most commonly, addresses can encode different things, but the most common are public keys and scripts. If coins are locked to an address representing a public key, the act of spending coins becomes a problem of proving knowledge of the associated private key. The most simple way to do this would be to simply announce the private key to the network.This is undesirable, and prompts search for a better method. Digital Signatures are this alternative.

A digital signature scheme is one which allows for a party to "sign" a message
A digital signaure allows for a party to prove knowledge of a secret without divulging that secret. This rather interesting
property enables data to be encumbered or associated with that Party ...
The mechanism used to enforce the association of Bitcoin to an address is a simple scripting language that is understood and
executed on every Bitcoin network node. 

    - Authorization
    - Nonrepudiation

The following description assumes some basic familiarity with how mathematical operation of addition/multiplication
is defined on an elliptic curve over a finite field.


## SCHNORR SIGNATURE

Schnorr Signatures is a particular type of digital signature scheme with notable improvements over the more traditionally
used and current signature scheme in Bitcoin - Elliptic Curve Digital Signatures. Like ECDSA, a Schnorr Signature is characterized by the tuple:

    (s, R)

    s - Blinded secret key. The value that lets you demonstrate knowledge of a secret without revealing the secret
    R - A public nonce. This appears to be used to help blind the secret

where s is given by:

    s = k + e*x

with

    k - a uniformally random 256 bit integer known as a nonce,
    e - a commitment to the associated public key P, the x-cordinate of the public nonce R, and the message m
    x - a a uniformally random 256 bit integer known in this context as the private key

and R given by:

    R = k*G

with

    G - generator point from an elliptic curve


### KEY PAIR
let (x, P) be a public/private key pair with the public key P given by:

    P = x*G

for private key x and elliptic curve generator G

### NONCE
let (k, R) represent a public nonce and its associated private nonce according to:

    R = k*G

for nonce k, and elliptic curve generator G

The private key x is analogous to nonce k and public key P is analogous to the public nonce R.

    x --> k
    P --> R

Mathematically the contructs are the same. They differ in both semantics and how they are used in signature construction.
(k, R) is essentially an ephemeral public/private key pair used only for the purposes of creating a digital signature.

NOTE: A single linear equation with one unknown variable is solvable for that variable.


    s = k + e*x  --> x = (s - k)/e

I like to think of k as a blinding factor which allows us to make the value s public. As result, just like the private key,
the nonce k must be kept secret.

### COMMITMENT

Unlike handwritten signatures (if you have terrible penmenship like I do) the digital signatures should be specific to a message. This is important as we need to show that not only did the owner of a private key produce the signature, but also that the signature was produced unforgeably for the document being signed. This can be achieved by choosing the scaling factor e, such that it incorporates the data we would like to commit to. In this case the message and public key. In practice, this is done by taking the 32 byte hash of the data

    Hash(P, R, m)

This incoporates the message to be signed into the signature equation. If the message changes, the signature changes.
The inclusion of P seems fine as we (Explore this further)

### SIGNING AND VERIFICATION

Given a signature (s, R), validation can be performed by computing the following:

    s*G = R + e*P

this comes from the putting the signing equation on the elliptic curve by multiplying by G

    s*G = k*G + e*x*G

but from (1) and (2) we can rewrite this as:

    s*G = R + e*P

Note that the equation above contains only information which is safe to be made public. Thus, it can act as method
of signature validation.
 - The signature (s,R) is by necessity public
 - The public key P is public. lol
 - The public nonce R is associated with the nonce k, but k is protected by the assumed hardness of the Discrete Logarithm problem (the same security as the private key).

The validator cannot construct R directly, as they do not know "k". They must solve

    R = s*G + e*P
    
    
<----------- Future Work ------------>
BATCH VERIFICATION
CROSS INPUT AGGREGATION
