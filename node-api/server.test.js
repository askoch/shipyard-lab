const test = require('node:test');
const assert = require('node:assert');
const { server, links } = require('./server');

test('create, follow, and stats respond', async (t) => {
  await new Promise((resolve) => server.listen(0, resolve));
  t.after(() => server.close());
  const base = `http://127.0.0.1:${server.address().port}`;

  const created = await fetch(`${base}/links`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url: 'https://example.com' }),
  });
  assert.strictEqual(created.status, 201);
  const { code } = await created.json();

  const followed = await fetch(`${base}/${code}`, { redirect: 'manual' });
  assert.strictEqual(followed.status, 307);

  const stats = await fetch(`${base}/links/${code}/stats`);
  assert.strictEqual(stats.status, 200);
  assert.ok(links.has(code));
});

test('unknown code returns 404', async (t) => {
  await new Promise((resolve) => server.listen(0, resolve));
  t.after(() => server.close());
  const base = `http://127.0.0.1:${server.address().port}`;
  const res = await fetch(`${base}/nope`, { redirect: 'manual' });
  assert.strictEqual(res.status, 404);
});
