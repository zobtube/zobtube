"""Authenticated (admin) user: login and API access."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, login_admin


def test_login(page: Page):
    login_admin(page)


def test_authenticated_api_access(page: Page):
    """After login, API returns 200 for protected routes."""
    login_admin(page)
    r = page.request.get(BASE_URL + "/api/auth/me")
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/home")
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/actor")
    assert r.status == 200
