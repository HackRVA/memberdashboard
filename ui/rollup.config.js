import filesize from 'rollup-plugin-filesize';
import { terser } from 'rollup-plugin-terser';
import replace from '@rollup/plugin-replace';
import { nodeResolve } from '@rollup/plugin-node-resolve';

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
                process: JSON.stringify({ env: { API: '/api' } }),
                preventAssignment: true
            }
        ),
        nodeResolve(),
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