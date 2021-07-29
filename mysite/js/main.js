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
	console.log("Download_Test...")
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
	console.log(Math.ceil(mbs)+" Mbps")
	console.log("Success")
	document.getElementById("down1").innerHTML = "Download: "+Math.ceil(mbs)+" Mbps";
	document.getElementById("down2").innerHTML = "Download: "+Math.ceil(mbs)+" Mbps";
	document.getElementById("down3").innerHTML = "Download: "+Math.ceil(mbs)+" Mbps";
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
    //console.log(data)
    let v = performance.now() - start
    //console.log(v+ " ms")
    let bms = (initialMessageSize*8)/v
    //console.log(bms+" bits/ms")
    let bs = bms*1000
    //console.log(bs+" bits/s")
    let mbs = bs/1000000
    console.log(Math.ceil(mbs)+" Mbps")
    console.log(data)
    document.getElementById("up1").innerHTML = "Upload: "+Math.ceil(mbs)+" Mbps";
    document.getElementById("up2").innerHTML = "Upload: "+Math.ceil(mbs)+" Mbps";
    document.getElementById("up3").innerHTML = "Upload: "+Math.ceil(mbs)+" Mbps";
    }).catch(err=>console.log(err));
}



function uploadA(){
	const initialMessageSize = 1 << 23 /* (1<<13) */
    const databuf = new Uint8Array(initialMessageSize) 
	// generate some data
	//console.log(typeof databuf)
	//console.log(databuf)
	console.log("Upload_Test...")
	var bl = new Blob([databuf], {type: "application/octet-stream"});
	const controller = new AbortController();
	const { signal } = controller;

	fetch("https://monitor.uac.bj:4448/upload", { method: 'post', body: bl, signal }).then(response => {
	    console.log("Request 1 is complete!");
	}).catch(e => {
	    console.warn(`Fetch 1 error: ${e.message}`);
	});


	// Wait 10 seconds to abort both requests
	setTimeout(() => controller.abort(), 10000);

}

