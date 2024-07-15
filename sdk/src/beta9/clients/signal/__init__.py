# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: signal.proto
# plugin: python-betterproto
# This file has been @generated

from dataclasses import dataclass
from typing import (
    TYPE_CHECKING,
    AsyncIterator,
    Dict,
    Iterator,
    Optional,
)

import betterproto
import grpc
from betterproto.grpcstub.grpcio_client import SyncServiceStub
from betterproto.grpcstub.grpclib_server import ServiceBase


if TYPE_CHECKING:
    import grpclib.server
    from betterproto.grpcstub.grpclib_client import MetadataLike
    from grpclib.metadata import Deadline


@dataclass(eq=False, repr=False)
class SignalSetRequest(betterproto.Message):
    name: str = betterproto.string_field(1)
    ttl: int = betterproto.int64_field(2)


@dataclass(eq=False, repr=False)
class SignalSetResponse(betterproto.Message):
    ok: bool = betterproto.bool_field(1)


@dataclass(eq=False, repr=False)
class SignalClearRequest(betterproto.Message):
    name: str = betterproto.string_field(1)


@dataclass(eq=False, repr=False)
class SignalClearResponse(betterproto.Message):
    ok: bool = betterproto.bool_field(1)


@dataclass(eq=False, repr=False)
class SignalMonitorRequest(betterproto.Message):
    name: str = betterproto.string_field(1)


@dataclass(eq=False, repr=False)
class SignalMonitorResponse(betterproto.Message):
    ok: bool = betterproto.bool_field(1)
    set: bool = betterproto.bool_field(2)


class SignalServiceStub(SyncServiceStub):
    def signal_set(self, signal_set_request: "SignalSetRequest") -> "SignalSetResponse":
        return self._unary_unary(
            "/signal.SignalService/SignalSet",
            SignalSetRequest,
            SignalSetResponse,
        )(signal_set_request)

    def signal_clear(
        self, signal_clear_request: "SignalClearRequest"
    ) -> "SignalClearResponse":
        return self._unary_unary(
            "/signal.SignalService/SignalClear",
            SignalClearRequest,
            SignalClearResponse,
        )(signal_clear_request)

    def signal_monitor(
        self, signal_monitor_request: "SignalMonitorRequest"
    ) -> Iterator["SignalMonitorResponse"]:
        for response in self._unary_stream(
            "/signal.SignalService/SignalMonitor",
            SignalMonitorRequest,
            SignalMonitorResponse,
        )(signal_monitor_request):
            yield response