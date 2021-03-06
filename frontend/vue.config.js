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
                target: 'http://localhost:10000',
                pathRewrite: {
                    '^/api': '/revsuit/api'
                }
            }
        },
    },
    productionSourceMap: false,
    runtimeCompiler: true,
    filenameHashing: false,
    publicPath: '/revsuit/admin/',
}
