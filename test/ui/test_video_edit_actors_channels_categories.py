"""Video edit page: add/remove actors, channel, categories via UI; reload and verify persistence."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, SEEDED_IDS, login_admin


def test_video_edit_add_actor_and_persists_after_reload(page: Page):
    """On /video/UUID/edit, add actor via UI, reload, actor still present."""
    login_admin(page)
    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"] + "/edit")
    page.wait_for_load_state("networkidle")

    # Open actor modal and add seeded actor by actor-id (filter(has_text="test") matches multiple)
    page.locator('button[data-bs-target="#actorSelectionModal"]').click()
    actor_chip = f'#actorSelectionModal .add-actor-list[actor-id="{SEEDED_IDS["actor_id"]}"]'
    page.locator(actor_chip).locator(".add-actor-add").click()
    page.locator("#actorSelectionModal").get_by_role("button", name="Close").click()

    # Actor appears in main area
    main_chip = f'.video-actor-list[actor-id="{SEEDED_IDS["actor_id"]}"]'
    expect(page.locator(main_chip)).to_be_visible(timeout=5000)

    # Reload and ensure actor persists
    page.reload()
    page.wait_for_load_state("networkidle")
    expect(page.locator(main_chip)).to_be_visible()

    # Cleanup
    page.request.delete(
        BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/actor/" + SEEDED_IDS["actor_id"]
    )


def test_video_edit_remove_actor_and_persists_after_reload(page: Page):
    """On /video/UUID/edit, remove actor via UI, reload, actor still gone."""
    login_admin(page)
    page.request.put(
        BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/actor/" + SEEDED_IDS["actor_id"]
    )

    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"] + "/edit")
    page.wait_for_load_state("networkidle")

    actor_chip = f'.video-actor-list[actor-id="{SEEDED_IDS["actor_id"]}"]'
    expect(page.locator(actor_chip)).to_be_visible()
    page.locator(actor_chip).get_by_role("button").click()

    expect(page.locator(actor_chip)).to_have_count(0, timeout=5000)

    page.reload()
    page.wait_for_load_state("networkidle")
    # After reload, unselected chips may stay in DOM with display:none
    expect(page.locator(actor_chip)).not_to_be_visible()


def test_video_edit_set_channel_and_persists_after_reload(page: Page):
    """On /video/UUID/edit, set channel via UI, reload, channel still present."""
    login_admin(page)
    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"] + "/edit")
    page.wait_for_load_state("networkidle")

    page.locator("#video-channel-edit").click()
    page.wait_for_selector("#editChannelModal.show", timeout=5000)
    page.locator("#channel-list").select_option(label="test")
    page.locator("#channel-send").click()

    # Page reloads after channel change
    page.wait_for_load_state("networkidle")
    expect(page.locator("#video-channel")).to_have_value("test")

    # Reload again and ensure channel persists
    page.reload()
    page.wait_for_load_state("networkidle")
    expect(page.locator("#video-channel")).to_have_value("test")

    # Cleanup: set channel back to None
    page.locator("#video-channel-edit").click()
    page.wait_for_selector("#editChannelModal.show", timeout=5000)
    page.locator("#channel-list").select_option(value="x")
    page.locator("#channel-send").click()


def test_video_edit_clear_channel_and_persists_after_reload(page: Page):
    """On /video/UUID/edit, clear channel via UI, reload, channel still None."""
    login_admin(page)
    # Ensure video has a channel first
    page.request.post(
        BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/channel",
        data={"channelID": SEEDED_IDS["channel_id"]},
    )

    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"] + "/edit")
    page.wait_for_load_state("networkidle")
    expect(page.locator("#video-channel")).to_have_value("test")

    page.locator("#video-channel-edit").click()
    page.wait_for_selector("#editChannelModal.show", timeout=5000)
    page.locator("#channel-list").select_option(value="x")
    page.locator("#channel-send").click()

    page.wait_for_load_state("networkidle")
    expect(page.locator("#video-channel")).to_have_value("None")

    page.reload()
    page.wait_for_load_state("networkidle")
    expect(page.locator("#video-channel")).to_have_value("None")


def _ensure_video_category_exists(page: Page, suffix: str = "0") -> str:
    """Create category and sub via API, return category-sub ID. suffix avoids collisions."""
    cat_name = "E2E Video Cat " + suffix
    sub_name = "E2E Video Sub " + suffix
    r = page.request.post(
        BASE_URL + "/api/category",
        data={"Name": cat_name},
    )
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/category")
    assert r.status == 200
    items = r.json().get("items") or r.json().get("categories") or []
    parent_id = None
    for cat in items:
        n = cat.get("name") or cat.get("Name") or ""
        if n == cat_name:
            parent_id = cat.get("ID") or cat.get("id")
            break
    assert parent_id, "Category not found"
    r = page.request.post(
        BASE_URL + "/api/category-sub",
        data={"Name": sub_name, "Parent": parent_id},
    )
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/category")
    items = r.json().get("items") or r.json().get("categories") or []
    for cat in items:
        if (cat.get("name") or cat.get("Name")) == cat_name:
            subs = cat.get("sub") or cat.get("Sub") or []
            for s in subs:
                if (s.get("name") or s.get("Name")) == sub_name:
                    return s.get("ID") or s.get("id")
    raise AssertionError("Category sub not found")


def test_video_edit_add_category_and_persists_after_reload(page: Page):
    """On /video/UUID/edit, add category via UI, reload, category still present."""
    login_admin(page)
    sub_name = "E2E Video Sub add"
    _ensure_video_category_exists(page, "add")

    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"] + "/edit")
    page.wait_for_load_state("networkidle")

    page.locator('button[data-bs-target="#categorySelectionModal"]').click()
    page.locator("#categorySelectionModal .add-category-list").filter(
        has_text=sub_name
    ).locator(".add-category-add").click()
    page.locator("#categorySelectionModal").get_by_role("button", name="Close").click()

    expect(
        page.locator(".video-category-list").filter(has_text=sub_name)
    ).to_be_visible(timeout=5000)

    page.reload()
    page.wait_for_load_state("networkidle")
    expect(
        page.locator(".video-category-list").filter(has_text=sub_name)
    ).to_be_visible()

    # Cleanup: remove category from video
    r = page.request.get(BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/edit")
    assert r.status == 200
    data = r.json()
    v = data.get("video") or data.get("Video") or {}
    cats = v.get("categories") or v.get("Categories") or []
    for c in cats:
        if (c.get("name") or c.get("Name")) == sub_name:
            cid = c.get("ID") or c.get("id")
            page.request.delete(
                BASE_URL + "/api/video/" + SEEDED_IDS["video_id"] + "/category/" + cid
            )
            break


def test_video_edit_remove_category_and_persists_after_reload(page: Page):
    """On /video/UUID/edit, remove category via UI, reload, category still gone."""
    login_admin(page)
    sub_name = "E2E Video Sub remove"
    sub_id = _ensure_video_category_exists(page, "remove")
    page.request.put(
        BASE_URL
        + "/api/video/"
        + SEEDED_IDS["video_id"]
        + "/category/"
        + sub_id
    )

    page.goto(BASE_URL + "/video/" + SEEDED_IDS["video_id"] + "/edit")
    page.wait_for_load_state("networkidle")

    expect(
        page.locator(".video-category-list").filter(has_text=sub_name)
    ).to_be_visible()
    page.locator(".video-category-list").filter(
        has_text=sub_name
    ).get_by_role("button").click()

    expect(
        page.locator(".video-category-list").filter(has_text=sub_name)
    ).to_have_count(0, timeout=5000)

    page.reload()
    page.wait_for_load_state("networkidle")
    # After reload, unselected chips stay in DOM with display:none; assert not visible
    expect(
        page.locator(".video-category-list").filter(has_text=sub_name)
    ).not_to_be_visible()
