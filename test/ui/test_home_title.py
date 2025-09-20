import pytest
from playwright.sync_api import Page, expect
from xprocess import ProcessStarter

BASE_URL = 'http://127.0.0.1:8069'

def test_has_title(page: Page):
    page.goto(BASE_URL+'/')

    # Expect a title "to contain" a substring.
    expect(page).to_have_title("ZobTube")
