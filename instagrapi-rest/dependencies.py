from typing import Generator

from storages import ClientStorage


async def get_clients() -> Generator:
    try:
        clients = ClientStorage()
        yield clients
    finally:
        clients.close()
