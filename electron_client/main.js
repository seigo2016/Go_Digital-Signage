const electron = require('electron');
const remote = electron.remote;
const net = require('net');
const timeout = 10 * 1000
const jimp = require('jimp')
// const host = '192.168.0.5';
const host = 'signage-server.local'
// const host = 'localhost'
const port = '30000'
function socket_connect(){
    let client = new net.Socket();
    client.connect(port, host, () => {;
        console.log('Connect: ' + host + ':' + port);
    })
    return client
}

function socket_close(){
    setTimeout(main, timeout);
}

function main(){
    let client = socket_connect()
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
                    })
                }).catch(function (err) {
                    console.error(err);
                });
                prev_data = new Buffer("", 'base64')
            }else{
                prev_data = Buffer.concat([prev_data, s]);
            }
        })
    });
    client.on('error', (err) => {
        console.log(err)
        socket_close()
    });
    client.on('close', function(){socket_close()});
}
main()