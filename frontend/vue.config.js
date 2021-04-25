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
                    '^/revsuit/api': '/revsuit/api'
                }
            }
        },
    },
    runtimeCompiler: true,
    filenameHashing: false,
    publicPath: ''
}
