import filesize from 'rollup-plugin-filesize';
import { terser } from 'rollup-plugin-terser';
import resolve from 'rollup-plugin-node-resolve';
import replace from '@rollup/plugin-replace';

export default {
    input: 'build/index.js',
    output: {
        file: 'dist/index.js',
        format: 'esm',
        inlineDynamicImports: true
    },
    onwarn(warning) {
        if (warning.code !== 'THIS_IS_UNDEFINED') {
            console.error(`(!) ${warning.message}`);
        }
    },
    plugins: [
        replace(
            { 
                'Reflect.decorate': 'undefined',
                process: JSON.stringify({ env: { API: '/api' } }) 
            }
        ),
        resolve(),
        terser({
            module: true,
            warnings: true,
            mangle: {
                properties: {
                    regex: /^__/,
                },
            },
        }),
        filesize({
            showBrotliSize: true,
        })
    ],
    watch: ["js", "build/*"]
}