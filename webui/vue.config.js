module.exports = {
    publicPath: "./",
    devServer: {
        port: 80,
        proxy: 'https://127.0.0.1:8888'
        //proxy: 'http://localhost:8080'
        //proxy: 'https://vpn.gz.junqirao.icu:60203'
    },
    pages: {
        index: {
            // entry for the page
            entry: 'src/main.js',
            // the source template
            template: 'public/index.html',
            // output as dist/index.html
            filename: 'index.html',
            // when using title option,
            // template title tag needs to be <title><%= htmlWebpackPlugin.options.title %></title>
            title: 'TunnHub WebUI',
        },
    }
}

