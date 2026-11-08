// Node-like globals for lessons. Output formatting follows Node's console.log
// for the common cases (strings raw at top level, quoted inside arrays/objects).
(function (host) {
  'use strict';
  function fmt(v, depth, nested) {
    if (typeof v === 'string') return nested ? "'" + v.replace(/'/g, "\\'") + "'" : v;
    if (v === null) return 'null';
    if (v === undefined) return 'undefined';
    if (typeof v === 'function') return '[Function: ' + (v.name || '(anonymous)') + ']';
    if (typeof v === 'bigint') return String(v) + 'n';
    if (typeof v !== 'object') return String(v);
    if (v instanceof Error) return v.stack || String(v);
    if (depth > 2) return Array.isArray(v) ? '[Array]' : '[Object]';
    if (v instanceof Map) {
      var m = []; v.forEach(function (val, k) { m.push(fmt(k, depth + 1, true) + ' => ' + fmt(val, depth + 1, true)); });
      return 'Map(' + v.size + ') ' + (m.length ? '{ ' + m.join(', ') + ' }' : '{}');
    }
    if (v instanceof Set) {
      var s = []; v.forEach(function (val) { s.push(fmt(val, depth + 1, true)); });
      return 'Set(' + v.size + ') ' + (s.length ? '{ ' + s.join(', ') + ' }' : '{}');
    }
    if (Array.isArray(v)) {
      return v.length ? '[ ' + v.map(function (x) { return fmt(x, depth + 1, true); }).join(', ') + ' ]' : '[]';
    }
    var keys = Object.keys(v);
    if (!keys.length) return '{}';
    return '{ ' + keys.map(function (k) {
      var key = /^[A-Za-z_$][\w$]*$/.test(k) ? k : "'" + k + "'";
      return key + ': ' + fmt(v[k], depth + 1, true);
    }).join(', ') + ' }';
  }
  function line(args) { return Array.prototype.map.call(args, function (a) { return fmt(a, 0, false); }).join(' ') + '\n'; }
  var g = globalThis;
  g.console = {
    log: function () { host.write(1, line(arguments)); },
    info: function () { host.write(1, line(arguments)); },
    error: function () { host.write(2, line(arguments)); },
    warn: function () { host.write(2, line(arguments)); },
  };
  var stdinText = null;
  function readStdin() { if (stdinText === null) stdinText = host.stdin(); return stdinText; }
  g.process = {
    argv: ['node', 'main.js'],
    env: {},
    exit: function (code) { host.exit(code === undefined ? 0 : code | 0); },
    stdout: { write: function (s) { host.write(1, String(s)); return true; } },
    stderr: { write: function (s) { host.write(2, String(s)); return true; } },
  };
  var fs = {
    readFileSync: function (path) {
      if (path === 0 || path === '/dev/stdin') return readStdin();
      throw new Error("ENOENT: no such file or directory, open '" + path + "'");
    },
  };
  g.require = function (name) {
    if (name === 'fs' || name === 'node:fs') return fs;
    throw new Error("Cannot find module '" + name + "'");
  };
  delete g.__host;
})(__host);
