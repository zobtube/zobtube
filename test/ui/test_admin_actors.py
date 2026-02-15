"""Admin: actor CRUD operations."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, SEEDED_IDS, login_admin


def test_admin_create_actor_via_api(page: Page):
    """Admin can create actor via API and navigate to edit page."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/actor/",
        data={"name": "E2E Test Actor", "sex": "f"},
    )
    assert r.status == 200
    body = r.json()
    actor_id = body["result"]
    page.goto(BASE_URL + "/actor/" + actor_id + "/edit")
    page.wait_for_load_state("networkidle")
    expect(page.locator("#actor-name")).to_have_value("E2E Test Actor")


def test_admin_navigate_to_actor_edit(page: Page):
    """Admin can navigate to actor edit page."""
    login_admin(page)
    page.goto(BASE_URL + "/actor/" + SEEDED_IDS["actor_id"] + "/edit")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Edit actor information")).to_be_visible()
    expect(page.locator("#actor-name")).to_be_visible()


def test_admin_actor_rename_via_api(page: Page):
    """Admin can rename actor via API."""
    login_admin(page)
    # Create a new actor to rename (avoid mutating shared fixture)
    r = page.request.post(
        BASE_URL + "/api/actor/",
        data={"name": "Rename Target", "sex": "m"},
    )
    assert r.status == 200
    actor_id = r.json()["result"]
    r = page.request.post(
        BASE_URL + "/api/actor/" + actor_id + "/rename",
        data={"name": "Renamed Actor"},
    )
    assert r.status == 200
    page.goto(BASE_URL + "/actor/" + actor_id)
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Renamed Actor")).to_be_visible()
    # Cleanup
    page.request.delete(BASE_URL + "/api/actor/" + actor_id)


def test_admin_actor_add_alias_via_api(page: Page):
    """Admin can add alias via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/actor/",
        data={"name": "Alias Target", "sex": "f"},
    )
    assert r.status == 200
    actor_id = r.json()["result"]
    r = page.request.post(
        BASE_URL + "/api/actor/" + actor_id + "/alias",
        data={"alias": "Alt Name"},
    )
    assert r.status in (200, 201)
    page.goto(BASE_URL + "/actor/" + actor_id)
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Alt Name")).to_be_visible()
    page.request.delete(BASE_URL + "/api/actor/" + actor_id)


def test_admin_actor_delete_via_api(page: Page):
    """Admin can delete actor via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/actor/",
        data={"name": "Delete Target", "sex": "m"},
    )
    assert r.status == 200
    actor_id = r.json()["result"]
    r = page.request.delete(BASE_URL + "/api/actor/" + actor_id)
    assert r.status == 204
    r = page.request.get(BASE_URL + "/api/actor/" + actor_id)
    assert r.status == 404


def test_admin_actor_list_shows_actor(page: Page):
    """Admin can see actors in list."""
    login_admin(page)
    page.goto(BASE_URL + "/actors")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("test", exact=True)).to_be_visible()


def test_admin_actor_view_page(page: Page):
    """Admin can view actor and click edit."""
    login_admin(page)
    page.goto(BASE_URL + "/actor/" + SEEDED_IDS["actor_id"])
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("test", exact=True)).to_be_visible()
    page.get_by_role("link", name="Edit profile").click()
    page.wait_for_load_state("networkidle")
    expect(page).to_have_url(BASE_URL + "/actor/" + SEEDED_IDS["actor_id"] + "/edit")
