import logging
import tempfile
from pathlib import Path
from typing import List, Optional

from dependencies import ClientStorage, get_clients
from fastapi import APIRouter, Depends, Form, File, UploadFile
from instagrapi import Client
from instagrapi.exceptions import UserNotFound
from instagrapi.extractors import extract_account
from instagrapi.types import (
    User
)

router = APIRouter(
    prefix="/user",
    tags=["user"],
    responses={404: {"description": "Not found"}},
)

logger = logging.getLogger(__name__)


@router.post("/check/landings", response_model=List[str])
async def user_followers(sessionid: str = Form(...),
                         usernames: List[str] = Form(...),
                         clients: ClientStorage = Depends(get_clients)) -> List[str]:
    """Get user's followers
    """
    cl = clients.get(sessionid)
    checked_landing_accounts: List[str] = []
    if len(usernames) == 1:
        usernames = usernames[0].split(',')

    for username in usernames:
        try:
            user: User = cl.user_info_by_username_v1(username)
        except UserNotFound:
            logger.warning(f" checking landing account {username}: no client with this user name, skipping it")
            continue

        if not user.external_url:
            logger.warning(f" checking landing account {username}: got empty external link ")
            continue

        checked_landing_accounts.append(user.username)

    return checked_landing_accounts


@router.post("/edit_profile")
async def photo_upload(sessionid: str = Form(...),
                       file: Optional[UploadFile] = File(...),
                       full_name: Optional[str] = Form(...),
                       clients: ClientStorage = Depends(get_clients)
                       ):
    """Обновить фотографию профиля"""
    cl: Client = clients.get(sessionid)

    # if full_name:
    #     result = cl.private_request(
    #         "accounts/edit_profile/", cl.with_default_data({'first_name': full_name})
    #     )
    #
    #     print(extract_account(result["user"]))

        # cl.account_edit(data={'first_name': full_name, 'email': full_name+'@gmail.com'})

    if file:
        content = await file.read()
        with tempfile.NamedTemporaryFile(suffix='.jpg') as fp:
            fp.write(content)
            cl.account_change_picture(Path(fp.name))
