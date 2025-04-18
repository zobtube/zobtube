import pytest
import requests
from xprocess import ProcessStarter

@pytest.fixture(autouse=True)
def zotbue_server(xprocess):
    class Starter(ProcessStarter):
        # startup pattern
        #pattern = "PATTERN"

        # command to start process
        args = ['/tmp/zt']

        env = {
            "ZT_SERVER_BIND": '127.0.0.1:8080',
            "ZT_DB_DRIVER": 'sqlite',
            "ZT_DB_CONNSTRING": '/tmp/zt-db.sqlite3',
            "ZT_MEDIA_PATH": '/tmp/zt-data',
        }

        def startup_check(self):
            r = requests.get('http://127.0.0.1:8080/ping')
            return r.status_code == 200

    # ensure process is running and return its logfile
    logfile = xprocess.ensure("zobtube", Starter)
    yield
    xprocess.getinfo("zobtube").terminate()
