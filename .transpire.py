from pathlib import Path

from transpire.resources import Deployment, Ingress, Service
from transpire.types import Image
from transpire.utils import get_image_tag

name = "jukebox"


def objects():
    dep = Deployment(
        name="jukebox",
        image=get_image_tag("jukebox"),
        ports=[8080],
    )

    dep.obj.spec.template.spec.containers[0].args = [
        "-mpdhost",
        "tv.ocf.berkeley.edu",
        "-host",
        "0.0.0.0",
    ]

    svc = Service(
        name="jukebox",
        selector=dep.get_selector(),
        port_on_pod=8080,
        port_on_svc=80,
    )

    ing = Ingress.from_svc(
        svc=svc,
        host="jukebox.ocf.berkeley.edu",
        path_prefix="/",
    )

    yield dep.build()
    yield svc.build()
    yield ing.build()


def images():
    yield Image(name="jukebox", path=Path("/"))
