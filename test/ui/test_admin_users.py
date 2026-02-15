"""Admin: user management."""
import json

from playwright.sync_api import Page, expect

from conftest import BASE_URL, login_admin


def test_admin_users_page_loads(page: Page):
    """Admin can load adm users list."""
    login_admin(page)
    page.goto(BASE_URL + "/adm/users")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("User list")).to_be_visible()


def test_admin_user_new_page_loads(page: Page):
    """Admin can load add user page."""
    login_admin(page)
    page.goto(BASE_URL + "/adm/user")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("heading", name="Add a new user", exact=True)).to_be_visible()
    # Use placeholder to target add-user form (avoid duplicate #username from login form in SPA)
    expect(page.get_by_placeholder("my-new-user")).to_be_visible()
    expect(page.get_by_role("button", name="Create new user")).to_be_visible()


def test_admin_create_user_via_api(page: Page):
    """Admin can create user via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/adm/user",
        data=json.dumps({
            "username": "e2e_test_user",
            "password": "e2e_password",
            "admin": False,
        }),
        headers={"Content-Type": "application/json"},
    )
    assert r.status == 201
    body = r.json()
    user_id = body["id"]
    page.goto(BASE_URL + "/adm/users")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("e2e_test_user")).to_be_visible()
    # Cleanup - delete the created user (do not delete validation or non-admin)
    r = page.request.delete(BASE_URL + "/api/adm/user/" + user_id)
    assert r.status == 204


def test_admin_users_list_shows_existing(page: Page):
    """Admin sees existing users in list."""
    login_admin(page)
    page.goto(BASE_URL + "/adm/users")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("cell", name="validation")).to_be_visible()
    expect(page.get_by_role("cell", name="non-admin")).to_be_visible()


def test_admin_add_user_link_from_list(page: Page):
    """Admin can navigate to add user from list."""
    login_admin(page)
    page.goto(BASE_URL + "/adm/users")
    page.wait_for_load_state("networkidle")
    add_link = page.locator("a[href='/adm/user']")
    expect(add_link).to_be_visible()
    add_link.click()
    page.wait_for_load_state("networkidle")
    expect(page).to_have_url(BASE_URL + "/adm/user")
