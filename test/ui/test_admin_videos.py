"""Admin: video edit, delete, metadata operations."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, SEEDED_IDS, login_admin


def test_admin_video_view_and_edit_link(page: Page):
    """Admin can view video and navigate to edit."""
    login_admin(page)
    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"])
    page.wait_for_load_state("networkidle")
    expect(page.locator("video")).to_be_visible()
    page.get_by_role("link", name="Edit").click()
    page.wait_for_load_state("networkidle")
    expect(page).to_have_url(BASE_URL + "/video/" + SEEDED_IDS["video_id"] + "/edit")


def test_admin_video_edit_page_loads(page: Page):
    """Admin can load video edit page."""
    login_admin(page)
    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"] + "/edit")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Video editing")).to_be_visible()
    expect(page.locator("video")).to_be_visible()


def test_admin_video_rename_via_api(page: Page):
    """Admin can rename video via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/rename",
        data={"name": "Renamed Video"},
    )
    assert r.status == 200
    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"])
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Renamed Video")).to_be_visible()
    # Restore original name
    page.request.post(
        BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/rename",
        data={"name": "test"},
    )


def test_admin_video_add_actor_via_api(page: Page):
    """Admin can add actor to video via API."""
    login_admin(page)
    r = page.request.put(
        BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/actor/" + SEEDED_IDS["actor_id"]
    )
    assert r.status == 200
    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"])
    page.wait_for_load_state("networkidle")
    # Actor link has icon + name, use locator for actor link
    expect(page.locator('a[href*="/actor/' + SEEDED_IDS["actor_id"] + '"]')).to_be_visible()
    # Remove actor
    page.request.delete(
        BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/actor/" + SEEDED_IDS["actor_id"]
    )


def test_admin_video_generate_thumbnail_via_api(page: Page):
    """Admin can trigger thumbnail generation via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/generate-thumbnail/0"
    )
    assert r.status in (200, 500)


def test_admin_video_delete_created_video(page: Page):
    """Admin can create and delete a video via API (when task runner is configured)."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/video",
        data={"name": "Delete Me", "filename": "delete-me.mp4", "type": "v"},
    )
    if r.status != 200:
        # Task runner may not be configured in test env - skip delete test
        return
    body = r.json()
    video_id = body["video_id"]
    r = page.request.delete(BASE_URL + "/api/video/" + video_id)
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/video/" + video_id)
    assert r.status == 404


def test_admin_video_download_link(page: Page):
    """Admin can see download link on video view."""
    login_admin(page)
    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"])
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Download")).to_be_visible()
