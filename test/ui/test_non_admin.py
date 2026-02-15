"""Non-admin user: can access normal API routes (200); admin API routes return 403."""
from playwright.sync_api import Page, expect

from conftest import BASE_URL


def login(page: Page):
    page.goto(BASE_URL + "/auth")
    page.get_by_role("textbox", name="Username").click()
    page.get_by_role("textbox", name="Username").fill("non-admin")
    page.get_by_role("textbox", name="Password").click()
    page.get_by_role("textbox", name="Password").fill("non-admin")
    page.get_by_role("button", name="Sign in").click()
    expect(page).to_have_url(BASE_URL + "/")


def test_login(page: Page):
    return login(page)


def test_access(page: Page):
    login(page)

    rc = {
        200: {
            "GET": [
                "",
                "/actors",
                "/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d",
                "/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/thumb",
                "/categories",
                #"/category/8c50735e-1dc4-11f0-b1fc-305a3a05e04d/thumb",
                "/channels",
                "/channel/8c50735e-1dc4-11f0-b1fc-305a3a05e04d",
                "/channel/8c50735e-1dc4-11f0-b1fc-305a3a05e04d/thumb",
                "/clips",
                "/movies",
                "/videos",
                "/video/d8045d56-1dc4-11f0-9970-305a3a05e04d",
                "/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/stream",
                #"/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/thumb",
                "/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/thumb_xs",
                "/profile",
            ],
            'POST': [
                #"/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/count-view",
            ],
        },
        403: {
            "GET": [
                "/api/adm",
                "/api/adm/category",
                "/api/adm/video",
                "/api/adm/actor",
                "/api/adm/channel",
                "/api/adm/config/auth",
                "/api/adm/config/auth/enable",
                "/api/adm/config/offline",
                "/api/adm/config/offline/enable",
                "/api/adm/config/provider",
                "/api/adm/config/provider/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/switch",
                "/api/adm/task",
                "/api/adm/task/home",
                "/api/adm/task/045e1b0e-1dc4-11f0-a04a-305a3a05e04d",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/edit",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/provider/mockprov",
                "/api/actor/link/9ec8e7be-1dc9-11f0-b02c-305a3a05e04d/thumb",
                "/api/adm/user",
            ],
            "DELETE": [
                "/api/actor/link/9ec8e7be-1dc9-11f0-b02c-305a3a05e04d",
                "/api/actor/alias/9ec8e7be-1dc9-11f0-b02c-305a3a05e04d",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/category/d17232c1-1512-46c9-8b7a-158c4d89df6b",
                "/api/category/d17232c1-1512-46c9-8b7a-158c4d89df6b",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/category/d17232c1-1512-46c9-8b7a-158c4d89df6b",
                "/api/upload/file",
                "/api/upload/triage/mass-action",
                "/api/adm/user/045e1b0e-1dc4-11f0-a04a-305a3a05e04d",
            ],
            "PUT": [
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/category/d17232c1-1512-46c9-8b7a-158c4d89df6b",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/category/d17232c1-1512-46c9-8b7a-158c4d89df6b",
            ],
            "POST": [
                "/api/actor/",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/alias",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/link",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/thumb",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/rename",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/description",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/merge",
                "/api/category",
                "/api/category-sub",
                "/api/category-sub/d17232c1-1512-46c9-8b7a-158c4d89df6b/thumb",
                "/api/category-sub/d17232c1-1512-46c9-8b7a-158c4d89df6b/rename",
                "/api/video",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/upload",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/thumb",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/migrate",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/generate-thumbnail/0",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/rename",
                "/api/upload/import",
                "/api/upload/triage/folder",
                "/api/upload/triage/file",
                "/api/upload/triage/mass-action",
                "/api/upload/file",
                "/api/upload/folder",
                "/api/adm/user",
            ],
            "HEAD": [
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d",
            ],
        },
    }

    for rc, methods in rc.items():
        for method, urls in methods.items():
            for url in urls:
                url = BASE_URL+url
                response = None
                match method:
                    case 'GET':
                        response = page.request.get(url)
                    case 'POST':
                        response = page.request.post(url)
                    case 'DELETE':
                        response = page.request.delete(url)
                    case 'HEAD':
                        response = page.request.head(url)
                    case 'PUT':
                        response = page.request.put(url)
                print(f"{method} {url} want: {rc} / got {response.status}")
                assert response.status == rc
