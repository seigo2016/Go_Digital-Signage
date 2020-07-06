const builder = require('electron-builder');

builder.build({
    config: {
        'appId': 'com.seigo2016.signage',
        'linux':{
            'target': {
                'target': 'deb',
                'arch': 'arm64'
            }
        }
    }
});