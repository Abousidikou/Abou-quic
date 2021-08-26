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

function newWorker(){
	let worker = new Worker("js/down.js");

	worker.onmessage = function (ev) {
        if (ev.data === null) {
          console.log("ev.data error")
          return
        }
        console.log(ev.data);
      }

    // Kill the worker after the timeout. This force the browser to
      // close the WebSockets and prevent too-long tests.
      setTimeout(function () {
        worker.terminate()
      }, 2000)
      worker.postMessage({
        cmd: "start",
      })
}

function downTest(){
	console.log("Download_Test")
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
	document.getElementById("down1").innerHTML = Math.ceil(mbs)+" Mbps"; 
    document.getElementById("down2").innerHTML = Math.ceil(mbs)+" Mbps"; 
    document.getElementById("down3").innerHTML = Math.ceil(mbs)+" Mbps"; 
	}).catch(err=>console.log(err))
}

function upTest(){
    console.log("Upload_Test...")
    var initialMessageSize = 1 << 18 /* (1<<13) */
    var databuf = new Uint8Array(initialMessageSize)
    var bl = new Blob([databuf], {type: "application/octet-stream"});
    //var table = []
    let start = performance.now()
    let worker = new Worker("js/up.js");
    worker.postMessage({
        cmd: "start",
        toSend: bl,
        initialMessageSize: initialMessageSize,
    })
    console.log("Initial: "+initialMessageSize);
    worker.onmessage = function (ev) {
        if (ev.data === null) {
          console.log("ev.data error")
          return
        }
        //table.push(ev.data.AppInfo.Speed)
        document.getElementById("up1").innerHTML = ev.data.AppInfo.Speed;
        document.getElementById("up2").innerHTML = ev.data.AppInfo.Speed;
        document.getElementById("up3").innerHTML = ev.data.AppInfo.Speed;
        console.log("Speed: "+ev.data.AppInfo.Speed);
        console.log("Time: "+ev.data.AppInfo.ElapsedTime);
        console.log("Initial: "+initialMessageSize);
        if (ev.data.AppInfo.ElapsedTime <= 250 && initialMessageSize <= 1<<24 ) {
            initialMessageSize *= 2
            databuf = new Uint8Array(initialMessageSize)
            bl = new Blob([databuf], {type: "application/octet-stream"});
            worker.postMessage({
                cmd: "start",
                toSend: bl,
                initialMessageSize: initialMessageSize,
            })
        }else{
                worker.postMessage({
                    cmd: "start",
                    toSend: bl,
                    initialMessageSize: initialMessageSize,
                })
        }  
      }
    let duration = (performance.now() - start)
    console.log("Duration: "+duration)
    // Kill the worker after the timeout. This force the browser to
    // close the WebSockets and prevent too-long tests.
  setTimeout(function () {
    worker.terminate()
  }, 13000)

/*var moy = 0
table = table.slice(-3)
for(var i = 0, i < table.length; i++){
    moy += table[i]
}
moy /= table.length
console.log("Speed Up: "+moy)
      

	/*const initialMessageSize = 1 << 24 /* (1<<13) */
    /*const databuf = new Uint8Array(initialMessageSize) 
	// generate some data
	//console.log(typeof databuf)
	//console.log(databuf)
	
	var bl = new Blob([databuf], {type: "application/octet-stream"});
    let start = performance.now()
    fetch("https://monitor.uac.bj:4448/upload",{method: 'post', body: bl}).then(response=>response.text()).then(data=> {
    let v = performance.now() - start
    //console.log(v+ " ms")
    let bms = (initialMessageSize*8)/v
    //console.log(bms+" bits/ms")
    let bs = bms*1000
    //console.log(bs+" bits/s")
    let mbs = bs/1000000
    console.log(Math.ceil(mbs)+" Mbps")
    console.log(data)
    document.getElementById("up1").innerHTML = Math.ceil(mbs)+" Mbps"; 
    document.getElementById("up2").innerHTML = Math.ceil(mbs)+" Mbps";
    document.getElementById("up3").innerHTML = Math.ceil(mbs)+" Mbps";  
    }).catch(err=>console.log(err));*/
}





