const { DefinePlugin } = require('webpack');
const packageJson = require('./package.json');
const crypto = require('crypto');
const path = require('path');

// package.jsonからバージョンを取得、もしくはランダムなハッシュ値を生成
const version = packageJson.version || crypto.randomBytes(4).toString('hex').substring(0, 7);
const hash = crypto.randomBytes(4).toString('hex').substring(0, 7);

module.exports = {
    // rules: [
    //   {
    //     test: /\.module\.css$/,
    //     use: [
    //       'style-loader',
    //       {
    //         loader: 'css-loader',
    //         options: {
    //           modules: true,
    //         },
    //       },
    //     ],
    //   },
    //   {
    //     test: /\.css$/,
    //     exclude: /\.module\.css$/,
    //     use: ['style-loader', 'css-loader'],
    //   },
    // ],
    webpack: {
        alias: {
            '@': path.resolve(__dirname, 'src'),
        },
        plugins: {
            add: [
              new DefinePlugin({
                'process.env.BUILD_TIME': JSON.stringify(new Date().toISOString()),
                'process.env.VERSION': JSON.stringify(version),
                'process.env.HASH': JSON.stringify(hash),
              }),
            ],
        },
    },
};
