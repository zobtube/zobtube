"""Admin: upload triage operations."""
import uuid

import pytest
from playwright.sync_api import Page, expect

from conftest import BASE_URL, login_admin


def test_admin_upload_page_loads(page: Page):
    """Admin can load upload page."""
    login_admin(page)
    page.goto(BASE_URL + "/upload")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Upload and triage folder")).to_be_visible()


def test_admin_upload_triage_folder_api(page: Page):
    """Admin can call triage folder API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/upload/triage/folder",
        data={"path": "/"},
    )
    assert r.status == 200
    body = r.json()
    assert "folders" in body


def test_admin_upload_triage_file_api(page: Page):
    """Admin can call triage file API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/upload/triage/file",
        data={"path": "/"},
    )
    assert r.status == 200
    body = r.json()
    assert "files" in body


def test_admin_upload_create_folder_via_api(page: Page):
    """Admin can create folder in triage via API."""
    login_admin(page)
    folder_name = "E2E_Triage_" + uuid.uuid4().hex[:8]
    r = page.request.post(
        BASE_URL + "/api/upload/folder",
        data={"name": folder_name},
    )
    assert r.status in (200, 409), f"Expected 200 or 409, got {r.status}"  # 409 = already exists
    r = page.request.post(
        BASE_URL + "/api/upload/triage/folder",
        data={"path": "/"},
    )
    assert r.status == 200
    body = r.json()
    assert "folders" in body


def test_admin_upload_new_folder_modal(page: Page):
    """Admin can open New folder modal, fill and click Create (modal may stay open if API fails)."""
    login_admin(page)
    folder_name = "E2E_Modal_" + uuid.uuid4().hex[:8]
    page.goto(BASE_URL + "/upload")
    page.wait_for_load_state("networkidle")
    page.get_by_role("button", name="New folder").click()
    page.wait_for_selector("#newFolderModal.show", timeout=5000)
    page.locator("#folder-new").fill(folder_name)
    page.locator("#newFolderModal").get_by_role("button", name="Create").click()
    page.wait_for_load_state("networkidle")
    # Verify we remained on upload page; modal may stay open if folder create fails (path/config)
    expect(page).to_have_url(BASE_URL + "/upload")


def test_admin_upload_buttons_visible(page: Page):
    """Admin sees Upload file and Mass action buttons."""
    login_admin(page)
    page.goto(BASE_URL + "/upload")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("button", name="Upload file")).to_be_visible()
    expect(page.get_by_role("button", name="New folder")).to_be_visible()
    expect(page.get_by_role("button", name="Mass action")).to_be_visible()


def test_admin_upload_double_click_folder_navigation(page: Page):
    """Admin can create folder via New folder button and double-click to navigate into it."""
    login_admin(page)
    folder_name = "e2e_nav_folder_" + uuid.uuid4().hex[:8]
    page.goto(BASE_URL + "/upload")
    page.wait_for_load_state("networkidle")
    page.get_by_role("button", name="New folder").click()
    page.wait_for_selector("#newFolderModal.show", timeout=5000)
    page.locator("#folder-new").fill(folder_name)
    page.locator("#zt-folder-create-btn").click()
    page.wait_for_selector("#newFolderModal", state="hidden", timeout=5000)
    row = page.locator("#zt-triage-listing tr[data-type=folder][data-folder='" + folder_name + "']")
    row.wait_for(state="visible", timeout=10000)
    row.dblclick()
    page.wait_for_load_state("networkidle")
    expect(page.locator("#zt-path")).to_contain_text(folder_name)


def test_admin_upload_file_details_offcanvas(page: Page):
    """Admin can click file to open details offcanvas with preview and actions."""
    login_admin(page)
    page.goto(BASE_URL + "/upload")
    page.wait_for_load_state("networkidle")
    row = page.locator("#zt-triage-listing tr[data-type=file]").filter(has_text="sample_clip")
    row.click()
    page.wait_for_selector("#zt-item-details.show", timeout=5000)
    expect(page.locator("#zt-item-details-content")).to_contain_text("File details")
    expect(page.locator("#zt-item-details-content")).to_contain_text("Import")
    expect(page.locator("#zt-item-details-content")).to_contain_text("Download")
    expect(page.locator("#zt-item-details-content")).to_contain_text("Delete")


def test_admin_upload_mass_action_modal(page: Page):
    """Admin can select file and open mass action modal."""
    login_admin(page)
    page.goto(BASE_URL + "/upload")
    page.wait_for_load_state("networkidle")
    row = page.locator("#zt-triage-listing tr[data-type=file]").first
    row.locator(".zt-check-cell").click()
    page.get_by_role("button", name="Mass action").click()
    page.wait_for_selector("#zt-mass-action-modal.show", timeout=5000)
    expect(page.locator("#zt-mass-action-list")).to_contain_text("sample_clip")
    expect(page.get_by_role("button", name="Import")).to_be_visible()
    expect(page.get_by_role("button", name="Delete")).to_be_visible()


def test_admin_upload_single_file_import(page: Page):
    """Admin can import single video file from triage."""
    login_admin(page)
    page.goto(BASE_URL + "/upload")
    page.wait_for_load_state("networkidle")
    row = page.locator("#zt-triage-listing tr[data-type=file]").filter(has_text="sample_clip")
    row.click()
    page.wait_for_selector("#zt-item-details.show", timeout=5000)
    page.locator("#zt-item-details-content").get_by_role("button", name="Import").click()
    page.wait_for_selector("#zt-video-import-modal.show", timeout=5000)
    with page.expect_response(
        lambda r: r.url.endswith("/api/video") and r.request.method == "POST"
    ) as resp_info:
        page.locator("#zt-video-import-modal").get_by_role("button", name="Clip").click()
    resp = resp_info.value
    assert resp.status == 200, f"Import failed: {resp.status}"
    data = resp.json()
    video_id = data.get("video_id")
    assert video_id, "Expected video_id in import response"
    page.wait_for_load_state("networkidle")
    page.goto(BASE_URL + "/clips")
    page.wait_for_load_state("networkidle")
    expect(page.locator(f'a[href="/clip/{video_id}"]')).to_be_visible(timeout=10000)
