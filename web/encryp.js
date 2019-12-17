// some constants
const enc = new TextEncoder();
const dec = new TextDecoder();
const info = enc.encode("pspk_info");
let write = new Uint8Array(1);
write[0] = 0x01;


// Generate keys by ECDH on x25518
let seed = new Uint8Array(32);
crypto.getRandomValues(seed);
const keyPair = axlsign.generateKeyPair(seed);
// just for debug
// console.log("pub", ToBase64(keyPair.public));
// console.log("priv", ToBase64(keyPair.private));

function _base64ToUint8Array(base64) {
	var binary_string = window.atob(base64);
	var len = binary_string.length;
	var bytes = new Uint8Array(len);
	for (var i = 0; i < len; i++) {
		bytes[i] = binary_string.charCodeAt(i);
	}
	return bytes;
}

function ToBase64(u8) {
	return btoa(String.fromCharCode.apply(null, u8));
}

// Import shared key for HMAC usage
function importShared(shared) {
	return crypto.subtle.importKey(
		"raw", // raw format of the key - should be Uint8Array
		shared,
		{ // algorithm details
			name: "HMAC",
			hash: {name: "SHA-256"}
		},
		false, // export = false
		["sign", "verify"], // what this key can do
	);
}

// sigh shred key by write code
function signByShared(importedKey) {
	return crypto.subtle.sign("HMAC", importedKey, write)
}

// import material key for HKDF
function importMaterialKey(materialKey) {
	return crypto.subtle.importKey(
		"raw", // raw format of the key - should be Uint8Array
		new Uint8Array(materialKey),
		{ // algorithm details
			name: "HKDF",
			hash: {name: "SHA-256"}
		},
		false, // export = false
		["deriveBits", "deriveKey"], // what this key can do
	);
}

// Derived key by material key
function derivedKey(key) {
	const salt = new Uint8Array(32);
	return crypto.subtle.deriveBits(
		{
			name: "HKDF",
			hash: "SHA-256",
			salt: salt.buffer,
			info: info.buffer,
		},
		key,
		80 * 8, // 80 * 8 bytes or 640 bit
	);
}

// Import AES key for encryption
function importAESKey(key) {
	return crypto.subtle.importKey(
		"raw", // raw format of the key - should be Uint8Array
		key.slice(0, 32),
		{name: "AES-CBC"},
		false, // export = false
		["decrypt", "encrypt"], // what this key can do
	);
}

// Encrypt text as string by IV and imported AES key
function encryptByAESKeyOption(iv, text) {
	return function (key) {
		return crypto.subtle.encrypt(
			{
				name: "AES-CBC",
				iv: new Uint8Array(iv).slice(64),
			},
			key,
			enc.encode(text),
		);
	}
}

function decryptByAESKeyOption(iv, data) {
	return function (key) {
		return crypto.subtle.decrypt(
			{
				name: "AES-CBC",
				iv: new Uint8Array(iv).slice(64),
			},
			key,
			data,
		);
	}
}

function copyEncToClipboard() {
	let copyText = window.document.getElementById("copy_enc");
	copyText.select();
	copyText.setSelectionRange(0, 99999);
	document.execCommand('copy');

	const delEle = document.getElementById('copy_enc_alert');
	if (delEle !== null) {
		document.body.removeChild(delEle);
	}

	const ele = document.createElement('div');
	ele.id = 'copy_enc_alert'
	ele.textContent = "Copied encoded to clipboard!";
	ele.className = 'alert alert-primary';
	ele.setAttribute('role', 'alert');
	document.body.appendChild(ele);
}

function copyDecToClipboard() {
	let copyText = window.document.getElementById("copy_dec");
	copyText.select();
	copyText.setSelectionRange(0, 99999);
	document.execCommand('copy');

	const delEle = document.getElementById('copy_dec_alert');
	if (delEle !== null) {
		document.body.removeChild(delEle);
	}

	const ele = document.createElement('div');
	ele.id = 'copy_dec_alert'
	ele.textContent = "Copied text to clipboard!";
	ele.className = 'alert alert-primary';
	ele.setAttribute('role', 'alert');
	document.body.appendChild(ele);
}

// Encode encrypted data
function encodeEncrypt(data) {
	const bdata = new Uint8Array(data);
	let pub = new Uint8Array(keyPair.public.length + bdata.length);
	pub.set(keyPair.public);
	pub.set(bdata, keyPair.public.length);

	let copyText = window.document.getElementById("copy_enc");
	copyText.value = ToBase64(pub);
	copyText.select();
	copyText.setSelectionRange(0, 99999);

	console.log("encrypted data", ToBase64(pub));
}

// Encode encrypted data
function decodeDecrypt(data) {
	let copyText = window.document.getElementById("copy_dec");
	copyText.value = dec.decode(data);
	copyText.select();
	copyText.setSelectionRange(0, 99999);

	console.log("encrypted data", dec.decode(data));
}

function aesEncCatch(e) {
	console.error("encrypt failed:", e.message)
}

function aesDecCatch(e) {
	console.error("decrypt failed:", e.message)
}

// Encrtyp text by derived key via AES
function encryptByAES(text) {
	return function (key) {
		// just for debug
		// console.log("key", ToBase64(new Uint8Array(gkey)))
		let aesEnc = importAESKey(key).then(encryptByAESKeyOption(key, text));
		aesEnc.then(encodeEncrypt, aesEncCatch);
	}
}

// Encrtyp text by derived key via AES
function decryptByAES(data) {
	return function (key) {
		// just for debug
		// console.log("key", ToBase64(new Uint8Array(gkey)))
		let aesDec = importAESKey(key).then(decryptByAESKeyOption(key, data));
		aesDec.then(decodeDecrypt, aesDecCatch);
	}
}

