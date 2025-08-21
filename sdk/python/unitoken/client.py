import os
import subprocess
import sys
import json
import requests


class UniTokenClient:
    def __init__(self):
        self.server_url = None  # type: str | None
        self.setup_root_path()
        self.detect_running_url_from_file()
        if not self.server_url:
            self.start_service()
        if not self.server_url:
            raise Exception("Failed to start UniToken service.")

    def setup_root_path(self):
        if sys.platform == "win32":
            self.root_path = os.path.join(os.environ["LOCALAPPDATA"], "UniToken")
        else:
            self.root_path = os.path.join(os.environ["HOME"], ".local/share/uni-token")

        if not os.path.exists(self.root_path):
            os.makedirs(self.root_path)

    def detect_running_url_from_file(self):
        try:
            path = os.path.join(self.root_path, "service.json")

            print(f"Checking for service info at: {path}")

            if not os.path.exists(path):
                return

            with open(path, "r") as file:
                service_info = json.load(file)

            self.server_url = service_info.get("url")

            if not self.get("").json().get("__uni_token", None):
                self.server_url = None
        except Exception as _:
            self.server_url = None

    def start_service(self):
        if sys.platform == "win32":
            exec_path = os.path.join(self.root_path, "service.exe")
        else:
            exec_path = os.path.join(self.root_path, "service")

        if not os.path.isfile(exec_path):
            self.download_service(exec_path)

        subprocess.run(
            [exec_path, "sudo", "setup"],
            check=True,
        )

        self.detect_running_url_from_file()

    def download_service(self, exec_path):
        filename = {
            "linux": "service-linux-amd64",
            "darwin": "service-darwin-amd64",
            "win32": "service-windows-amd64.exe",
        }[sys.platform]

        url = f"https://uni-token.app/release/{filename}"

        print(f"Downloading service from {url} to {exec_path}")

        response = requests.get(url)
        if response.status_code != 200:
            raise Exception(f"Failed to download service: {response.status_code}")
        with open(exec_path, "wb") as file:
            file.write(response.content)
        os.chmod(exec_path, 0o755)

    def get(self, path, **kwargs):
        return requests.get(self.server_url + path, **kwargs)

    def post(self, path, **kwargs):
        return requests.post(self.server_url + path, **kwargs)
