const { spawn } = require('child_process');
const os = require('os');

let watch = false;
let mode = 'development';

for (let i = 0; i < process.argv.length; i++) {
    const arg = process.argv[i];
    if (arg === '--watch') {
        watch = true;
    } else if (arg === '--mode') {
        mode = process.argv[i + 1];
        i + 1;
    }
}

const startWebpack = (configFile) => {
    return new Promise(resolve => {
        let file = 'npx';
        const args = ['webpack', '--config', configFile];
        if (os.platform() === 'win32') {
            file = 'node_modules\\.bin\\webpack.cmd';
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

        const ps = spawn(file, args, { stdio: 'inherit', env: env });
        ps.on('close', code => {
            console.log(file, args, 'exit ' + code);
            if (code === 0) {
                resolve([true]);
            } else {
                resolve([false, 'Exit ' + code]);
            }
        });
    });
};

const startApp = () => {
    return startWebpack('webpack.app.js');
};

const startLogin = () => {
    return startWebpack('webpack.login.js');
};

(async () => {
    const results = await Promise.all([startApp(), startLogin()]);
    const appResult = results[0];
    const loginResult = results[1];

    if (!appResult[0]) {
        console.error('building app failed', appResult[1]);
        throw 'failed';
    }
    if (!loginResult[0]) {
        console.error('building login failed', loginResult[1]);
        throw 'failed';
    }
})();
