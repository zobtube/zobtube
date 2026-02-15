from playwright.sync_api import Page, expect

from conftest import BASE_URL


def test_has_title(page: Page):
    page.goto(BASE_URL + "/")
    expect(page).to_have_title("ZobTube")
