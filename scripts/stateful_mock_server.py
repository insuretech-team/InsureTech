#!/usr/bin/env python3
"""Stateful OpenAPI-backed mock server for Postman and frontend testing."""

from __future__ import annotations

import argparse
import copy
import json
import re
import threading
import uuid
from collections import defaultdict
from dataclasses import dataclass
from datetime import datetime, timezone
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from typing import Any
from urllib.parse import parse_qs, urlparse


SUCCESS_STATUS_PRIORITY = ("200", "201", "202", "204")
JSON_CONTENT_TYPES = ("application/json", "application/*+json")
ACTION_STATUS_MAP = {
    "activate": "ACTIVE",
    "deactivate": "INACTIVE",
    "discontinue": "DISCONTINUED",
    "cancel": "CANCELLED",
    "verify": "VERIFIED",
    "confirm-payment": "PAID",
    "pay": "PAID",
    "review": "REVIEWED",
    "approve": "APPROVED",
    "reject": "REJECTED",
    "renew": "RENEWED",
    "issue": "ISSUED",
}


def now_iso() -> str:
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def clone(value: Any) -> Any:
    return copy.deepcopy(value)


def singularize(name: str) -> str:
    if name.endswith("ies") and len(name) > 3:
        return name[:-3] + "y"
    if name.endswith("s") and len(name) > 1:
        return name[:-1]
    return name


def deep_merge(base: Any, incoming: Any) -> Any:
    if isinstance(base, dict) and isinstance(incoming, dict):
        merged = clone(base)
        for key, value in incoming.items():
            if key in merged:
                merged[key] = deep_merge(merged[key], value)
            else:
                merged[key] = clone(value)
        return merged
    return clone(incoming)


def is_placeholder_value(value: Any) -> bool:
    if value is None:
        return True
    if value in ("", "string", {}, []):
        return True
    return isinstance(value, str) and value.startswith("example")


def normalize_json_body(raw: bytes) -> Any:
    if not raw:
        return {}
    text = raw.decode("utf-8").strip()
    if not text:
        return {}
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        return None


@dataclass
class RouteSpec:
    template: str
    regex: re.Pattern[str]
    methods: dict[str, dict[str, Any]]
    param_names: list[str]
    static_score: int


