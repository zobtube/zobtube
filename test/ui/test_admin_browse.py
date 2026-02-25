"""Admin user: browse and navigate to all sections including adm and upload."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, SEEDED_IDS, login_admin


def test_admin_navigation_main_sections(page: Page):
    """Admin can navigate to Home, Movies, Categories, Channels, Actors, Clips."""
    login_admin(page)
    checks = [
        ("/", "heading", "Latest Trending Videos"),
        ("/movies", "heading", "Movies"),
        ("/categories", "heading", "Categories"),
        ("/channels", "heading", "Channels"),
        ("/actors", "heading", "Actors"),
        ("/clips", "text", "clip"),  # "All clips" or "No clip available"
    ]
    for path, role, name in checks:
        page.goto(BASE_URL + path)
        page.wait_for_load_state("networkidle")
        expect(page).to_have_url(BASE_URL + path)
        if role == "text":
            expect(page.get_by_text(name).first).to_be_visible()
        else:
            expect(page.get_by_role(role, name=name).first).to_be_visible()


def test_admin_adm_tab_navigation(page: Page):
    """Admin can navigate to adm sections."""
    login_admin(page)
    adm_paths = [
        ("/adm", "heading", "Overview"),
        ("/adm/videos", "heading", "Videos"),
        ("/adm/actors", "heading", "Actors"),
        ("/adm/channels", "heading", "Channels"),
        ("/adm/categories", "heading", "Categories"),
        ("/adm/users", "text", "User list"),
        ("/adm/task/home", "heading", "Task list"),
        ("/adm/tasks", "heading", "Task list"),
    ]
    for path, role, name in adm_paths:
        page.goto(BASE_URL + path)
        page.wait_for_load_state("networkidle")
        expect(page).to_have_url(BASE_URL + path)
        if role == "text":
            expect(page.get_by_text(name).first).to_be_visible()
        else:
            expect(page.get_by_role(role, name=name).first).to_be_visible()


def test_admin_upload_link_visible_and_navigable(page: Page):
    """Admin sees upload link and can navigate to /upload."""
    login_admin(page)
    page.goto(BASE_URL + "/")
    page.wait_for_load_state("networkidle")
    upload_link = page.get_by_role("link", name="uploads")
    expect(upload_link).to_be_visible()
    upload_link.click()
    page.wait_for_load_state("networkidle")
    expect(page).to_have_url(BASE_URL + "/upload")


def test_admin_admin_link_visible_and_navigable(page: Page):
    """Admin sees admin link and can navigate to /adm."""
    login_admin(page)
    page.goto(BASE_URL + "/")
    page.wait_for_load_state("networkidle")
    admin_link = page.get_by_role("link", name="admin")
    expect(admin_link).to_be_visible()
    admin_link.click()
    page.wait_for_load_state("networkidle")
    expect(page).to_have_url(BASE_URL + "/adm")


def test_admin_browse_actor_view_and_edit_link(page: Page):
    """Admin can view actor and see edit link."""
    login_admin(page)
    page.goto(BASE_URL + "/actor/" + SEEDED_IDS["actor_id"])
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("test")).to_be_visible()
    expect(page.get_by_role("link", name="Edit profile")).to_be_visible()


def test_admin_browse_channel_view_and_edit_link(page: Page):
    """Admin can view channel and see edit link."""
    login_admin(page)
    page.goto(BASE_URL + "/channel/" + SEEDED_IDS["channel_id"])
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("test")).to_be_visible()
    expect(page.get_by_role("link", name="Edit channel")).to_be_visible()


def test_admin_actors_list_shows_add_link(page: Page):
    """Admin sees Add actor link on actors list."""
    login_admin(page)
    page.goto(BASE_URL + "/actors")
    page.wait_for_load_state("networkidle")
    add_link = page.locator("a[href='/actor/new']")
    expect(add_link).to_be_visible()


def test_admin_channels_list_shows_add_link(page: Page):
    """Admin sees Add channel link on channels list."""
    login_admin(page)
    page.goto(BASE_URL + "/channels")
    page.wait_for_load_state("networkidle")
    add_link = page.locator("a[href='/channel/new']")
    expect(add_link).to_be_visible()
