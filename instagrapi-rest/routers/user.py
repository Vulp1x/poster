import random
import tempfile
from pathlib import Path
from time import sleep
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
    User, UserShort, Media
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


@router.post("/similar/full", response_model=List[User])
async def similar_full(sessionid: str = Form(...),
                       username: str = Form(...),
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

    similar_bloggers: List[User] = []

    try:
        blogger = cl.user_info_by_username_v1(username)
        similar_bloggers.append(blogger)

    except Exception as e:
        return PlainTextResponse(status_code=404, content=f'blogger {username} not found: {e}')

    suggested_users = cl.private_request("discover/chaining/", params={"target_id": blogger.pk})
    extracted_users = [extract_user_short(user) for user in suggested_users.get('users', [])]

    logger.info(f'got {len(extracted_users)} similar accounts')

    for i, user in enumerate(extracted_users):
        if i > 30:
            break
        sleep(20 + random.randint(0, 15))

        try:
            user_info = cl.user_info(user.pk)
            similar_bloggers.append(user_info)
        except Exception as e:
            logger.warning(f"got exception {e} when tried to get info about user {user}")
            break

    logger.info(f'returning {len(similar_bloggers) - 1} similar accounts')

    return similar_bloggers


@router.post("/similar", response_model=List[UserShort])
async def similar(sessionid: str = Form(...),
                  username: str = Form(...),
                  clients: ClientStorage = Depends(get_clients)) -> Union[PlainTextResponse, List[UserShort]]:
    """Get user's followers
    """
    try:
        cl: Client = clients.get(sessionid)
    except instagrapi.exceptions.ChallengeRequired:
        return PlainTextResponse(status_code=400,
                                 content=f"required challenge on init")

    # all posts: 'https://i.instagram.com/api/v1/users/web_profile_info/?username={0}'

    # data = cl.private_request("users/403353154/info/?include_suggested_users=true")

    similar_bloggers: List[UserShort] = []
    blogger: Optional[UserShort] = None

    try:
        bloggers: List[UserShort] = cl.search_users_v1(username, 5)
        for blogger_ in bloggers:
            if blogger_.username == username:
                blogger = blogger_
                break

        if not blogger:
            blogger_full = cl.user_info_by_username_v1(username)
            blogger = extract_user_short(
                dict(pk=blogger_full.pk, username=blogger_full.username, full_name=blogger_full.full_name,
                     is_private=blogger_full.is_private, is_verified=blogger_full.is_verified))

        similar_bloggers.append(blogger)


    except Exception as e:
        return PlainTextResponse(status_code=404, content=f'blogger {username} not found: {e}')

    suggested_users = cl.private_request("discover/chaining/", params={"target_id": blogger.pk})
    similar_bloggers.extend([extract_user_short(user) for user in suggested_users.get('users', [])])

    logger.info(f'got {len(similar_bloggers)} similar accounts')

    return similar_bloggers


@router.post("/parse", response_model=List[UserShort])
async def parse_blogger(sessionid: str = Form(...),
                        user_id: int = Form(...),
                        posts_count: int = Form(...),
                        comments_count: int = Form(...),
                        likes_count: int = Form(...),
                        clients: ClientStorage = Depends(get_clients)) -> Union[PlainTextResponse, List[UserShort]]:
    """Get user's followers
    """
    try:
        cl: Client = clients.get(sessionid)
    except instagrapi.exceptions.ChallengeRequired:
        return PlainTextResponse(status_code=400,
                                 content=f"required challenge on init")

    # all posts: 'https://i.instagram.com/api/v1/users/web_profile_info/?username={0}'

    # data = cl.private_request("users/403353154/info/?include_suggested_users=true")

    parsed_users: List[UserShort] = []
    # blogger = cl.user_info_v1(str(user_id))
    # if not blogger.is_private:
    #     return PlainTextResponse(status_code=403, content=f"user {blogger} is private")
    #

    medias = cl.user_medias(user_id, posts_count)
    for media in medias:
        new_users = parse_post(cl, media, comments_count, likes_count)
        parsed_users.extend(new_users)

    logger.info(f'returning {len(parsed_users)} parsed users')

    return parsed_users


def parse_post(cl: Client, post: Media, commenters_to_parse: int, likers_to_parse) -> List[UserShort]:
    users_from_post: List[UserShort] = []

    log = logger.bind(media_id=post.pk)

    comments = cl.media_comments(post.pk, 3 * commenters_to_parse)
    comments = [comment for comment in comments if not comment.user.is_private]
    random_comments = random.sample(comments, commenters_to_parse)
    log.info(f'choose {len(random_comments)} from {len(comments)} comments')
    users_from_post.extend([comment.user for comment in random_comments])

    likes = cl.media_likers(post.pk)
    random_likes: List[UserShort] = random.sample(likes, likers_to_parse)
    log.info(f'choose {len(random_likes)} from {len(likes)} likes')
    users_from_post.extend(random_likes)

    users_from_post = list(set(users_from_post))
    log.info(f'got {len(users_from_post)} users from post')


    return users_from_post
