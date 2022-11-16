import json
from pathlib import Path
from typing import List, Optional, Union

import instagrapi.exceptions
import loguru
import requests
# noinspection PyUnresolvedReferences
from custom_logging import CustomizeLogger
# noinspection PyUnresolvedReferences
from dependencies import ClientStorage, get_clients
from fastapi import APIRouter, Depends, File, UploadFile, Form
# noinspection PyUnresolvedReferences
from helpers import photo_upload_post
from instagrapi.types import (
    Media, Location, Usertag
)
from pydantic import AnyHttpUrl, ValidationError
from starlette.responses import PlainTextResponse

config_path = Path(__file__).parent.with_name("logging_config.json")
logger: "loguru.Logger" = CustomizeLogger.make_logger(config_path)

router = APIRouter(
    prefix="/photo",
    tags=["photo"],
    responses={404: {"description": "Not found"}},
)


@router.post("/upload", response_model=Media)
async def photo_upload(sessionid: str = Form(...),
                       file: UploadFile = File(...),
                       caption: str = Form(...),
                       upload_id: Optional[str] = Form(""),
                       usertags: Optional[str] = Form(""),
                       cheap_proxy: str = Form(""),
                       # location: Optional[Location] = Form(None),
                       clients: ClientStorage = Depends(get_clients)
                       ) -> Union[PlainTextResponse, Media]:
    """Upload photo and configure to feed
    """
    try:
        cl = clients.get(sessionid)
    except instagrapi.exceptions.ChallengeRequired as ex:
        return PlainTextResponse(content="challenge required at start", status_code=400)

    usernames_tags = []

    if usertags is not None and usertags != "":
        try:
            usertags_json = json.loads(usertags)
            for usertag in usertags_json:
                usernames_tags.append(Usertag(user=usertag.get('user'), x=usertag.get('x'), y=usertag.get('y')))
        except ValidationError as ex:
            logger.exception(ex)
            return PlainTextResponse(content=ex, status_code=400)

    content = await file.read()
    return await photo_upload_post(
        cl, content, cheap_proxy, caption=caption,
        upload_id=upload_id,
        usertags=usernames_tags)
    # location=location)


@router.post("/upload/by_url", response_model=Media)
async def photo_upload(sessionid: str = Form(...),
                       url: AnyHttpUrl = Form(...),
                       caption: str = Form(...),
                       upload_id: Optional[str] = Form(""),
                       usertags: Optional[List[str]] = Form([]),
                       location: Optional[Location] = Form(None),
                       clients: ClientStorage = Depends(get_clients)
                       ) -> Media:
    """Upload photo and configure to feed
    """
    cl = clients.get(sessionid)

    usernames_tags = []
    for usertag in usertags:
        usertag_json = json.loads(usertag)
        usernames_tags.append(Usertag(user=usertag_json['user'], x=usertag_json['x'], y=usertag_json['y']))

    content = requests.get(url).content
    return await photo_upload_post(
        cl, content, caption=caption,
        upload_id=upload_id,
        usertags=usernames_tags,
        location=location)
