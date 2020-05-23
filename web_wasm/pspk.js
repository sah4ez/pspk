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

function publishKey() {
	callGo("./publish");
	// console.log(name);
	// console.log(pub_key);
	// console.log(priv_key);
}
