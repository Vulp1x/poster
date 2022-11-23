import time
from pathlib import Path
from typing import Optional, Dict, List, Union

import loguru
# noinspection PyUnresolvedReferences
from custom_logging import CustomizeLogger
# noinspection PyUnresolvedReferences
from dependencies import ClientStorage, get_clients
from fastapi import APIRouter, Depends, Body
from instagrapi import Client
from instagrapi.exceptions import ChallengeRequired, ChallengeError, LoginRequired, ClientError, UserNotFound
from pydantic import BaseModel
from starlette.responses import PlainTextResponse

router = APIRouter(
    prefix="/auth",
    tags=["auth"],
    responses={404: {"description": "Not found"}}
)

config_path = Path(__file__).parent.with_name("logging_config.json")
logger: "loguru.Logger" = CustomizeLogger.make_logger(config_path)

"""
{
  "session_id": "55421762746%3ALGcXzRrjhSjwE4%3A12%3AAYfqa6OVkj3BlaYZwks_WiewL5fWn-LIwmguOOnSsQ",
  "uuids": {
    "android_id": "android-ed0d3b157e361500", 
    "phone_id": "dc36ada3-c9ce-488c-89e0-48b564c8f060", 
    "uuid": "249160d7-3663-42ee-9e6e-c5d64eeb4ec4", 
    "advertising_id": "748a2e89-3fda-4c36-bf12-405d86557897"
  },
  "device_settings": {
    "app_version": "252.0.0.17.111",
    "android_version": 28,
    "android_release": "9.0.0",
    "dpi": "320dpi",
    "resolution": "720x1402",
    "manufacturer": "samsung",
    "device": "a10e",
    "model":  "SM-S102DL",
    "cpu": "exynos7885",
    "version_code": "397702079"
  },
  "user_agent": "Instagram 252.0.0.17.111 Android (28/9; 320dpi; 720x1402; samsung; SM-S102DL; a10e; exynos7885; en_IN; 397702079)",
  "proxy": "http://dmitrijkholodkov7815:21e49b@109.248.7.220:10475",
}
"""


class DeviceSettings(BaseModel):
    app_version: str
    android_version: int
    android_release: str
    dpi: str
    resolution: str
    manufacturer: str
    device: str
    model: str
    cpu: str
    version_code: str

    def as_dict(self) -> dict:
        return {
            "app_version": self.app_version,
            "android_version": self.android_version,
            "android_release": self.android_release,
            "dpi": self.dpi,
            "resolution": self.resolution,
            "manufacturer": self.manufacturer,
            "device": self.device,
            "model": self.model,
            "cpu": self.cpu,
            "version_code": self.version_code,
        }


class Uuids(BaseModel):
    android_id: str
    phone_id: str
    uuid: str
    advertising_id: str

    def as_dict(self) -> dict:
        return {"android_id": self.android_id, "phone_id": self.phone_id, "uuid": self.uuid,
                "advertising_id": self.advertising_id}


@router.post("/add")
async def auth_add(session_id: str = Body(...),
                   uuids: Uuids = Body(None),
                   device_settings: Optional[DeviceSettings] = Body(None),
                   user_agent: str = Body(...),
                   proxy: str = Body(""),
                   locale: Optional[str] = Body(""),
                   timezone: Optional[str] = Body(""),
                   clients: ClientStorage = Depends(get_clients)) -> PlainTextResponse:
    """Login by username and password with 2FA
    """
    try:
        cl = await clients.get(session_id)
        return PlainTextResponse(cl.sessionid)
    except (ChallengeError, ClientError) as ex:
        return PlainTextResponse(f"account is blocked: {ex}", status_code=400)

    except Exception as e:
        logger.warning(e)
        pass

    cl: Client = clients.client(proxy)

    if locale != "":
        cl.set_locale(locale)

    if timezone != "":
        cl.set_timezone_offset(timezone)

    cl.set_user_agent(user_agent)
    cl.set_device(device_settings.as_dict())
    cl.set_uuids(uuids.as_dict())

    try:
        result = cl.login_by_sessionid(session_id)
        if not result:
            return PlainTextResponse(result)
    except (UserNotFound, ClientError) as e:
        return PlainTextResponse(f"account is blocked: {e}", status_code=400)
    except AssertionError as e:
        if "sessionid" in str(e):
            return PlainTextResponse(f"invalid session id '{session_id}': {e}", status_code=400)

    await clients.set(cl)

    return PlainTextResponse(cl.sessionid)


@router.post("/follow_targets")
async def auth_add(session_id: str = Body(...),
                   target_user_ids: List[int] = Body(None),
                   clients: ClientStorage = Depends(get_clients)) -> PlainTextResponse:
    cl: Client = await clients.get(session_id)

    followers = cl.user_following(str(cl.user_id), use_cache=False, amount=0)
    followed_count = 0

    if len(followers) > 50:
        logger.info(f"already got {len(followers)} followers, skipping others")
        return PlainTextResponse(content=f'got {len(followers)} followings')

    for i, user_id in enumerate(target_user_ids):
        if followers.get(str(user_id), None) is not None:
            logger.info(f"bot is already a follower of {user_id}")
            followed_count += 1
            continue

        time.sleep(2)

        try:
            ok = cl.user_follow(str(user_id))
            logger.debug(f"user  {user_id}:  {ok}")

        except ChallengeRequired:
            return PlainTextResponse(status_code=400,
                                     content=f"after {i} followers got challenge required from user {user_id}")
        except Exception as e:
            logger.warning(f"got exception {e} when attempted to follow user {user_id}")
            continue

        if not ok:
            logger.warning(
                f"failed to follow user '{user_id}', followed {followed_count}/{i} users, skipping it")
            continue
        followed_count += 1

        if followed_count >= 25:
            break
        # time.sleep(2)

    return PlainTextResponse(content=f'got {followed_count} followings')


@router.get("/settings/get")
async def settings_get(sessionid: str,
                       clients: ClientStorage = Depends(get_clients)) -> Union[PlainTextResponse, Dict]:
    """Get client's settings
    """

    try:
        cl: Client = await clients.get(sessionid)
    except (ChallengeError, LoginRequired, ClientError) as e:
        return PlainTextResponse(status_code=400, content=f'bot is blocked: {e}')
    except IndexError as e:
        return PlainTextResponse(status_code=404, content=f'{e}')

    settings = cl.get_settings()
    return settings


#
# @router.post("/settings/set")
# async def settings_set(settings: str = Form(...),
#                        sessionid: Optional[str] = Form(""),
#                        clients: ClientStorage = Depends(get_clients)) -> str:
#     """Set client's settings
#     """
#     if sessionid != "":
#         cl = await clients.get(sessionid)
#     else:
#         cl = clients.client()
#     cl.set_settings(json.loads(settings))
#     cl.expose()
#     await clients.set(cl)
#     return cl.sessionid


@router.get("/timeline_feed")
async def timeline_feed(sessionid: str,
                        clients: ClientStorage = Depends(get_clients)) -> Union[PlainTextResponse, Dict]:
    """Get your timeline feed
    """
    cl: Client = await clients.get(sessionid)
    try:
        return cl.get_timeline_feed()

    except ChallengeRequired as e:
        return PlainTextResponse(status_code=400, content=f'bot {cl.username} is blocked')
