"""Profile: Your stats page."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, login_non_admin


def test_profile_stats_page_loads(page: Page):
    """User can load profile stats page."""
    login_non_admin(page)
    page.goto(BASE_URL + "/profile/stats")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("heading", name="Your stats")).to_be_visible()
    expect(page.get_by_text("Unique videos")).to_be_visible()
    expect(page.get_by_text("Total view time")).to_be_visible()
