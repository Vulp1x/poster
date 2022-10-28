# Custom Logger Using Loguru

import json
import logging
import sys
from pathlib import Path

import loguru
from loguru import logger


class InterceptHandler(logging.Handler):
    loglevel_mapping = {
        50: 'CRITICAL',
        40: 'ERROR',
        30: 'WARNING',
        20: 'INFO',
        10: 'DEBUG',
        0: 'NOTSET',
    }

    def emit(self, record):
        try:
            level = logger.level(record.levelname).name
        except AttributeError:
            level = self.loglevel_mapping[record.levelno]

        # frame, depth = logging.currentframe(), 2
        # while frame.f_code.co_filename == logging.__file__:
        #     frame = frame.f_back
        #     depth += 1

        log = logger.bind(request_id='app')
        log.opt(
            depth=4,
            exception=record.exc_info
        ).log(level, record.getMessage())


class CustomizeLogger:

    @classmethod
    def make_logger(cls, config_path: Path) -> "loguru.Logger":
        config = cls.load_logging_config(config_path)
        logging_config = config.get('logger')

        logger = cls.customize_logging(
            logging_config.get('path'),
            level=logging_config.get('level'),
            retention=logging_config.get('retention'),
            rotation=logging_config.get('rotation'),
            format=logging_config.get('format')
        )

        return logger

    @classmethod
    def customize_logging(cls,
                          filepath: Path,
                          level: str,
                          rotation: str,
                          retention: str,
                          format: str
                          ):
        logger.remove()

        def not_too_long(record: "loguru.Record") -> bool:
            return len(record.get('message')) < 1000

        logger.add(
            sys.stdout,
            enqueue=True,
            backtrace=False,
            level='TRACE',
            format=format, filter=not_too_long,
        )
        logger.add(
            str(filepath),
            rotation=rotation,
            retention=retention,
            enqueue=True,
            backtrace=True,
            level=level.upper(),
            format=format
        )
        logging.basicConfig(handlers=[InterceptHandler()], level=10, force=True)
        logging.getLogger("uvicorn.access").handlers = [InterceptHandler()]

        # _logger = logging.getLogger('public_request')
        # _logger.setLevel('ERROR')
        _logger = logging.getLogger('public_request')
        _logger.handlers = [InterceptHandler()]

        for _log in ['uvicorn', 'uvicorn.error']:
            _logger = logging.getLogger(_log)
            _logger.handlers = [InterceptHandler()]

        return logger.bind(method=None, user_id=None)

    @classmethod
    def load_logging_config(cls, config_path):
        config = None
        with open(config_path) as config_file:
            config = json.load(config_file)
        return config
