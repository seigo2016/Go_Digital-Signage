const electron = require('electron');
const remote = electron.remote;
const net = require('net');
const timeout = 30000
const jimp = require('jimp')
// const host = '192.168.0.5';
const host = 'signage-server.local'
// const host = 'localhost'
const port = '30000'
function socket_connect(client){
    client.connect(port, host, () => {;
        console.log('Connect: ' + host + ':' + port);
    })
}

function socket_close(client){
    setTimeout(socket_connect, timeout, client);
}

function main(){
    let client = new net.Socket();
    socket_connect(client);
    client.on('connect', ()=>{
        let prev_data = new Buffer("", 'base64')
        let all_data;
        client.on('data', data=>{
            let endindex = data.indexOf('\n\n')
            let s = new Buffer(data, 'base64')
            let image_element = document.getElementById("image");
            if (endindex != -1){
                all_data = Buffer.concat([prev_data, s])
                jimp.read(all_data).then(image => {
                    image.rotate(90).getBase64(jimp.MIME_PNG, function (err, src) {
                        image_element.src = src;
                    }).catch(function (err){
                        console.error(err)
                    });
                }).catch(function (err) {
                    console.error(err);
                });
                prev_data = new Buffer("", 'base64')
            }else{
                prev_data = Buffer.concat([prev_data, s]);
            }
        })
        client.on('close', function(){socket_close(client)});
    });
}
main()