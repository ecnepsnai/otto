const { spawn } = require('child_process');
const os = require('os');

let watch = false;
let mode = 'development';

for (var i = 0; i < process.argv.length; i++) {
    const arg = process.argv[i];
    if (arg === '--watch') {
        watch = true;
    } else if (arg === '--mode') {
        mode = process.argv[i + 1];
        i + 1;
    }
}

var startWebpack = (configFile) => {
    return new Promise(resolve => {
        let file = 'npx';
        const args = ['webpack', '--config', configFile];
        if (os.platform() === 'win32') {
            file = "node_modules\\.bin\\webpack.cmd";
            args.splice(0, 1);
        }
        const env = process.env;

        if (mode === 'production') {
            args.push('--mode', 'production');
            env['NODE_ENV'] = 'production';
        }
        if (watch) {
            args.push('--watch');
        }

        console.log(file, args);

        const electron = spawn(file, args, { stdio: 'inherit', env: env });
        electron.on('close', () => {
            resolve();
        });
    });
}

var startApp = () => {
    const confFile = mode == 'production' ? 'webpack.app.production.js' : 'webpack.app.development.js';
    return startWebpack(confFile);
};

var startLogin = () => {
    const confFile = mode == 'production' ? 'webpack.login.production.js' : 'webpack.login.development.js';
    return startWebpack(confFile);
};

(async () => {
    await Promise.all([startApp(), startLogin()]);
})();