function derivationCatch(e) {
	console.error("hmmm... derived failed:", e.message)
}

function encrypt(key, text) {
	const pub = _base64ToUint8Array(key);

	const shared = axlsign.sharedKey(keyPair.private, pub);
	// just for debug
	// console.log("shared", ToBase64(shared));
	let sign = importShared(shared).then(signByShared);

	let materialKey = sign.then(importMaterialKey);

	let derivation = materialKey.then(derivedKey);

	derivation.then(encryptByAES(text), derivationCatch)
}

function decrypt(priv, pub, data) {

	const shared = axlsign.sharedKey(priv, pub);
	// just for debug
	// console.log("shared", ToBase64(shared));
	let sign = importShared(shared).then(signByShared);

	let materialKey = sign.then(importMaterialKey);

	let derivation = materialKey.then(derivedKey);

	derivation.then(decryptByAES(data), derivationCatch)
}

function encryptText() {
	// Init elements of UI
	const pubName = window.document.getElementById("pub_name").value;
	const text = window.document.getElementById("text_enc").value;

	var xhr = new XMLHttpRequest();
	var url = "https://pspk.now.sh/";
	xhr.open("POST", url, true);
	xhr.onreadystatechange = function () {
		if (xhr.readyState === 4 && xhr.status === 200) {
			encrypt(JSON.parse(xhr.responseText).key, text)
		}
	};

	const data = JSON.stringify({"name": pubName});
	xhr.send(data);
}

function decryptText() {
	// Init elements of UI
	const keyBase64 = window.document.getElementById("private_key").value;
	const textBase64 = window.document.getElementById("text_dec").value;

	const messageBytes = _base64ToUint8Array(textBase64);
	const pubKey = messageBytes.slice(0, 32);
	const data = messageBytes.slice(32);

	const privKey = _base64ToUint8Array(keyBase64);

	decrypt(privKey, pubKey, data);
}

function findByName() {
	// Init elements of UI
	const name = window.document.getElementById("pub_name").value;


	var url = "https://pspk.now.sh/";
	var params = "?name_regex="+name;

	var table = new DataTable(".table", {
		ajax: {
			url: url+params,
			content: {
				type: "json",
				headings: true
			}
		}
	});
}

function clearPulbishKeyAlert() {
	const alertDel = document.getElementById('publish_alert');
	if (alertDel !== null) {
		document.body.removeChild(alertDel);
	}
	const publicDel = document.getElementById('public_key');
	if (publicDel !== null) {
		document.body.removeChild(publicDel);
	}
	const privateDel = document.getElementById('private_key');
	if (privateDel !== null) {
		document.body.removeChild(privateDel);
	}
	const delEle = document.getElementById('priv_key_alert');
	if (delEle !== null) {
		document.body.removeChild(delEle);
	}
}

function copyPrivateKey() {
	let copyText = window.document.getElementById("priv_key");
	copyText.select();
	copyText.setSelectionRange(0, 99999);
	document.execCommand('copy');

	clearPulbishKeyAlert();

	const ele = document.createElement('div');
	ele.id = 'priv_key_alert'
	ele.textContent = "Copied encoded!";
	ele.className = 'alert alert-primary';
	ele.setAttribute('role', 'alert');
	document.body.appendChild(ele);
}

function publishKey() {
	// Init elements of UI
	const name = window.document.getElementById("pub_name").value;
	const pub = ToBase64(new Uint8Array(keyPair.public));
	const key = ToBase64(new Uint8Array(keyPair.private));

	var xhr = new XMLHttpRequest();
	var url = "https://pspk.sah4ez.now.sh";
	xhr.open("POST", url, true);
	xhr.onreadystatechange = function () {
		if (xhr.readyState === 2 && xhr.status === 201) {
			clearPulbishKeyAlert();

			const alertEle = document.createElement('div');
			alertEle.id = 'publish_alert'
			alertEle.textContent = "Key published!";
			alertEle.className = 'alert alert-primary';
			alertEle.setAttribute('role', 'alert');
			document.body.appendChild(alertEle);

			const pubKey = document.getElementById('pub_key');
			pubKey.value = pub

			const privKey = document.getElementById('priv_key');
			privKey.value = key
		}

		if (xhr.readyState === 2 && xhr.status === 400) {
			clearPulbishKeyAlert();

			const alertEle = document.createElement('div');
			alertEle.id = 'publish_alert'
			alertEle.textContent = "Key with name: "+name+" exists";
			alertEle.className = 'alert alert-warning';
			alertEle.setAttribute('role', 'alert');
			document.body.appendChild(alertEle);
		}
	};

	const data = JSON.stringify({"name": name, "key": pub});
	xhr.send(data);
}

function scanPubKey(){
	const fileSelector = document.getElementById('file-selector');
	const text = window.document.getElementById("text_enc").value;

	const file = fileSelector.files[0];
	if (!file) {
		return;
	}

	var reader  = new FileReader();
	reader.addEventListener("load", function () {
		qrcode.callback = function(key) {
			if (key === 'error decoding QR Code') {
				const del = document.getElementById('dec_alert');
				if (del !== null) {
					document.body.removeChild(del);
				}
				const ele = document.createElement('div');
				ele.id = 'dec_alert'
				ele.textContent = key;
				ele.className = 'alert alert-warning';
				ele.setAttribute('role', 'alert');
				document.body.appendChild(ele);
			}
			console.log(key, text)
			encrypt(key, text)
		};
		qrcode.decode(reader.result);
	}, false);
	reader.readAsDataURL(file);
}
