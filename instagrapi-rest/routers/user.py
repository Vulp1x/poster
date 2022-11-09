import tempfile
from pathlib import Path
from typing import List, Optional, Union

import instagrapi.exceptions
import loguru
# noinspection PyUnresolvedReferences
from custom_logging import CustomizeLogger
# noinspection PyUnresolvedReferences
from dependencies import ClientStorage, get_clients
from fastapi import APIRouter, Depends, Form, File, UploadFile
from instagrapi import Client
from instagrapi.exceptions import UserNotFound
from instagrapi.extractors import extract_user_short
from instagrapi.types import (
    User
)
from starlette.responses import JSONResponse, PlainTextResponse

router = APIRouter(
    prefix="/user",
    tags=["user"],
    responses={404: {"description": "Not found"}},
)

config_path = Path(__file__).parent.with_name("logging_config.json")
logger: "loguru.Logger" = CustomizeLogger.make_logger(config_path)


@router.post("/check/landings", response_model=List[str])
async def check_landings(sessionid: str = Form(...),
                         usernames: List[str] = Form(...),
                         clients: ClientStorage = Depends(get_clients)) -> JSONResponse:
    """Get user's followers
    """
    try:
        cl: Client = clients.get(sessionid)
    except instagrapi.exceptions.ChallengeRequired:
        return JSONResponse(status_code=400,
                            content=f"required challenge on init")

    checked_landing_accounts: List[str] = []
    if len(usernames) == 1:
        usernames = usernames[0].split(',')

    # log_: loguru.Logger = logging.getLogger("private_request")
    with logger.contextualize(user_id=cl.user_id):
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

    return JSONResponse(checked_landing_accounts)


@router.post("/edit_profile")
async def edit_profile(sessionid: str = Form(...),
                       file: Optional[UploadFile] = File(...),
                       full_name: Optional[str] = Form(...),
                       clients: ClientStorage = Depends(get_clients)
                       ):
    """Обновить фотографию профиля"""
    try:
        cl: Client = clients.get(sessionid)
    except instagrapi.exceptions.ChallengeRequired:
        return PlainTextResponse(status_code=400,
                                 content=f"required challenge on init")

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


@router.post("/similar", response_model=List[User])
async def similar(sessionid: str = Form(...),
                  user_id: int = Form(...),
                  clients: ClientStorage = Depends(get_clients)) -> Union[PlainTextResponse, List[User]]:
    """Get user's followers
    """
    try:
        cl: Client = clients.get(sessionid)
    except instagrapi.exceptions.ChallengeRequired:
        return PlainTextResponse(status_code=400,
                                 content=f"required challenge on init")

    # all posts: 'https://i.instagram.com/api/v1/users/web_profile_info/?username={0}'

    # data = cl.private_request("users/403353154/info/?include_suggested_users=true")

    suggested_users = cl.private_request("discover/chaining/", params={"target_id": user_id})
    extracted_users = [extract_user_short(user) for user in suggested_users.get('users', [])]

    logger.info(f'got {len(extracted_users)} similar accounts')

    similar_bloggers: List[User] = []

    for i, user in enumerate(extracted_users):
        if i > 1:
            break

        try:
            user_info = cl.user_info(user.pk)
            similar_bloggers.append(user_info)
        except Exception as e:
            logger.warning(f"got exception {e} when tried to get info about user {user}")

    return similar_bloggers
