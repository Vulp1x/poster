import uuid
from pathlib import Path

import pkg_resources
import uvicorn
from fastapi import FastAPI, Depends
from fastapi.openapi.utils import get_openapi
from starlette.requests import Request
from starlette.responses import RedirectResponse, JSONResponse
from starlette.types import Message

from custom_logging import CustomizeLogger
from routers import (
    auth, media, video, photo, user,
    igtv, clip, album, story,
    insights
)


# sys.tracebacklimit = 6


async def set_body(request: Request, body: bytes):
    async def receive() -> Message:
        return {"type": "http.request", "body": body}

    request._receive = receive


async def get_body(request: Request) -> bytes:
    body = await request.body()
    await set_body(request, body)

    return body


async def logging_dependency(request: Request):
    logger.debug(f"{request.method} {request.url} started, {await request.form()}")
    # logger.debug(f"{request.method} {request.url} body: {await get_body(request)}")


config_path = Path(__file__).with_name("logging_config.json")

app = FastAPI()
logger = CustomizeLogger.make_logger(config_path)
app.logger = logger
app.include_router(auth.router, dependencies=[Depends(logging_dependency)])

app.include_router(media.router, dependencies=[Depends(logging_dependency)])

app.include_router(video.router, dependencies=[Depends(logging_dependency)])

app.include_router(photo.router, dependencies=[Depends(logging_dependency)])

app.include_router(user.router, dependencies=[Depends(logging_dependency)])

app.include_router(igtv.router, dependencies=[Depends(logging_dependency)])

app.include_router(clip.router, dependencies=[Depends(logging_dependency)])

app.include_router(album.router, dependencies=[Depends(logging_dependency)])

app.include_router(story.router, dependencies=[Depends(logging_dependency)])

app.include_router(insights.router, dependencies=[Depends(logging_dependency)])


@app.middleware("http")
async def request_middleware(request: Request, call_next):
    request_id = request.headers.get("X-Request-ID")
    if request_id is None:
        request_id = str(uuid.uuid4())

    with logger.contextualize(request_id=request_id):
        try:
            response = await call_next(request)
            response.headers["X-Request-ID"] = request_id
            return response

        except Exception as ex:
            logger.exception(f"Request failed: {ex}")
            response = JSONResponse({
                "detail": str(ex),
                "exc_type": str(type(ex).__name__)
            }, status_code=500)

            response.headers["X-Request-ID"] = request_id
            return response


@app.get("/", tags=["system"], summary="Redirect to /docs")
async def root():
    """Redirect to /docs
    """
    return RedirectResponse(url="/docs")


@app.get("/version", tags=["system"], summary="Get dependency versions")
async def version():
    """Get dependency versions
    """
    versions = {}
    for name in ('instagrapi',):
        item = pkg_resources.require(name)
        if item:
            versions[name] = item[0].version
    return versions


# @app.exception_handler(Exception)
# async def handle_exception(request, exc: Exception):
#     return JSONResponse({
#         "detail": str(exc),
#         "exc_type": str(type(exc).__name__)
#     }, status_code=500)


def custom_openapi():
    if app.openapi_schema:
        return app.openapi_schema
    # for route in app.routes:
    #     body_field = getattr(route, 'body_field', None)
    #     if body_field:
    #         body_field.type_.__name__ = 'name'
    openapi_schema = get_openapi(
        title="instagrapi-rest",
        version="1.0.0",
        description="RESTful API Service for instagrapi",
        routes=app.routes,
    )
    app.openapi_schema = openapi_schema
    return app.openapi_schema


app.openapi = custom_openapi

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, workers=4)
