from playwright.sync_api import Page, expect

BASE_URL = 'http://127.0.0.1:8080'

def login(page):
    page.goto(BASE_URL+'/auth')
    page.get_by_role("textbox", name="Username").click()
    page.get_by_role("textbox", name="Username").fill("non-admin")
    page.get_by_role("textbox", name="Password").click()
    page.get_by_role("textbox", name="Password").fill("non-admin")
    page.get_by_role("button", name="Sign in").click()
    expect(page).to_have_url(BASE_URL+'/')

def test_login(page: Page):
    return login(page)

def test_access(page: Page):
    # login
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
        401: {
            "GET": [
                "/adm",
                "/adm/categories",
                "/adm/videos",
                "/adm/actors",
                "/adm/channels",
                "/adm/tasks",
                "/adm/task/:id",
                "/actor/new",
                "/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/delete",
                "/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/edit",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/provider/:provider_slug",
                "/api/actor/link/:id/thumb",
                "/channel/new",
                "/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/edit",
                "/upload/",
                "/upload/preview/:filepath",
                "/adm/users",
                "/adm/user",
                "/adm/user/:id/delete",

            ],
            'DELETE': [
                "/api/actor/link/9ec8e7be-1dc9-11f0-b02c-305a3a05e04d",
                "/api/actor/alias/:id",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/category/:category_id",
                "/api/category/:id",
                "/api/category-sub/d17232c1-1512-46c9-8b7a-158c4d89df6b/thumb",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/actor/:actor_id",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/category/:category_id",
                "/api/upload/file",
            ],
            'PUT': [
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/category/:category_id",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/actor/:actor_id",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/category/:category_id",
            ],
            'POST': [
                "/actor/new",
                "/api/actor/",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/alias",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/link",
                "/api/actor/045e1b0e-1dc4-11f0-a04a-305a3a05e04d/thumb",
                "/api/category",
                "/api/category-sub",
                "/api/category-sub/d17232c1-1512-46c9-8b7a-158c4d89df6b/thumb",
                "/api/category-sub/d17232c1-1512-46c9-8b7a-158c4d89df6b/rename",
                "/channel/new",
                "/api/video",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/upload",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/thumb",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/migrate",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/generate-thumbnail/:timing",
                "/api/video/d8045d56-1dc4-11f0-9970-305a3a05e04d/rename",
                "/upload/import",
                "/api/upload/triage/folder",
                "/api/upload/triage/file",
                "/api/upload/file",
                "/api/upload/folder",
                "/adm/user",
            ],
            'HEAD': [
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
