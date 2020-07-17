const builder = require('electron-builder');

builder.build({
    config: {
        'appId': 'com.seigo2016.signage',
        'linux':{
            'target': {
                'target': 'zip',
                'arch': 'arm64'
            }
        }
        // 'mac':{
        //     'target':{
        //         'target': 'zip'
        //     }
        // }
    }
});