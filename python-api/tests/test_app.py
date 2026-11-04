from fastapi.testclient import TestClient

from app import main
from app.main import LINKS, app

client = TestClient(app)


def setup_function():
    LINKS.clear()


def test_create_and_follow():
    r = client.post("/links", json={"url": "https://example.com"})
    assert r.status_code == 201
    code = r.json()["code"]

    r2 = client.get(f"/{code}", follow_redirects=False)
    assert r2.status_code == 307

    stats = client.get(f"/links/{code}/stats").json()
    assert stats["clicks"] == 1


def test_unknown_code():
    assert client.get("/nope", follow_redirects=False).status_code == 404


def test_healthz():
    r = client.get("/healthz", follow_redirects=False)
    assert r.status_code == 200
    assert r.json() == {"status": "ok"}


def test_healthz_not_shadowed_by_code():
    LINKS["healthz"] = {"code": "healthz", "url": "https://example.com", "clicks": 0}
    r = client.get("/healthz", follow_redirects=False)
    assert r.status_code == 200
    assert r.json() == {"status": "ok"}
    assert LINKS["healthz"]["clicks"] == 0


def test_reserved_code_not_generated(monkeypatch):
    codes = iter(["healthz", "abc1234"])
    monkeypatch.setattr(main, "make_code", lambda: next(codes))
    r = client.post("/links", json={"url": "https://example.com"})
    assert r.status_code == 201
    assert r.json()["code"] == "abc1234"
