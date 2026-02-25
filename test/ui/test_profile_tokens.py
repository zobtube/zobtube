"""Profile: API tokens page and Bearer auth."""
import json

import requests
from playwright.sync_api import Page, expect

from conftest import BASE_URL, login_non_admin


def test_profile_tokens_page_loads(page: Page):
    """User can load profile API tokens page."""
    login_non_admin(page)
    page.goto(BASE_URL + "/profile/tokens")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("heading", name="API tokens")).to_be_visible()
    expect(page.get_by_placeholder("Token name (e.g. My script)")).to_be_visible()
    expect(page.get_by_role("button", name="Create token")).to_be_visible()


def test_profile_tokens_create_and_show_once(page: Page):
    """User can create a token and see it once in a modal."""
    login_non_admin(page)
    page.goto(BASE_URL + "/profile/tokens")
    page.wait_for_load_state("networkidle")
    page.get_by_placeholder("Token name (e.g. My script)").fill("E2E Test Token")
    page.get_by_role("button", name="Create token").click()
    page.wait_for_selector("#zt-token-show-modal.show", timeout=5000)
    token_input = page.locator("#zt-token-show-value")
    expect(token_input).to_be_visible()
    token_value = token_input.input_value()
    assert len(token_value) == 64, f"Expected 64-char hex token, got length {len(token_value)}"
    expect(page.get_by_text("Copy this token now")).to_be_visible()
    page.locator("#zt-token-show-modal").get_by_role("button", name="Close").click()
    page.wait_for_selector("#zt-token-show-modal", state="hidden", timeout=2000)
    # Token should appear in the list
    expect(page.get_by_role("cell", name="E2E Test Token")).to_be_visible(timeout=5000)


def test_profile_tokens_delete(page: Page):
    """User can delete a token from the list."""
    login_non_admin(page)
    # Create token via API
    r = page.request.post(
        BASE_URL + "/api/profile/tokens",
        data=json.dumps({"name": "E2E Token To Delete"}),
        headers={"Content-Type": "application/json"},
    )
    assert r.status == 201, f"Create token failed: {r.status}"
    body = r.json()
    token_id = body["id"]
    # Go to tokens page and delete it
    page.goto(BASE_URL + "/profile/tokens")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("cell", name="E2E Token To Delete")).to_be_visible(timeout=5000)
    page.once("dialog", lambda d: d.accept())
    page.get_by_role("row", name="E2E Token To Delete").get_by_role("button", name="Delete").click()
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("cell", name="E2E Token To Delete")).not_to_be_visible()


def test_profile_tokens_bearer_auth(page: Page):
    """API accepts Bearer token for authentication (without session cookie)."""
    login_non_admin(page)
    r = page.request.post(
        BASE_URL + "/api/profile/tokens",
        data=json.dumps({"name": "E2E Bearer Test"}),
        headers={"Content-Type": "application/json"},
    )
    assert r.status == 201, f"Create token failed: {r.status}"
    token = r.json()["token"]
    # Request with only Bearer token (no cookies) to verify Bearer auth works
    r2 = requests.get(
        BASE_URL + "/api/profile/tokens",
        headers={"Authorization": f"Bearer {token}"},
        timeout=10,
    )
    assert r2.status_code == 200, f"GET with Bearer failed: {r2.status_code}"
    data = r2.json()
    assert "tokens" in data
    names = [t["name"] for t in data["tokens"] if t.get("name") == "E2E Bearer Test"]
    assert len(names) >= 1, "Expected created token in list when using Bearer auth"
