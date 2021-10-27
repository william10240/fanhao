module.exports = {
    publicPath: process.env.NODE_ENV === 'production'
        ? '/static/'
        : '/',
    productionSourceMap: false,
}