const electron = require('electron');
const app = electron.app;
const BrowserWindow = electron.BrowserWindow;

let mainWindow = null;

app.on('window-all-closed', function() {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('ready', function() {
  mainWindow = new BrowserWindow({width: 800, height: 600,
    webPreferences: {
        nodeIntegration: true
    },
    'fullscreen': true, 'frame': false
    });
  mainWindow.loadFile('index.html')
  mainWindow.webContents.openDevTools()
  mainWindow.on('closed', function() {
    mainWindow = null;
  });
  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
});