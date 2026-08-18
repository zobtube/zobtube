"""Stored XSS: names coming back from the API must never execute as markup.

Components build their markup by string concatenation and assign it to
innerHTML, so a name containing HTML has to come out as text. Each test plants
a payload that sets window.__xss when it executes, renders the page that shows
the name, and asserts the flag is still unset and the payload is visible as
literal text.
"""
import json

from playwright.sync_api import Page, expect

from conftest import BASE_URL, login_admin

# Fires on render (no click needed) and is unambiguous: <img> with a broken src.
PAYLOAD = '<img src=x onerror="window.__xss=1">'


def _xss_fired(page: Page) -> bool:
    return page.evaluate("() => !!window.__xss")


def test_actor_name_is_escaped_in_admin_list(page: Page):
    """zt-adm-object-list renders actor names as text, not markup."""
    login_admin(page)
    r = page.request.post(BASE_URL + "/api/actor/", data={"name": PAYLOAD, "sex": "f"})
    assert r.status == 200
    actor_id = r.json()["result"]
    try:
        page.goto(BASE_URL + "/adm/actors")
        page.wait_for_load_state("networkidle")
        assert not _xss_fired(page), "payload executed in the admin actor list"
        expect(page.get_by_text(PAYLOAD, exact=True)).to_be_visible()
    finally:
        page.request.delete(BASE_URL + "/api/actor/" + actor_id)


def test_actor_name_is_escaped_on_actor_page(page: Page):
    """zt-actor-view renders the actor name as text, not markup."""
    login_admin(page)
    r = page.request.post(BASE_URL + "/api/actor/", data={"name": PAYLOAD, "sex": "f"})
    assert r.status == 200
    actor_id = r.json()["result"]
    try:
        page.goto(BASE_URL + "/actor/" + actor_id)
        page.wait_for_load_state("networkidle")
        assert not _xss_fired(page), "payload executed on the actor page"
        expect(page.locator(".actor_name")).to_have_text(PAYLOAD)
    finally:
        page.request.delete(BASE_URL + "/api/actor/" + actor_id)


def test_channel_name_is_escaped_on_channel_page(page: Page):
    """zt-channel-view renders the channel name as text, not markup."""
    login_admin(page)
    r = page.request.post(
        BASE_URL + "/api/channel",
        data=json.dumps({"name": PAYLOAD}),
        headers={"Content-Type": "application/json"},
    )
    assert r.status in (200, 201)
    channel_id = r.json()["id"]
    try:
        page.goto(BASE_URL + "/channel/" + channel_id)
        page.wait_for_load_state("networkidle")
        assert not _xss_fired(page), "payload executed on the channel page"
        expect(page.locator(".actor_name")).to_have_text(PAYLOAD)
    finally:
        page.request.delete(BASE_URL + "/api/channel/" + channel_id)


def test_toast_message_is_not_parsed_as_html(page: Page):
    """sendToast treats its message as text unless the caller opts into HTML."""
    login_admin(page)
    page.goto(BASE_URL + "/")
    page.wait_for_load_state("networkidle")
    page.evaluate(
        "payload => sendToast('Test', '', 'bg-danger', payload)", PAYLOAD
    )
    assert not _xss_fired(page), "payload executed through the toast body"
    expect(page.locator(".zt-toast-body").last).to_have_text(PAYLOAD)


def test_toast_renders_a_working_link_when_html_is_allowed(page: Page):
    """The link-bearing toasts (import, assign image) must still produce a real <a>.

    Escaping the message wholesale broke these: the markup showed up as literal
    text instead of a clickable link. Callers that build a link opt in via the
    allowHtml flag, and the href has to survive escaping intact.
    """
    login_admin(page)
    page.goto(BASE_URL + "/")
    page.wait_for_load_state("networkidle")
    page.evaluate(
        """() => {
            var href = window.ztEscHtml(window.ztSafeUrl('/actor/abc-123?tab=photosets'));
            sendToast('Assign image', '', 'bg-success',
                      'Image assigned successfully <a href="' + href + '" target="_blank">View</a>',
                      true);
        }"""
    )
    link = page.locator(".zt-toast-body a").last
    expect(link).to_have_text("View")
    # &amp; / &#39; leaking into the href would break navigation
    expect(link).to_have_attribute("href", "/actor/abc-123?tab=photosets")


def test_escape_helper_leaves_ordinary_urls_intact(page: Page):
    """Escaping a normal URL must not alter it, or every generated link breaks."""
    login_admin(page)
    page.goto(BASE_URL + "/")
    page.wait_for_load_state("networkidle")
    for url in (
        "/video/abc-123/edit",
        "/photoset/9f8e/edit",
        "/actor/abc?tab=photosets",
        "/api/upload/preview/foo%2Fbar.mp4",
    ):
        assert page.evaluate("u => window.ztEscHtml(u)", url) == url


def test_escape_helper_covers_all_five_metacharacters(page: Page):
    """The shared helper escapes ' and > too, which the old per-file copies missed."""
    login_admin(page)
    page.goto(BASE_URL + "/")
    page.wait_for_load_state("networkidle")
    escaped = page.evaluate("""() => window.ztEscHtml('&<>"\\'')""")
    assert escaped == "&amp;&lt;&gt;&quot;&#39;"


def test_safe_url_neutralises_script_schemes(page: Page):
    """safeUrl blocks javascript: URLs while leaving ordinary links alone."""
    login_admin(page)
    page.goto(BASE_URL + "/")
    page.wait_for_load_state("networkidle")
    blocked = page.evaluate(
        """() => [
            'javascript:alert(1)',
            'JaVaScRiPt:alert(1)',
            'java\\tscript:alert(1)',
            ' javascript:alert(1)',
            'data:text/html,<script>alert(1)</script>',
        ].map(u => window.ztSafeUrl(u))"""
    )
    assert blocked == ["#"] * 5, blocked

    allowed = page.evaluate(
        """() => [
            '/video/abc/stream',
            'https://example.com/x',
            '?tab=photosets',
            '/foo/bar:baz',
        ].map(u => window.ztSafeUrl(u))"""
    )
    assert allowed == [
        "/video/abc/stream",
        "https://example.com/x",
        "?tab=photosets",
        "/foo/bar:baz",
    ], allowed
