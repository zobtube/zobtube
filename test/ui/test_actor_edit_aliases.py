"""Actor edit page: create and delete aliases via UI."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, login_admin


def test_actor_edit_create_alias_and_persists_after_reload(page: Page):
    """On /actor/UUID/edit, create alias via UI, reload, alias still present."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/actor/",
        data={"name": "Alias Create Target", "sex": "f"},
    )
    assert r.status == 200
    actor_id = r.json()["result"]

    page.goto(BASE_URL + "/actor/" + actor_id + "/edit")
    page.wait_for_load_state("networkidle")

    # Click "create alias" (button that opens Add alias modal)
    page.locator('button[data-bs-target="#addActorAliasModal"]').click()
    page.get_by_label("Alias").fill("E2E Test Alias")
    page.get_by_role("button", name="Add").click()

    # Alias appears immediately (page may refresh)
    expect(page.get_by_text("E2E Test Alias", exact=True)).to_be_visible(timeout=5000)

    # Reload and ensure alias persists
    page.reload()
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("E2E Test Alias", exact=True)).to_be_visible()

    # Cleanup
    page.request.delete(BASE_URL + "/api/actor/" + actor_id)


def test_actor_edit_delete_alias_and_persists_after_reload(page: Page):
    """On /actor/UUID/edit, delete alias via UI, reload, alias still gone."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/actor/",
        data={"name": "Alias Delete Target", "sex": "m"},
    )
    assert r.status == 200
    actor_id = r.json()["result"]
    r = page.request.post(
        BASE_URL + "/api/actor/" + actor_id + "/alias",
        data={"alias": "E2E Alias To Delete"},
    )
    assert r.status in (200, 201)

    page.goto(BASE_URL + "/actor/" + actor_id + "/edit")
    page.wait_for_load_state("networkidle")

    # Verify alias is present, then click its delete button
    expect(page.get_by_text("E2E Alias To Delete", exact=True)).to_be_visible()
    page.locator(".chip").filter(has_text="E2E Alias To Delete").locator(".zt-alias-remove").click()

    # Alias disappears (page may refresh)
    expect(page.get_by_text("E2E Alias To Delete", exact=True)).to_have_count(0, timeout=5000)

    # Reload and ensure deletion persists
    page.reload()
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("E2E Alias To Delete", exact=True)).to_have_count(0)

    # Cleanup
    page.request.delete(BASE_URL + "/api/actor/" + actor_id)
