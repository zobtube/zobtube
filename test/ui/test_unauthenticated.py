"""Unauthenticated: SPA (GET /, GET /auth, GET /path) returns 200; API returns 401 JSON."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, SEEDED_IDS


def test_spa_served_without_auth(page: Page):
    """GET / and GET /auth serve the SPA shell (200)."""
    for path in ["", "/auth", "/actors", "/videos", "/movies", "/clips", "/categories", "/channels"]:
        r = page.request.get(BASE_URL + path)
        assert r.status == 200, f"GET {path} expected 200, got {r.status}"
        assert "ZobTube" in (r.text() or "")


def test_api_returns_401_when_unauthenticated(page: Page):
    """All /api/* endpoints return 401 JSON when not authenticated."""
    api_get_urls = [
        "/api/home",
        "/api/actor",
        "/api/actor/" + SEEDED_IDS["actor_id"],
        "/api/category",
        "/api/channel",
        "/api/channel/" + SEEDED_IDS["channel_id"],
        "/api/clip",
        "/api/movie",
        "/api/video",
        "/api/video/" + SEEDED_IDS["video_id"],
        "/api/profile",
        "/api/auth/me",
    ]
    for url in api_get_urls:
        r = page.request.get(BASE_URL + url)
        assert r.status == 401, f"GET {url} expected 401, got {r.status}"
        body = r.text() or ""
        assert "error" in body.lower() or body == "{}"

    # POST/DELETE/PUT without auth -> 401
    r = page.request.post(BASE_URL + "/api/auth/logout")
    assert r.status == 401
    r = page.request.post(BASE_URL + "/api/actor/", data={})
    assert r.status == 401
    r = page.request.delete(BASE_URL + "/api/actor/" + SEEDED_IDS["actor_id"])
    assert r.status == 401


def test_login_form_visible_at_auth(page: Page):
    """Login form is visible when visiting /auth unauthenticated."""
    page.goto(BASE_URL + "/auth")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("textbox", name="Username")).to_be_visible()
    expect(page.get_by_role("textbox", name="Password")).to_be_visible()
    expect(page.get_by_role("button", name="Sign in")).to_be_visible()
