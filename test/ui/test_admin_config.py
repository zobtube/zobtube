"""Admin: config pages."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, login_admin


def test_admin_config_auth_page_loads(page: Page):
    """Admin can load auth config page."""
    login_admin(page)
    page.goto(BASE_URL + "/adm/config/auth")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Authentication settings")).to_be_visible()


def test_admin_config_auth_api(page: Page):
    """Admin can fetch auth config API."""
    login_admin(page)
    r = page.request.get(BASE_URL + "/api/adm/config/auth")
    assert r.status == 200


def test_admin_config_provider_page_loads(page: Page):
    """Admin can load provider config page."""
    login_admin(page)
    page.goto(BASE_URL + "/adm/config/provider")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("heading", name="Providers")).to_be_visible()


def test_admin_config_provider_api(page: Page):
    """Admin can fetch provider config API."""
    login_admin(page)
    r = page.request.get(BASE_URL + "/api/adm/config/provider")
    assert r.status == 200
    body = r.json()
    assert "providers" in body


def test_admin_config_offline_page_loads(page: Page):
    """Admin can load offline config page."""
    login_admin(page)
    page.goto(BASE_URL + "/adm/config/offline")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("heading", name="Offline mode").first).to_be_visible()


def test_admin_config_offline_api(page: Page):
    """Admin can fetch offline config API."""
    login_admin(page)
    r = page.request.get(BASE_URL + "/api/adm/config/offline")
    assert r.status == 200
    body = r.json()
    assert "offline_mode" in body


def test_admin_config_auth_toggle_endpoint_exists(page: Page):
    """Auth config endpoint responds."""
    login_admin(page)
    r = page.request.get(BASE_URL + "/api/adm/config/auth")
    assert r.status == 200
