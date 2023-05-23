// vue.config.js
const path = require('path');

// eslint-disable-next-line no-unused-vars
function resolve(dir) {
    return path.join(__dirname, dir)
}

module.exports = {
    devServer: {
        proxy: {
            '/api': {
                target: 'http://127.0.0.1:10000/revsuit/',
            }
        },
    },
    productionSourceMap: false,
    runtimeCompiler: true,
    filenameHashing: false,
    publicPath: '',
}
