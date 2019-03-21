var EC = require('elliptic').ec;

var ec = new EC('curve25519');

var alice = ec.genKeyPair();
var bob = ec.genKeyPair();

var aliceBob = alice.derive(bob.getPublic());
var bobAlice = bob.derive(alice.getPublic());

console.log('Both shared secrets are BN instances');
console.log(Buffer.from(aliceBob.toArray(), 'binary').toString('base64'));
console.log(Buffer.from(bobAlice.toArray(), 'binary').toString('base64'));