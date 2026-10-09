from fastapi.testclient import TestClient

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
