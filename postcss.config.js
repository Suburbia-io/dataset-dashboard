module.exports = function({ options }) {
    const prod = options.mode === 'production'

    const config = {
        plugins: {
            'autoprefixer': {},
        },
    }

    if (prod) config.plugins['cssnano'] = { preset: 'default' }

    return config
}
