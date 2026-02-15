"""Admin: channel CRUD operations."""
import json

from playwright.sync_api import Page, expect

from conftest import BASE_URL, SEEDED_IDS, login_admin


def test_admin_create_channel_via_api(page: Page):
    """Admin can create channel via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/channel",
        data=json.dumps({"name": "E2E Test Channel"}),
        headers={"Content-Type": "application/json"},
    )
    assert r.status in (200, 201)
    body = r.json()
    channel_id = body["id"]
    page.goto(BASE_URL + "/channel/" + channel_id)
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("E2E Test Channel")).to_be_visible()


def test_admin_navigate_to_channel_create(page: Page):
    """Admin can navigate to channel create page."""
    login_admin(page)
    page.goto(BASE_URL + "/channel/new")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("New channel")).to_be_visible()
    # Use form-scoped selector: get_by_label("Name") can match hidden login form's Username
    expect(page.locator("#channel-create-form input[name='name']")).to_be_visible()


def test_admin_channel_view_and_edit(page: Page):
    """Admin can view channel and navigate to edit."""
    login_admin(page)
    page.goto(BASE_URL + "/channel/" + SEEDED_IDS["channel_id"])
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("test")).to_be_visible()
    page.get_by_role("link", name="Edit channel").click()
    page.wait_for_load_state("networkidle")
    expect(page).to_have_url(BASE_URL + "/channel/" + SEEDED_IDS["channel_id"] + "/edit")


def test_admin_channel_update_via_api(page: Page):
    """Admin can update channel name via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/channel",
        data=json.dumps({"name": "Channel To Rename"}),
        headers={"Content-Type": "application/json"},
    )
    assert r.status in (200, 201)
    channel_id = r.json()["id"]
    r = page.request.put(
        BASE_URL + "/api/channel/" + channel_id,
        data=json.dumps({"name": "Renamed Channel"}),
        headers={"Content-Type": "application/json"},
    )
    assert r.status == 200
    page.goto(BASE_URL + "/channel/" + channel_id)
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Renamed Channel")).to_be_visible()


def test_admin_channels_list(page: Page):
    """Admin can see channels in list."""
    login_admin(page)
    page.goto(BASE_URL + "/channels")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("link", name="test", exact=True)).to_be_visible()
