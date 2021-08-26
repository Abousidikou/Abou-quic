onmessage = function (ev) {
	let start = performance.now()
    fetch("https://monitor.uac.bj:4448/upload",{method: 'post', body: ev.data.toSend}).then(response=>response.text()).then(data=> {
    let v = performance.now() - start
    //console.log(v+ " ms")
    let bms = (ev.data.initialMessageSize*8)/v
    //console.log(bms+" bits/ms")
    let bs = bms*1000
    //console.log(bs+" bits/s")
    let mbs = bs/1000000
    //console.log(Math.ceil(mbs)+" Mbps")
    //console.log(data)
     postMessage({
        'AppInfo': {
          'ElapsedTime': v,  // ms
          'Speed': Math.ceil(mbs)+ " Mbps",
        },
      })
    }).catch(err=>console.log(err));

}