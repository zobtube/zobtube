import os
import subprocess
from pathlib import Path

import pytest
import requests
from xprocess import ProcessStarter
from playwright.sync_api import Page, expect

# Per-worker URLs: master uses 6969, gw0->6970, gw1->6971, etc.
_worker = os.environ.get("PYTEST_XDIST_WORKER", "master")
if _worker == "master":
    _port = 6969
else:
    _num = int(_worker.replace("gw", ""))
    _port = 6970 + _num

BASE_URL = f"http://127.0.0.1:{_port}"

# IDs seeded by test/db/generate-empty.sh
SEEDED_IDS = {
    "actor_id": "045e1b0e-1dc4-11f0-a04a-305a3a05e04d",
    "channel_id": "8c50735e-1dc4-11f0-b1fc-305a3a05e04d",
    "video_id": "d8045d56-1dc4-11f0-9970-305a3a05e04d",
    "clip_id": "e8045d56-1dc4-11f0-9970-305a3a05e04d",
}


def login_admin(page: Page, base_url: str = BASE_URL) -> None:
    """Log in as admin user (validation/validation)."""
    page.goto(base_url + "/auth")
    page.get_by_role("textbox", name="Username").fill("validation")
    page.get_by_role("textbox", name="Password").fill("validation")
    page.get_by_role("button", name="Sign in").click()
    expect(page).to_have_url(base_url + "/")


def login_non_admin(page: Page, base_url: str = BASE_URL) -> None:
    """Log in as non-admin user (non-admin/non-admin)."""
    page.goto(base_url + "/auth")
    page.get_by_role("textbox", name="Username").fill("non-admin")
    page.get_by_role("textbox", name="Password").fill("non-admin")
    page.get_by_role("button", name="Sign in").click()
    expect(page).to_have_url(base_url + "/")


def _worker_resources():
    """Return port, db_path, media_path, xprocess_name for current worker."""
    worker = os.environ.get("PYTEST_XDIST_WORKER", "master")
    if worker == "master":
        return 6969, "/tmp/zt-db.sqlite3", "/tmp/zt-data", "zobtube"
    num = int(worker.replace("gw", ""))
    port = 6970 + num
    suffix = worker
    return (
        port,
        f"/tmp/zt-db-{suffix}.sqlite3",
        f"/tmp/zt-data-{suffix}",
        f"zobtube_{suffix}",
    )


@pytest.fixture(scope="session", autouse=True)
def zotbue_server(xprocess):
    port, db_path, media_path, xprocess_name = _worker_resources()

    # Prepare DB and media for this worker (generate-empty.sh)
    project_root = Path(__file__).resolve().parent.parent.parent
    env = os.environ.copy()
    env["ZT_SERVER_BIND"] = f"127.0.0.1:{port}"
    env["ZT_DB_CONNSTRING"] = db_path
    env["ZT_MEDIA_PATH"] = media_path
    subprocess.run(
        [str(project_root / "test" / "db" / "generate-empty.sh")],
        cwd=project_root,
        env=env,
        check=True,
    )

    class Starter(ProcessStarter):
        args = ["/tmp/zt"]
        env = {
            "ZT_SERVER_BIND": f"127.0.0.1:{port}",
            "ZT_DB_DRIVER": "sqlite",
            "ZT_DB_CONNSTRING": db_path,
            "ZT_MEDIA_PATH": media_path,
        }

        def startup_check(self):
            try:
                r = requests.get(f"http://127.0.0.1:{port}/ping", timeout=5)
                return r.status_code == 200
            except (requests.ConnectionError, requests.Timeout):
                return False

    logfile = xprocess.ensure(xprocess_name, Starter)
    yield
    xprocess.getinfo(xprocess_name).terminate()
