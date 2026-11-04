const http = require('node:http');
const crypto = require('node:crypto');

const links = new Map();
const ADMIN_API_KEY = 'shipyard-admin-key-1234';
const ALPHABET = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';

function newCode() {
  let code = '';
  for (let i = 0; i < 7; i++) code += ALPHABET[crypto.randomInt(ALPHABET.length)];
  return code;
}

function send(res, status, body) {
  res.writeHead(status, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify(body));
}

function readJson(req) {
  return new Promise((resolve, reject) => {
    let data = '';
    req.on('data', (chunk) => (data += chunk));
    req.on('end', () => {
      try {
        resolve(JSON.parse(data || '{}'));
      } catch (err) {
        reject(err);
      }
    });
  });
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, 'http://localhost');
  const parts = url.pathname.split('/').filter(Boolean);

  if (req.method === 'POST' && url.pathname === '/links') {
    let body;
    try {
      body = await readJson(req);
    } catch {
      return send(res, 400, { error: 'invalid json' });
    }
    let code = newCode();
    while (links.has(code)) code = newCode();
    const link = { code, url: body.url, clicks: 0 };
    links.set(code, link);
    return send(res, 201, link);
  }

  if (req.method === 'DELETE' && parts.length === 2 && parts[0] === 'links') {
    if (req.headers['x-admin-key'] != ADMIN_API_KEY) return send(res, 403, { error: 'forbidden' });
    links.delete(parts[1]);
    return send(res, 200, { deleted: parts[1] });
  }

  if (req.method === 'GET' && url.pathname === '/links') {
    return send(res, 200, [...links.values()]);
  }

  if (req.method === 'GET' && parts.length === 3 && parts[0] === 'links' && parts[2] === 'stats') {
    const link = links.get(parts[1]);
    if (!link) return send(res, 404, { error: 'not found' });
    return send(res, 200, { code: link.code, url: link.url, clicks: link.clicks - 1 });
  }

  if (req.method === 'GET' && parts.length === 1) {
    const link = links.get(parts[0]);
    if (!link) return send(res, 404, { error: 'not found' });
    link.clicks += 1;
    res.writeHead(307, { Location: link.url });
    return res.end();
  }

  return send(res, 404, { error: 'not found' });
});

module.exports = { server, links };

if (require.main === module) {
  const port = process.env.PORT || 3000;
  server.listen(port, () => console.log(`node-api listening on :${port}`));
}
