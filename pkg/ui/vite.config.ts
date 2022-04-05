import {defineConfig, searchForWorkspaceRoot} from 'vite'
import react from '@vitejs/plugin-react'
// @ts-ignore
import svgrPlugin from 'vite-plugin-svgr'


export default defineConfig({
    plugins: [
        react(),
        svgrPlugin({
            svgrOptions: {
                icon: true
            }
        })
    ],
    resolve: {
        alias: [
            {find: /^~/, replacement: ''}
        ],
    },
    css: {
        preprocessorOptions: {
            less: {
                javascriptEnabled: true,
            }
        }
    }

})
