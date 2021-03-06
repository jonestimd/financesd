module.exports = {
    presets: [
        '@babel/preset-react',
        [
            '@babel/preset-env',
            {
                targets: {'chrome': '60', 'node': 'current'},
                shippedProposals: true,
                useBuiltIns: 'entry',
                corejs: {version: 3, proposals: true}
            }
        ],
        '@babel/preset-typescript'
    ],
    plugins: [
        ['babel-plugin-transform-imports', {
            '@material-ui/core': {
                transform: '@material-ui/core/esm/${member}',
                preventFullImport: true
            },
            lodash: {
                transform: 'lodash/${member}',
                preventFullImport: true
            }
        }],
        ['@babel/plugin-proposal-decorators', {legacy: true}],
        ['@babel/plugin-proposal-class-properties', {loose: true}],
        ['@babel/plugin-proposal-private-methods', { loose: true }],
        // '@babel/plugin-transform-async-to-generator',
        '@babel/plugin-proposal-optional-chaining',
        '@babel/plugin-proposal-nullish-coalescing-operator'
    ],
    sourceMaps: true,
    retainLines: true
}
