function fetchAndInstantiate(url, importObject) {
	return fetch(url).then(response =>
		response.arrayBuffer()
	).then(bytes =>
		WebAssembly.instantiate(bytes, importObject)
	).then(results =>
		results.instance
	);
}

function callGo(name) {
	var go = new Go();
	var mod = fetchAndInstantiate(name, go.importObject);
	mod.then(function(instance) {
	  go.run(instance);
	});
}

function ToBase64(u8) {
	return btoa(String.fromCharCode.apply(null, u8));
}

var publish_name = "html-key";
function publishKey() {
	callGo("./publish");
}
