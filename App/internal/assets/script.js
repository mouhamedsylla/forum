const goWasm = new Go()

WebAssembly.instantiateStreaming(fetch("/assets/main.wasm"), goWasm.importObject)
.then((result) => {
    goWasm.run(result.instance) 
    goWasm.instance.exports.InitSession()
    goWasm.instance.exports.InitApp()
    goWasm.instance.exports.BootstrapApplication()
});



