import re
from playwright.sync_api import Page, expect

BASE_URL = 'http://127.0.0.1:8080'

def test_pages_unusable_if_unauthenticated(page: Page):
    methods = {
        "GET": [
            "",
            "/adm",
            "/adm/videos",
            "/adm/actors",
            "/adm/categories",
            "/adm/channels",
            "/adm/tasks",
            "/adm/task/:id",
            "/actors",
            "/actor/new",
            "/actor/:id",
            "/actor/:id/edit",
            "/actor/:id/thumb",
            "/actor/:id/delete",
            "/api/actor/:id/provider/:provider_slug",
            "/api/actor/link/:id/thumb",
            "/categories",
            "/category/:id",
            "/category-sub/:id/thumb",
            "/channels",
            "/channel/new",
            "/channel/:id",
            "/channel/:id/thumb",
            "/clips",
            "/movies",
            "/videos",
            "/video/:id",
            "/video/:id/edit",
            "/video/:id/stream",
            "/video/:id/thumb",
            "/video/:id/thumb_xs",
            "/upload/",
            "/upload/preview/:filepath",
            "/profile",
            "/adm/users",
            "/adm/user",
            "/adm/user/:id/delete",

        ],
        'POST': [
            "/actor/new",
            "/api/actor/",
            "/api/actor/:id/link",
            "/api/actor/:id/rename",
            "/api/actor/:id/thumb",
            "/api/actor/:id/alias",
            "/api/category",
            "/api/category-sub",
            "/api/category-sub/d17232c1-1512-46c9-8b7a-158c4d89df6b/rename",
            "/api/category-sub/d17232c1-1512-46c9-8b7a-158c4d89df6b/thumb",
            "/channel/new",
            "/api/video",
            "/api/video/:id/upload",
            "/api/video/:id/thumb",
            "/api/video/:id/migrate",
            "/api/video/:id/generate-thumbnail/:timing",
            "/api/video/:id/rename",
            "/api/video/:id/count-view",
            "/upload/import",
            "/api/upload/triage/folder",
            "/api/upload/triage/file",
            "/api/upload/file",
            "/api/upload/folder",
            "/adm/user",
        ],
        'DELETE': [
            "/api/actor/link/:id",
            "/api/actor/alias/:id",
            "/api/category/:id",
            "/api/video/:id",
            "/api/video/:id/actor/:actor_id",
            "/api/video/:id/category/:category_id",
            "/api/upload/file",
            "/api/actor/:id/category/:category_id",
            "/api/category-sub/:id/thumb",
        ],
        'HEAD': [
            "/api/video/:id",
        ],
        'PUT': [
            "/api/actor/:id/category/:category_id",
            "/api/video/:id/actor/:actor_id",
            "/api/video/:id/category/:category_id",
        ],
    }

    for method, urls in methods.items():
        print("checking method: "+method)
        for url in urls:
            url = BASE_URL+url
            print(f"checking url: {url}")
            response = None
            match method:
                case 'GET':
                    response = page.request.get(url, max_redirects=0)
                case 'POST':
                    response = page.request.post(url, max_redirects=0)
                case 'DELETE':
                    response = page.request.delete(url, max_redirects=0)
                case 'HEAD':
                    response = page.request.head(url, max_redirects=0)
                case 'PUT':
                    response = page.request.put(url, max_redirects=0)
            assert response.status == 302
            print(f"    response: {response.status}")
