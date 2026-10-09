import secrets
import string

from fastapi import FastAPI, HTTPException
from fastapi.responses import RedirectResponse
from pydantic import BaseModel

app = FastAPI(title="shipyard-lab python-api")

LINKS: dict[str, dict] = {}
ALPHABET = string.ascii_letters + string.digits


class NewLink(BaseModel):
    url: str


def make_code() -> str:
    return "".join(secrets.choice(ALPHABET) for _ in range(7))


@app.post("/links", status_code=201)
def create_link(body: NewLink):
    code = make_code()
    while code in LINKS:
        code = make_code()
    LINKS[code] = {"code": code, "url": body.url, "clicks": 0}
    return LINKS[code]


@app.get("/links")
def list_links():
    return list(LINKS.values())


@app.get("/links/{code}/stats")
def link_stats(code: str):
    link = LINKS.get(code)
    if link is None:
        raise HTTPException(status_code=404, detail="not found")
    return link


@app.get("/{code}")
def follow(code: str):
    link = LINKS.get(code)
    if link is None:
        raise HTTPException(status_code=404, detail="not found")
    link["clicks"] += 1
    return RedirectResponse(link["url"], status_code=307)
