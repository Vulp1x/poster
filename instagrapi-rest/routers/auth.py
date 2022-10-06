import json
from typing import Optional, Dict
from fastapi import APIRouter, Depends, Form, Body
from dependencies import ClientStorage, get_clients
from pydantic import BaseModel

router = APIRouter(
    prefix="/auth",
    tags=["auth"],
    responses={404: {"description": "Not found"}}
)


@router.post("/login")
async def auth_login(username: str = Form(...),
                     password: str = Form(...),
                     verification_code: Optional[str] = Form(""),
                     proxy: Optional[str] = Form(""),
                     locale: Optional[str] = Form(""),
                     timezone: Optional[str] = Form(""),
                     clients: ClientStorage = Depends(get_clients)) -> str:
    """Login by username and password with 2FA
    """
    cl = clients.client()
    if proxy != "":
        cl.set_proxy(proxy)

    if locale != "":
        cl.set_locale(locale)

    if timezone != "":
        cl.set_timezone_offset(timezone)

    result = cl.login(
        username,
        password,
        verification_code=verification_code
    )
    if result:
        clients.set(cl)
        return cl.sessionid
    return result


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
                   proxy: Optional[str] = Body(""),
                   locale: Optional[str] = Body(""),
                   timezone: Optional[str] = Body(""),
                   clients: ClientStorage = Depends(get_clients)) -> str:
    """Login by username and password with 2FA
    """
    cl = clients.client()
    if proxy != "":
        cl.set_proxy(proxy)

    if locale != "":
        cl.set_locale(locale)

    if timezone != "":
        cl.set_timezone_offset(timezone)

    cl.set_user_agent(user_agent)
    cl.set_device(device_settings.as_dict())
    cl.set_uuids(uuids.as_dict())

    result = cl.login_by_sessionid(session_id)
    if result:
        clients.set(cl)
        return cl.sessionid
    return result


@router.post("/relogin")
async def auth_relogin(sessionid: str = Form(...),
                       clients: ClientStorage = Depends(get_clients)) -> str:
    """Relogin by username and password (with clean cookies)
    """
    cl = clients.get(sessionid)
    result = cl.relogin()
    return result


@router.get("/settings/get")
async def settings_get(sessionid: str,
                       clients: ClientStorage = Depends(get_clients)) -> Dict:
    """Get client's settings
    """
    cl = clients.get(sessionid)
    return cl.get_settings()


@router.post("/settings/set")
async def settings_set(settings: str = Form(...),
                       sessionid: Optional[str] = Form(""),
                       clients: ClientStorage = Depends(get_clients)) -> str:
    """Set client's settings
    """
    if sessionid != "":
        cl = clients.get(sessionid)
    else:
        cl = clients.client()
    cl.set_settings(json.loads(settings))
    cl.expose()
    clients.set(cl)
    return cl.sessionid


@router.get("/timeline_feed")
async def timeline_feed(sessionid: str,
                        clients: ClientStorage = Depends(get_clients)) -> Dict:
    """Get your timeline feed
    """
    cl = clients.get(sessionid)
    return cl.get_timeline_feed()
