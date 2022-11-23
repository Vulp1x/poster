from typing import Generator

from storages import ClientStorage


async def get_clients() -> Generator:
    try:
        clients = ClientStorage()
        await clients.init()
        yield clients
    finally:
        clients.close()
