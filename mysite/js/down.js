
onmessage = function (ev) {
	let data = ev.data.cmd
	console.log(data)
	let i = 0
		while(i<10){
			postMessage(i)
			i+=1
		}

}



