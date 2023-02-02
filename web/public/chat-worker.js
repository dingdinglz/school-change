function sleep(d) {
    for(var t = Date.now();Date.now()-t<=d;);
}
onmessage = (e) => {
    console.log(e.data);
    for(;;)
    {
        sleep(2000);
        var xlr = new XMLHttpRequest();
        xlr.open("GET","/api/get/chat",false);
        xlr.send(null);
        postMessage(xlr.responseText)
    }
}