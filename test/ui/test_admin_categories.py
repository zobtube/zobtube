"""Admin: category management."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL, login_admin


def test_admin_categories_page_loads(page: Page):
    """Admin can load adm categories page."""
    login_admin(page)
    page.goto(BASE_URL + "/adm/categories")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_role("heading", name="Categories")).to_be_visible()


def test_admin_add_category_via_api(page: Page):
    """Admin can add category via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/category",
        data={"Name": "E2E Test Category"},
    )
    assert r.status == 200
    page.goto(BASE_URL + "/adm/categories")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("E2E Test Category")).to_be_visible()


def test_admin_add_category_sub_via_api(page: Page):
    """Admin can add category sub via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/category",
        data={"Name": "E2E Category For Sub"},
    )
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/category")
    assert r.status == 200
    categories = r.json()["items"]
    parent_id = None
    for cat in categories:
        n = cat.get("name") or cat.get("Name") or ""
        if n == "E2E Category For Sub":
            parent_id = cat.get("ID") or cat.get("id")
            break
    assert parent_id, "Category not found"
    r = page.request.post(
        BASE_URL + "/api/category-sub",
        data={"Name": "E2E Sub Item", "Parent": parent_id},
    )
    assert r.status == 200
    page.goto(BASE_URL + "/adm/categories")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("E2E Sub Item")).to_be_visible()


def test_admin_rename_category_sub_via_api(page: Page):
    """Admin can rename category sub via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/category",
        data={"Name": "E2E Rename Category"},
    )
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/category")
    categories = r.json()["items"]
    parent_id = next((c.get("ID") or c.get("id")) for c in categories if (c.get("name") or c.get("Name")) == "E2E Rename Category")
    r = page.request.post(
        BASE_URL + "/api/category-sub",
        data={"Name": "Original Sub", "Parent": parent_id},
    )
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/category")
    categories = r.json()["items"]
    cat = next(c for c in categories if (c.get("name") or c.get("Name")) == "E2E Rename Category")
    subs = cat.get("sub") or cat.get("Sub") or []
    sub_id = subs[0].get("ID") or subs[0].get("id")
    r = page.request.post(
        BASE_URL + "/api/category-sub/" + sub_id + "/rename",
        data={"title": "Renamed Sub"},
    )
    assert r.status == 200
    page.goto(BASE_URL + "/adm/categories")
    page.wait_for_load_state("networkidle")
    expect(page.get_by_text("Renamed Sub")).to_be_visible()


def test_admin_delete_empty_category_via_api(page: Page):
    """Admin can delete empty category via API."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/category",
        data={"Name": "E2E Delete Category"},
    )
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/category")
    categories = r.json()["items"]
    cat_id = next((c.get("ID") or c.get("id")) for c in categories if (c.get("name") or c.get("Name")) == "E2E Delete Category")
    r = page.request.delete(BASE_URL + "/api/category/" + cat_id)
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/category")
    remaining = r.json()["items"]
    assert not any((c.get("name") or c.get("Name")) == "E2E Delete Category" for c in remaining)


def test_admin_category_modal_add(page: Page):
    """Admin can add category via modal."""
    login_admin(page)
    page.goto(BASE_URL + "/adm/categories")
    page.wait_for_load_state("networkidle")
    page.get_by_role("button", name="Add category").click()
    page.wait_for_selector("#zt-add-category-modal.show", timeout=5000)
    page.locator("#zt-add-category-name").fill("E2E Modal Category")
    page.locator("#zt-add-category-modal").get_by_role("button", name="Create").click()
    page.wait_for_selector("#zt-add-category-modal", state="hidden", timeout=5000)
    expect(page.get_by_role("heading", name="E2E Modal Category", level=4)).to_be_visible()


def test_admin_add_sub_category_via_modal(page: Page):
    """Admin can add a category item (sub-category) via modal."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/category",
        data={"Name": "E2E Category For Sub Modal"},
    )
    assert r.status == 200
    page.goto(BASE_URL + "/adm/categories")
    page.wait_for_load_state("networkidle")
    page.locator("section").filter(
        has=page.get_by_role("heading", name="E2E Category For Sub Modal", level=4)
    ).locator("a.category-new").click()
    page.wait_for_selector("#zt-add-sub-modal.show", timeout=5000)
    page.locator("#zt-add-sub-name").fill("E2E Modal Sub Item")
    page.locator("#zt-add-sub-modal").get_by_role("button", name="Create").click()
    page.wait_for_selector("#zt-add-sub-modal", state="hidden", timeout=5000)
    expect(page.get_by_text("E2E Modal Sub Item")).to_be_visible()


def test_admin_edit_category_via_modal(page: Page):
    """Admin can edit a parent category name via modal."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/category",
        data={"Name": "E2E Edit Category Original"},
    )
    assert r.status == 200
    page.goto(BASE_URL + "/adm/categories")
    page.wait_for_load_state("networkidle")
    page.locator("section").filter(
        has=page.get_by_role("heading", name="E2E Edit Category Original", level=4)
    ).locator(".zt-edit-category-btn").click()
    page.wait_for_selector("#zt-edit-category-modal.show", timeout=5000)
    page.locator("#zt-edit-category-name").fill("E2E Edit Category Renamed")
    page.locator("#zt-edit-category-modal").get_by_role("button", name="Save").click()
    page.wait_for_selector("#zt-edit-category-modal", state="hidden", timeout=5000)
    expect(page.get_by_role("heading", name="E2E Edit Category Renamed", level=4)).to_be_visible()


def test_admin_edit_sub_category_via_modal(page: Page):
    """Admin can edit a category item name via modal."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/category",
        data={"Name": "E2E Edit Sub Parent"},
    )
    assert r.status == 200
    r = page.request.get(BASE_URL + "/api/category")
    categories = r.json()["items"]
    parent_id = next(
        (c.get("ID") or c.get("id"))
        for c in categories
        if (c.get("name") or c.get("Name")) == "E2E Edit Sub Parent"
    )
    r = page.request.post(
        BASE_URL + "/api/category-sub",
        data={"Name": "E2E Sub Original", "Parent": parent_id},
    )
    assert r.status == 200
    page.goto(BASE_URL + "/adm/categories")
    page.wait_for_load_state("networkidle")
    page.get_by_role("link", name="E2E Sub Original").click()
    page.wait_for_selector("#zt-edit-sub-modal.show", timeout=5000)
    page.locator("#zt-edit-sub-name").fill("E2E Sub Renamed")
    page.locator("#zt-edit-sub-modal").get_by_role("button", name="Save").click()
    page.wait_for_selector("#zt-edit-sub-modal", state="hidden", timeout=5000)
    expect(page.get_by_text("E2E Sub Renamed")).to_be_visible()
