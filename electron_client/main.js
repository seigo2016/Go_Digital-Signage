const electron = require('electron');
const remote = electron.remote;
const net = require('net');
const timeout = 30000
// const host = 'digital-signage-server.local';
const host = 'localhost';
const port = '30000'

function socket_connect(client){
    client.connect(port, host, () => {;
        console.log('Connect: ' + host + ':' + port);
    })
}

function socket_data(data){
    console.log(data);
    let s = Buffer.from(data).toString('base64') ;
    let image_element = document.getElementById("image");
    image_element.src = "data:image/jpg;base64," + s;
}

function socket_close(client){
    setTimeout(socket_connect, timeout, client);
}

function main(){
    let host = '127.0.0.1';
    let port = 30000;

    let client = new net.Socket();
    socket_connect(client);
    client.on('data', function(data){socket_data(data)});
    client.on('close', function(){socket_close(client)});
}
main()