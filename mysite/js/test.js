/*function downTest() {
    console.log("Test Start")
    
    fetch("https://monitor.uac.bj:4448/download").then(response => {
            response.blob()
        }).then(data => {
            
            console.log(data.size + " bytes")
            console.log(v + " ms")
        }).catch(error => {
            console.log("err: ", error)
        });
}*/
function downTest(){
	console.log("Test Start")
	let start = performance.now()
	fetch("https://monitor.uac.bj:4448/download").then(response=>response.blob()).then(data=>{
	let v = performance.now() - start
	//console.log(data)
	//console.log(data.size*8)
	//console.log(v+" ms")
	let bms = (data.size*8)/v
	//console.log(bms+" bit/ms")
	let bs = bms*1000
	//console.log(bs + " bits/s")
	let mbs = bs/1000000
	console.log(mbs+" Mbps")
	}).catch(err=>console.log(err))
}

function upTest(){
	const initialMessageSize = 1 << 23 /* (1<<13) */
    const databuf = new Uint8Array(initialMessageSize) 
	// generate some data
	//console.log(typeof databuf)
	//console.log(databuf)
	console.log("Upload_Test...")
	var bl = new Blob([databuf], {type: "application/octet-stream"});
    let start = performance.now()
    fetch("https://monitor.uac.bj:4448/upload",{method: 'post', body: bl}).then(response=>response.text()).then(data=> {
    console.log(data)
    let v = performance.now() - start
    console.log(v+ " ms")
    let bms = (initilMessageSize*8)/v
    console.log(bms+" bits/ms")
    let bs = bms*1000
    console.log(bs+" bits/s")
    let mbs = bs/1000000
    console.log(mbs+" Mbps")
    }).catch(err=>console.log(err));
}



function uploadA(){
// abort in 1 second
	const initialMessageSize = 1 << 23 /* (1<<13) */
    const databuf = new Uint8Array(initialMessageSize) 
	// generate some data
	//console.log(typeof databuf)
	//console.log(databuf)
	console.log("Upload_Test...")
	var bl = new Blob([databuf], {type: "application/octet-stream"});
let controller = new AbortController();
setTimeout(() => controller.abort(), 10000);

try {
  let response = await fetch('https://monitor.uac.bj:4448/uploadA', {
    signal: controller.signal
  });
} catch(err) {
  if (err.name == 'AbortError') { // handle abort()
    alert("Aborted!");
    console.log(data)
    let v = performance.now() - start
    console.log(v+ " ms")
    let bms = (initilMessageSize*8)/v
    console.log(bms+" bits/ms")
    let bs = bms*1000
    console.log(bs+" bits/s")
    let mbs = bs/1000000
    console.log(mbs+" Mbps")
  } else {
    throw err;
  }
}
}
