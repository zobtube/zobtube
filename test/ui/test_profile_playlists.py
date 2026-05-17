"""Profile: user playlists."""
import json
import re

from playwright.sync_api import Page, expect

from conftest import BASE_URL, SEEDED_IDS, login_non_admin


def test_profile_playlists_page_loads(page: Page):
    """User can load profile playlists page."""
    login_non_admin(page)
    page.goto(BASE_URL + "/profile/playlists")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("heading", name="Your playlists")).to_be_visible()
    expect(page.get_by_placeholder("Playlist name")).to_be_visible()
    expect(page.get_by_role("button", name="Create playlist")).to_be_visible()
    expect(page.get_by_role("link", name="Your playlists")).to_be_visible()


def test_profile_playlists_virtual_unseen_lists(page: Page):
    """Automatic unseen playlists are listed and cannot be deleted."""
    login_non_admin(page)
    page.goto(BASE_URL + "/profile/playlists")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("link", name="Unseen videos")).to_be_visible()
    expect(page.get_by_role("link", name="Unseen clips")).to_be_visible()
    expect(page.get_by_role("link", name="Unseen movies")).to_be_visible()
    unseen_row = page.get_by_role("row", name="Unseen videos")
    expect(unseen_row.get_by_role("button", name="Delete")).to_have_count(0)


def test_profile_playlists_create_and_view(page: Page):
    """User can create a playlist and open its detail page."""
    login_non_admin(page)
    page.goto(BASE_URL + "/profile/playlists")
    page.wait_for_load_state("networkidle")
    page.get_by_placeholder("Playlist name").fill("E2E Test Playlist")
    page.get_by_role("button", name="Create playlist").click()
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("cell", name="E2E Test Playlist")).to_be_visible(timeout=5000)
    page.get_by_role("link", name="E2E Test Playlist").click()
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("heading", name="E2E Test Playlist")).to_be_visible()
    expect(page.get_by_text("No videos in this playlist yet")).to_be_visible()


def test_profile_playlists_add_video_via_api(page: Page):
    """Video added via API appears on playlist detail page."""
    login_non_admin(page)
    r = page.request.post(
        BASE_URL + "/api/playlists",
        data=json.dumps({"name": "E2E Playlist With Video"}),
        headers={"Content-Type": "application/json"},
    )
    assert r.status == 201, f"Create playlist failed: {r.status}"
    playlist_id = r.json()["id"]
    video_id = SEEDED_IDS["video_id"]
    r2 = page.request.post(
        BASE_URL + f"/api/playlists/{playlist_id}/videos",
        data=json.dumps({"video_id": video_id}),
        headers={"Content-Type": "application/json"},
    )
    assert r2.status == 200, f"Add video failed: {r2.status}"
    page.goto(BASE_URL + f"/playlist/{playlist_id}")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("heading", name="E2E Playlist With Video")).to_be_visible()
    expect(page.locator("zt-video-tile")).to_have_count(1, timeout=5000)


def test_profile_playlists_play_navigates_with_query(page: Page):
    """Play playlist opens first video with ?playlist= and autoplay."""
    login_non_admin(page)
    r = page.request.post(
        BASE_URL + "/api/playlists",
        data=json.dumps({"name": "E2E Play Queue"}),
        headers={"Content-Type": "application/json"},
    )
    assert r.status == 201
    playlist_id = r.json()["id"]
    video_id = SEEDED_IDS["video_id"]
    r2 = page.request.post(
        BASE_URL + f"/api/playlists/{playlist_id}/videos",
        data=json.dumps({"video_id": video_id}),
        headers={"Content-Type": "application/json"},
    )
    assert r2.status == 200
    page.goto(BASE_URL + f"/playlist/{playlist_id}")
    page.wait_for_load_state("networkidle")
    page.get_by_role("button", name="Play playlist").click()
    page.wait_for_load_state("networkidle")
    expect(page).to_have_url(
        re.compile(rf"/video/{video_id}\?playlist={playlist_id}.*autoplay"),
        timeout=10000,
    )
    expect(page.locator("text=Playing from")).to_be_visible(timeout=5000)


def test_profile_playlists_delete(page: Page):
    """User can delete a playlist from the profile list."""
    login_non_admin(page)
    r = page.request.post(
        BASE_URL + "/api/playlists",
        data=json.dumps({"name": "E2E Playlist To Delete"}),
        headers={"Content-Type": "application/json"},
    )
    assert r.status == 201
    playlist_id = r.json()["id"]
    page.goto(BASE_URL + "/profile/playlists")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("cell", name="E2E Playlist To Delete")).to_be_visible(timeout=5000)
    page.once("dialog", lambda d: d.accept())
    page.get_by_role("row", name="E2E Playlist To Delete").get_by_role("button", name="Delete").click()
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("cell", name="E2E Playlist To Delete")).not_to_be_visible()
