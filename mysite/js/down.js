
onmessage = function (ev) {
	console.log("In down  Worker")
	let start = performance.now()
	fetch("https://monitor.uac.bj:4450/download").then(response=>response.blob()).then(data=>{
		console.log('Fnished')
	}).catch(err=>console.log(err))
}