class MockEngine:
    def __init__(self, spec_path: Path) -> None:
        if not spec_path.exists():
            raise FileNotFoundError(f"OpenAPI JSON not found: {spec_path}")
        self.spec = json.loads(spec_path.read_text(encoding="utf-8"))
        self.routes = self._build_routes(self.spec.get("paths", {}))
        self.collection_to_item_param = self._build_collection_item_map()
        self.collection_state: dict[str, dict[str, dict[str, Any]]] = defaultdict(dict)
        self.global_index: dict[str, tuple[str, dict[str, Any]]] = {}
        self.counters: dict[str, int] = defaultdict(int)
        self.request_counter = 0
        self.lock = threading.RLock()

    def _segments(self, value: str) -> list[str]:
        return [segment for segment in value.strip("/").split("/") if segment]

    def _build_routes(self, paths: dict[str, Any]) -> list[RouteSpec]:
        routes: list[RouteSpec] = []
        for template, methods in paths.items():
            escaped = re.escape(template).replace(r"\{", "{").replace(r"\}", "}")
            pattern = re.sub(r"\{([^{}]+)\}", lambda m: f"(?P<{m.group(1)}>[^/]+)", escaped)
            routes.append(
                RouteSpec(
                    template=template,
                    regex=re.compile(f"^{pattern}$"),
                    methods={k.lower(): v for k, v in methods.items() if isinstance(v, dict)},
                    param_names=re.findall(r"\{([^{}]+)\}", template),
                    static_score=len(re.sub(r"\{[^{}]+\}", "", template)),
                )
            )
        routes.sort(key=lambda route: (-route.static_score, len(route.param_names), route.template))
        return routes

    def _build_collection_item_map(self) -> dict[str, str]:
        mapping: dict[str, str] = {}
        for route in self.routes:
            segments = self._segments(route.template)
            if not segments:
                continue
            last = segments[-1]
            if re.fullmatch(r"\{[^{}]+\}", last):
                collection = "/" + "/".join(segments[:-1]) if len(segments) > 1 else "/"
                mapping.setdefault(collection, last[1:-1])
        return mapping

    def match(self, path: str) -> tuple[RouteSpec | None, dict[str, str]]:
        for route in self.routes:
            matched = route.regex.match(path)
            if matched:
                return route, matched.groupdict()
        return None, {}

    def handle(self, method: str, path: str, query: dict[str, list[str]], body: Any) -> tuple[int, dict[str, str], Any]:
        if path == "/_mock/reset" and method == "POST":
            with self.lock:
                self.collection_state.clear()
                self.global_index.clear()
                self.counters.clear()
            return 200, {}, self._success_body({"reset": True, "message": "Mock state cleared"})

        if path == "/_mock/state" and method == "GET":
            with self.lock:
                state = {
                    "collections": {key: len(records) for key, records in self.collection_state.items()},
                    "total_records": sum(len(records) for records in self.collection_state.values()),
                }
            return 200, {}, self._success_body(state)

        route, path_params = self.match(path)
        if route is None:
            return 404, {}, self._error_body(404, "NOT_FOUND", "No mock route exists for this path")

        operation = route.methods.get(method.lower())
        if operation is None:
            return 405, {"Allow": ", ".join(sorted(route.methods.keys())).upper()}, self._error_body(
                405, "METHOD_NOT_ALLOWED", "This route exists but does not support the requested method"
            )

        if body is None:
            return 400, {}, self._error_body(400, "MALFORMED_REQUEST", "Request body is not valid JSON")

        with self.lock:
            return self._handle_operation(route, operation, method.upper(), path, path_params, query, body)

    def _handle_operation(
        self,
        route: RouteSpec,
        operation: dict[str, Any],
        method: str,
        path: str,
        path_params: dict[str, str],
        query: dict[str, list[str]],
        body: Any,
    ) -> tuple[int, dict[str, str], Any]:
        if method == "POST" and self._is_create_route(route):
            return self._create_resource(route, operation, path, path_params, body)
        if method == "GET" and self._is_collection_route(route):
            return self._list_resources(route, operation, path, query)
        if method == "GET":
            return self._get_resource(route, operation, path, path_params)
        if method in {"PATCH", "PUT"}:
            return self._update_resource(route, operation, path, path_params, body)
        if method == "DELETE":
            return self._delete_resource(route, operation, path, path_params)
        if method == "POST":
            return self._apply_action(route, operation, path, path_params, body)
        return self._success_from_example(operation)

    def _is_collection_route(self, route: RouteSpec) -> bool:
        segments = self._segments(route.template)
        if not route.param_names:
            return True
        if not segments:
            return False
        last = segments[-1]
        if ":" in last or re.fullmatch(r"\{[^{}]+\}", last):
            return False
        return last.endswith("s")

    def _is_create_route(self, route: RouteSpec) -> bool:
        return ":" not in route.template and self._is_collection_route(route)

    def _collection_template_for_route(self, route: RouteSpec) -> str:
        segments = self._segments(route.template)
        if not segments:
            return "/"
        last = segments[-1]
        if re.fullmatch(r"\{[^{}]+\}", last) or (":" in last and "{" in last):
            return "/" + "/".join(segments[:-1]) if len(segments) > 1 else "/"
        last_param_index = -1
        for index, segment in enumerate(segments):
            if "{" in segment:
                last_param_index = index
        if last_param_index >= 0:
            return "/" + "/".join(segments[:last_param_index]) if last_param_index > 0 else "/"
        return route.template

    def _collection_key_for_request(self, route: RouteSpec, path: str) -> str:
        if self._is_collection_route(route):
            return path
        actual = self._segments(path)
        template = self._segments(route.template)
        if not actual:
            return "/"
        last = template[-1]
        if re.fullmatch(r"\{[^{}]+\}", last) or (":" in last and "{" in last):
            return "/" + "/".join(actual[:-1]) if len(actual) > 1 else "/"
        last_param_index = -1
        for index, segment in enumerate(template):
            if "{" in segment:
                last_param_index = index
        if last_param_index >= 0:
            return "/" + "/".join(actual[:last_param_index]) if last_param_index > 0 else "/"
        return path

    def _resource_id_field(self, route: RouteSpec) -> str:
        collection = route.template if self._is_collection_route(route) else self._collection_template_for_route(route)
        param_name = self.collection_to_item_param.get(collection)
        if param_name:
            return param_name
        segments = self._segments(collection)
        last_static = next((segment for segment in reversed(segments) if "{" not in segment), "resource")
        return singularize(last_static.replace("-", "_")) + "_id"

    def _generate_id(self, field_name: str) -> str:
        prefix = singularize(field_name.removesuffix("_id"))
        self.counters[prefix] += 1
        return f"mock_{prefix}_{self.counters[prefix]:06d}"

    def _pick_success_response(self, operation: dict[str, Any]) -> tuple[str, dict[str, Any]]:
        responses = operation.get("responses", {})
        for status in SUCCESS_STATUS_PRIORITY:
            if status in responses:
                return status, responses[status]
        for status, response in responses.items():
            if status.startswith("2"):
                return status, response
        return "200", {}

    def _response_example(self, response: dict[str, Any]) -> Any | None:
        content = response.get("content", {})
        for content_type in JSON_CONTENT_TYPES:
            media = content.get(content_type)
            if media is None:
                continue
            if "example" in media:
                return clone(media["example"])
            examples = media.get("examples") or {}
            if examples:
                first = next(iter(examples.values()))
                if isinstance(first, dict) and "value" in first:
                    return clone(first["value"])
        return None

    def _pick_error_example(self, operation: dict[str, Any], status_code: int) -> Any | None:
        return self._response_example(operation.get("responses", {}).get(str(status_code), {}))

    def _with_fresh_meta(self, body: Any) -> Any:
        payload = clone(body)
        if not isinstance(payload, dict):
            return payload
        payload.setdefault("success", True)
        payload.setdefault("data", {})
        if payload.get("success") is True:
            payload["error"] = None
        meta = payload.get("meta")
        if not isinstance(meta, dict):
            meta = {}
            payload["meta"] = meta
        self.request_counter += 1
        meta["request_id"] = f"req_mock_{self.request_counter:06d}"
        meta["timestamp"] = now_iso()
        meta.setdefault("pagination", None)
        return payload

    def _success_body(self, data: Any) -> dict[str, Any]:
        self.request_counter += 1
        return {
            "success": True,
            "data": data,
            "error": None,
            "meta": {
                "request_id": f"req_mock_{self.request_counter:06d}",
                "timestamp": now_iso(),
                "pagination": None,
            },
        }

    def _error_body(self, status: int, code: str, message: str) -> dict[str, Any]:
        self.request_counter += 1
        return {
            "success": False,
            "data": None,
            "error": {
                "code": code,
                "message": message,
                "error_id": f"err_mock_{uuid.uuid4().hex[:12]}",
                "retryable": False,
                "http_status_code": status,
                "field_violations": [],
            },
            "meta": {
                "request_id": f"req_mock_{self.request_counter:06d}",
                "timestamp": now_iso(),
                "pagination": None,
            },
        }

    def _success_from_example(self, operation: dict[str, Any]) -> tuple[int, dict[str, str], Any]:
        status_text, response = self._pick_success_response(operation)
        status_code = int(status_text)
        if status_code == 204:
            return 204, {}, None
        example = self._response_example(response) or self._success_body({})
        return status_code, {}, self._with_fresh_meta(example)

    def _build_record_from_example(
        self,
        route: RouteSpec,
        operation: dict[str, Any],
        path_params: dict[str, str],
        body: Any,
    ) -> tuple[dict[str, Any], dict[str, Any], str]:
        _, response = self._pick_success_response(operation)
        example = self._with_fresh_meta(self._response_example(response) or self._success_body({}))
        data = example.get("data")
        if not isinstance(data, dict):
            data = {}
            example["data"] = data

        wrapper_key = ""
        nested_keys = [key for key, value in data.items() if isinstance(value, dict)]
        if len(nested_keys) == 1:
            wrapper_key = nested_keys[0]
            record = clone(data[wrapper_key])
        else:
            record = clone(data)
        if not isinstance(record, dict):
            record = {}

        if isinstance(body, dict):
            record = deep_merge(record, body)
        for key, value in path_params.items():
            record.setdefault(key, value)

        id_field = self._resource_id_field(route)
        if id_field not in record or is_placeholder_value(record.get(id_field)):
            record[id_field] = self._generate_id(id_field)
        if "id" not in record and id_field != "id":
            record["id"] = record[id_field]

        timestamp = now_iso()
        record.setdefault("created_at", timestamp)
        record["updated_at"] = timestamp
        return example, record, wrapper_key

    def _copy_identifier_fields(self, target: dict[str, Any], record: dict[str, Any]) -> None:
        for key, value in record.items():
            if key == "id" or key.endswith("_id"):
                target.setdefault(key, value)

    def _finalize_record_response(
        self,
        example: dict[str, Any],
        record: dict[str, Any],
        wrapper_key: str,
        list_mode: bool = False,
        list_key: str = "",
        items: list[dict[str, Any]] | None = None,
        page: int = 1,
        page_size: int = 20,
        total_items: int = 0,
    ) -> dict[str, Any]:
        payload = self._with_fresh_meta(example)
        data = payload.get("data")
        if not isinstance(data, dict):
            data = {}
            payload["data"] = data

        if list_mode:
            items = items or []
            if list_key:
                data[list_key] = items
            else:
                data["items"] = items
                list_key = "items"
            data["total_count"] = total_items
            data["next_page_token"] = None if page * page_size >= total_items else f"page_{page + 1}"
            payload["meta"]["pagination"] = {
                "page": page,
                "page_size": page_size,
                "total_pages": max(1, (total_items + page_size - 1) // page_size),
                "total_items": total_items,
                "has_next": page * page_size < total_items,
                "has_previous": page > 1,
            }
            return payload

        if wrapper_key:
            data[wrapper_key] = record
            self._copy_identifier_fields(data, record)
        else:
            payload["data"] = record
        return payload

    def _store_record(self, collection_key: str, record: dict[str, Any], id_field: str) -> None:
        record_id = str(record[id_field])
        stored = clone(record)
        self.collection_state[collection_key][record_id] = stored
        self.global_index[record_id] = (collection_key, stored)

    def _find_record(self, route: RouteSpec, path: str, path_params: dict[str, str]) -> tuple[str | None, dict[str, Any] | None]:
        collection_key = self._collection_key_for_request(route, path)
        for param_name in route.param_names:
            candidate_id = path_params.get(param_name)
            if not candidate_id:
                continue
            if candidate_id in self.collection_state.get(collection_key, {}):
                return collection_key, clone(self.collection_state[collection_key][candidate_id])
            indexed = self.global_index.get(candidate_id)
            if indexed:
                return indexed[0], clone(indexed[1])
        return collection_key, None

    def _remove_record(self, collection_key: str, record_id: str) -> None:
        self.collection_state.get(collection_key, {}).pop(record_id, None)
        self.global_index.pop(record_id, None)

    def _create_resource(
        self,
        route: RouteSpec,
        operation: dict[str, Any],
        path: str,
        path_params: dict[str, str],
        body: Any,
    ) -> tuple[int, dict[str, str], Any]:
        example, record, wrapper_key = self._build_record_from_example(route, operation, path_params, body)
        collection_key = self._collection_key_for_request(route, path)
        id_field = self._resource_id_field(route)
        self._store_record(collection_key, record, id_field)
        return 201, {"Location": f"{collection_key}/{record[id_field]}"}, self._finalize_record_response(example, record, wrapper_key)

    def _list_resources(
        self,
        route: RouteSpec,
        operation: dict[str, Any],
        path: str,
        query: dict[str, list[str]],
    ) -> tuple[int, dict[str, str], Any]:
        collection_key = self._collection_key_for_request(route, path)
        if collection_key not in self.collection_state:
            return self._success_from_example(operation)

        status_text, response = self._pick_success_response(operation)
        example = self._response_example(response) or self._success_body({})
        data = example.get("data")
        list_key = next((key for key, value in (data or {}).items() if isinstance(value, list)), "")
        records = list(self.collection_state.get(collection_key, {}).values())

        search_term = (query.get("search") or [""])[0].strip().lower()
        if search_term:
            records = [record for record in records if search_term in json.dumps(record, ensure_ascii=True).lower()]

        page = max(1, int((query.get("page") or ["1"])[0] or "1"))
        page_size = max(1, int((query.get("page_size") or ["20"])[0] or "20"))
        start = (page - 1) * page_size
        paged = records[start : start + page_size]
        payload = self._finalize_record_response(
            clone(example),
            {},
            "",
            list_mode=True,
            list_key=list_key,
            items=clone(paged),
            page=page,
            page_size=page_size,
            total_items=len(records),
        )
        return int(status_text), {}, payload

    def _inject_path_ids(self, data: Any, path_params: dict[str, str]) -> Any:
        payload = clone(data)
        if isinstance(payload, dict):
            for key, value in list(payload.items()):
                if isinstance(value, dict):
                    payload[key] = self._inject_path_ids(value, path_params)
            for param_name, param_value in path_params.items():
                payload.setdefault(param_name, param_value)
                payload.setdefault("id", param_value)
        return payload

    def _get_resource(
        self,
        route: RouteSpec,
        operation: dict[str, Any],
        path: str,
        path_params: dict[str, str],
    ) -> tuple[int, dict[str, str], Any]:
        collection_key, record = self._find_record(route, path, path_params)
        if record is not None:
            status_text, response = self._pick_success_response(operation)
            example = self._response_example(response) or self._success_body({})
            data = example.get("data")
            wrapper_key = next((key for key, value in (data or {}).items() if isinstance(value, dict)), "")
            return int(status_text), {}, self._finalize_record_response(example, record, wrapper_key)

        if collection_key and collection_key in self.collection_state:
            error_example = self._pick_error_example(operation, 404)
            return 404, {}, self._with_fresh_meta(error_example) if error_example else self._error_body(
                404, "NOT_FOUND", "The requested resource does not exist"
            )

        status_text, response = self._pick_success_response(operation)
        example = self._response_example(response) or self._success_body({})
        payload = self._with_fresh_meta(example)
        payload["data"] = self._inject_path_ids(payload.get("data"), path_params)
        return int(status_text), {}, payload

    def _update_resource(
        self,
        route: RouteSpec,
        operation: dict[str, Any],
        path: str,
        path_params: dict[str, str],
        body: Any,
    ) -> tuple[int, dict[str, str], Any]:
        collection_key, existing = self._find_record(route, path, path_params)
        if existing is None:
            collection_key = self._collection_key_for_request(route, path)
            existing = {key: value for key, value in path_params.items()}

        status_text, response = self._pick_success_response(operation)
        example = self._response_example(response) or self._success_body({})
        wrapper_key = next((key for key, value in (example.get("data") or {}).items() if isinstance(value, dict)), "")
        if wrapper_key:
            existing = deep_merge(clone((example.get("data") or {}).get(wrapper_key, {})), existing)
        if isinstance(body, dict):
            existing = deep_merge(existing, body)

        id_field = self._resource_id_field(route)
        existing.setdefault(id_field, path_params.get(id_field) or self._generate_id(id_field))
        existing.setdefault("id", existing[id_field])
        existing.setdefault("created_at", now_iso())
        existing["updated_at"] = now_iso()
        self._store_record(collection_key, existing, id_field)
        return int(status_text), {}, self._finalize_record_response(example, existing, wrapper_key)

    def _delete_resource(
        self,
        route: RouteSpec,
        operation: dict[str, Any],
        path: str,
        path_params: dict[str, str],
    ) -> tuple[int, dict[str, str], Any]:
        collection_key, record = self._find_record(route, path, path_params)
        if record is not None:
            record_id = str(next(iter(path_params.values())))
            self._remove_record(collection_key or "", record_id)

        if "204" in operation.get("responses", {}):
            return 204, {}, None
        if record is None and collection_key and collection_key in self.collection_state:
            error_example = self._pick_error_example(operation, 404)
            return 404, {}, self._with_fresh_meta(error_example) if error_example else self._error_body(
                404, "NOT_FOUND", "The requested resource does not exist"
            )
        return self._success_from_example(operation)

    def _apply_action(
        self,
        route: RouteSpec,
        operation: dict[str, Any],
        path: str,
        path_params: dict[str, str],
        body: Any,
    ) -> tuple[int, dict[str, str], Any]:
        collection_key, record = self._find_record(route, path, path_params)
        if record is None:
            collection_key = self._collection_key_for_request(route, path)
            record = {key: value for key, value in path_params.items()}
        if isinstance(body, dict):
            record = deep_merge(record, body)

        last_segment = self._segments(route.template)[-1] if self._segments(route.template) else ""
        action_name = last_segment.split(":", 1)[1] if ":" in last_segment else last_segment
        if action_name in ACTION_STATUS_MAP:
            record["status"] = ACTION_STATUS_MAP[action_name]

        id_field = self._resource_id_field(route)
        record.setdefault(id_field, path_params.get(id_field) or self._generate_id(id_field))
        record.setdefault("id", record[id_field])
        record.setdefault("created_at", now_iso())
        record["updated_at"] = now_iso()
        self._store_record(collection_key, record, id_field)

        status_text, response = self._pick_success_response(operation)
        if status_text == "204":
            return 204, {}, None
        example = self._response_example(response) or self._success_body({})
        wrapper_key = next((key for key, value in (example.get("data") or {}).items() if isinstance(value, dict)), "")
        return int(status_text), {}, self._finalize_record_response(example, record, wrapper_key)


class MockRequestHandler(BaseHTTPRequestHandler):
    server_version = "InsureTechStatefulMock/1.0"

    @property
    def engine(self) -> MockEngine:
        return self.server.engine  # type: ignore[attr-defined]

    def log_message(self, fmt: str, *args: Any) -> None:
        return

    def do_OPTIONS(self) -> None:
        self.send_response(204)
        self._write_common_headers()
        self.end_headers()

    def do_GET(self) -> None:
        self._dispatch("GET")

    def do_POST(self) -> None:
        self._dispatch("POST")

    def do_PATCH(self) -> None:
        self._dispatch("PATCH")

    def do_PUT(self) -> None:
        self._dispatch("PUT")

    def do_DELETE(self) -> None:
        self._dispatch("DELETE")

    def _dispatch(self, method: str) -> None:
        parsed = urlparse(self.path)
        content_length = int(self.headers.get("Content-Length", "0") or "0")
        raw_body = self.rfile.read(content_length) if content_length else b""
        status_code, headers, response_body = self.engine.handle(
            method,
            parsed.path,
            parse_qs(parsed.query),
            normalize_json_body(raw_body),
        )

        self.send_response(status_code)
        self._write_common_headers()
        for key, value in headers.items():
            self.send_header(key, value)
        if response_body is None:
            self.end_headers()
            return

        payload = json.dumps(response_body, ensure_ascii=True).encode("utf-8")
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(payload)))
        self.end_headers()
        self.wfile.write(payload)

    def _write_common_headers(self) -> None:
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
        self.send_header("Cache-Control", "no-store")
        self.send_header("X-Mock-Server", "insuretech-stateful")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Start the InsureTech stateful mock server")
    parser.add_argument("--port", type=int, default=4010, help="Port to listen on (default: 4010)")
    parser.add_argument("--spec", default="api/docs/openapi.json", help="Path to generated OpenAPI JSON")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    project_root = Path(__file__).resolve().parent.parent
    engine = MockEngine((project_root / args.spec).resolve())
    server = ThreadingHTTPServer(("0.0.0.0", args.port), MockRequestHandler)
    server.engine = engine  # type: ignore[attr-defined]

    print("")
    print("============================================================")
    print("  InsureTech Stateful Mock Server")
    print("============================================================")
    print(f"  Base URL:         http://localhost:{args.port}")
    print(f"  Reset endpoint:   POST http://localhost:{args.port}/_mock/reset")
    print(f"  State endpoint:   GET  http://localhost:{args.port}/_mock/state")
    print("============================================================")
    print("  Press Ctrl+C to stop")
    print("============================================================")
    print("")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        pass
    finally:
        server.server_close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
