const electron = require('electron');
const remote = electron.remote;
const net = require('net');
const timeout = 30000
const jimp = require('jimp')
// const host = 'digital-signage-server.local';
const host = '192.168.0.11';
const port = '30000'
// const { width, height } = require("screenz");
function socket_connect(client){
    client.connect(port, host, () => {;
        console.log('Connect: ' + host + ':' + port);
    })
}

function socket_data(data){
    console.log(data)
    // let s = Buffer.from(data).toString('base64') ;
    let s = Buffer.from(data, 'base64')
    let image_element = document.getElementById("image");
    // image_element.src = "data:image/png;base64," + s;
    jimp.read(s).then(image => {
        console.log(s)
        image.rotate(90).getBase64(jimp.MIME_PNG, function (err, src) {
            console.log(src)
            image_element.src = "" + src;
        })
    }).catch(function (err) {
        console.error(err);
    });
}

function socket_close(client){
    setTimeout(socket_connect, timeout, client);
}

function main(){
    let client = new net.Socket();
    socket_connect(client);
    client.on('data', function(data){socket_data(data)});
    client.on('close', function(){socket_close(client)});
}
main()