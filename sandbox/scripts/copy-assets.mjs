// Self-hosts the runtimes: copies Pyodide, sql.js and the Python harness into public/vendor.
// The sandbox CSP allows connections only to its own origin, so nothing loads from a CDN.
import { cpSync, mkdirSync, readdirSync, rmSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const out = join(root, 'public/vendor');
rmSync(out, { recursive: true, force: true });

const pyodide = join(root, 'node_modules/pyodide');
mkdirSync(join(out, 'pyodide'), { recursive: true });
for (const f of readdirSync(pyodide)) {
  if (/\.(mjs|js|wasm|zip|json)$/.test(f) && f !== 'package.json' && !f.endsWith('.map')) {
    cpSync(join(pyodide, f), join(out, 'pyodide', f));
  }
}

mkdirSync(join(out, 'sqljs'), { recursive: true });
// The browser build (sql-wasm-browser.js) asks for sql-wasm-browser.wasm; copy both non-debug builds.
for (const f of ['sql-wasm.wasm', 'sql-wasm-browser.wasm']) {
  cpSync(join(root, 'node_modules/sql.js/dist', f), join(out, 'sqljs', f));
}

mkdirSync(join(out, 'py/cq'), { recursive: true });
cpSync(join(root, 'src/python/browser_harness.py'), join(out, 'py/browser_harness.py'));
cpSync(join(root, '../runners/python/harness/cq/__init__.py'), join(out, 'py/cq/__init__.py'));
console.log('assets: copied to', out);
