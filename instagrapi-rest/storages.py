import json
from pathlib import Path
from urllib import parse

from instagrapi import Client
from tinydb import TinyDB, Query

from custom_logging import CustomizeLogger

config_path = Path(__file__).with_name("logging_config.json")
logger = CustomizeLogger.make_logger(config_path)


class ClientStorage:
    db = TinyDB('./db.json')

    def client(self, proxy: str):
        """Get new client (helper)
        """
        cl = Client(proxy=proxy)
        cl.request_timeout = 0.1
        return cl

    def get(self, sessionid: str) -> Client:
        """Get client settings
        """
        key = parse.unquote(sessionid.strip(" \""))
        try:
            settings = json.loads(self.db.search(Query().sessionid == key)[0]['settings'])
            cl = Client(settings=settings, proxy=settings['proxy'])
            cl.username = settings.get('username', 'username_not_set')
            cl.request_logger = logger
            cl.get_timeline_feed()
            return cl
        except IndexError:
            raise Exception('Session not found (e.g. after reload process), please relogin')

    def set(self, cl: Client) -> bool:
        """Set client settings
        """
        key = parse.unquote(cl.sessionid.strip(" \""))
        client_settings = cl.get_settings()
        client_settings['proxy'] = cl.proxy
        client_settings['username'] = cl.username
        self.db.insert({'sessionid': key, 'settings': json.dumps(client_settings)})

        return True

    def close(self):
        pass
