from playwright.sync_api import Page, expect

BASE_URL = 'http://127.0.0.1:8069'

def login(page):
    page.goto(BASE_URL+'/auth')
    page.get_by_role("textbox", name="Username").click()
    page.get_by_role("textbox", name="Username").fill("validation")
    page.get_by_role("textbox", name="Password").click()
    page.get_by_role("textbox", name="Password").fill("validation")
    page.get_by_role("button", name="Sign in").click()
    expect(page).to_have_url(BASE_URL+'/')

def test_login(page: Page):
    return login(page)

def upload_common(page: Page, video_import_button, video_edit_button):
    # login
    login(page)

    # go to upload
    page.get_by_role("link", name="uploads").click()
    expect(page).to_have_url(BASE_URL+'/upload/')

    # upload
    with page.expect_file_chooser() as fc_info:
        page.get_by_role("button", name="Upload file").click()

    file_chooser = fc_info.value
    file_chooser.set_files("test/video/Big_Buck_Bunny_360_10s_1MB.mp4")
    expect(page.get_by_role("cell", name="Big_Buck_Bunny_360_10s_1MB.mp4")).to_be_visible()

    # import
    page.get_by_role("link", name="uploads").click()
    page.get_by_role("cell", name="Big_Buck_Bunny_360_10s_1MB.mp4").click()
    page.get_by_role("button", name=" Import").click()
    page.get_by_role("button", name=video_import_button).click()

    # open edit window
    with page.expect_popup() as page1_info:
        page.get_by_role("link", name="You can edit the video more").click()
    page1 = page1_info.value

    # check if video is imported
    expect(page1.get_by_text("Imported")).to_be_visible()

    # check if the video have the correct type
    expect(page1.get_by_role("button").filter(has_text=video_edit_button)).to_be_disabled()


def test_upload_clip(page: Page):
    upload_common(page, "Clip", "Change to Clip")

def test_upload_movie(page: Page):
    upload_common(page, "Movie", "Change to Movie")

def test_upload_video(page: Page):
    upload_common(page, "Video", "Change to Video")

def test_view(page: Page):
    # login
    login(page)

    page.get_by_role("link", name="").first.click()
    expect(page.locator("video")).to_be_visible()
