# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: volume.proto
# plugin: python-betterproto
from dataclasses import dataclass

import betterproto
import grpclib


@dataclass
class GetOrCreateVolumeRequest(betterproto.Message):
    name: str = betterproto.string_field(1)


@dataclass
class GetOrCreateVolumeResponse(betterproto.Message):
    ok: bool = betterproto.bool_field(1)
    volume_id: str = betterproto.string_field(2)


class VolumeServiceStub(betterproto.ServiceStub):
    async def get_or_create_volume(
        self, *, name: str = ""
    ) -> GetOrCreateVolumeResponse:
        request = GetOrCreateVolumeRequest()
        request.name = name

        return await self._unary_unary(
            "/volume.VolumeService/GetOrCreateVolume",
            request,
            GetOrCreateVolumeResponse,
        )