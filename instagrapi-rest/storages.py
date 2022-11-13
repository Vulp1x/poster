import json
from pathlib import Path
from urllib import parse

from instagrapi import Client
from instagrapi.exceptions import ClientJSONDecodeError, LoginRequired
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
            cl = Client(settings=settings, proxy=settings['proxy'], logger=logger)
            cl.username = settings.get('username', 'username_not_set')
            cl.request_logger = logger

            try:
                cl.get_timeline_feed()
            except ClientJSONDecodeError as ex:
                logger.exception(ex)
            # except LoginRequired as ex:
            #     logger.exception(ex)
            #     if cl.login('_shalinicious_qhi', 'zwbm1q5a', relogin=True):
            #         logger.info('login succeeded')

                return cl

            return cl
        except IndexError as e:
            raise IndexError(f'Session not found (e.g. after reload process), please re-login {e}')

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
