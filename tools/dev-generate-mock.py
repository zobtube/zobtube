#!/usr/bin/env python3
"""Generate mock data for a local dev ZobTube instance."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request
from http.cookiejar import CookieJar
from typing import Any

# Categories: (category_name, sub_category_name)
CATEGORIES: list[tuple[str, str]] = [
    ("Countries", "Russia"),
    ("Countries", "United Kingdom"),
    ("Countries", "Estonia"),
    ("Countries", "Latvia"),
    ("Countries", "Netherlands"),
    ("Countries", "Serbia"),
    ("Countries", "Hungary"),
    ("Countries", "Czech Republic"),
    ("Countries", "Ukraine"),
    ("Countries", "France"),
    ("Countries", "Italy"),
]

# Actors: (name, sex, country sub-category name)
ACTORS: list[tuple[str, str, str]] = [
    ("Kwini Kim", "f", "Russia"),
    ("Skye Young", "f", "Russia"),
    ("Yasmina Khan", "f", "United Kingdom"),
    ("Lia Lin", "f", "Russia"),
    ("Hot Pearl", "f", "Russia"),
    ("Sweetie Fox", "f", "Russia"),
    ("Lily Phillips", "f", "United Kingdom"),
    ("Bonnie Blue", "f", "United Kingdom"),
    ("Katty West", "f", "Russia"),
    ("CarlaCute", "f", "Russia"),
    ("Octokuro", "f", "Russia"),
    ("Shinaryen", "f", "Estonia"),
    ("Matty", "f", "Latvia"),
    ("Taylor Sands", "f", "Netherlands"),
    ("Hello Siri3", "f", "Russia"),
    ("Maya Stone", "f", "Russia"),
    ("Foxy Di", "f", "Russia"),
    ("Alexis Crystal", "f", "Czech Republic"),
    ("Stacy Cruz", "f", "Czech Republic"),
    ("Britney Dutch", "f", "Netherlands"),
    ("Yoya Grey", "f", "Serbia"),
    ("Vale Nappi", "f", "Italy"),
    ("Jadilica", "f", "Russia"),
    ("Tiffany Tatum", "f", "Hungary"),
    ("kiraxxcherry", "f", "Ukraine"),
    ("Anissa Kate", "f", "France"),
    ("Stella Cox", "f", "Italy"),
    ("Milka", "f", "Russia"),
    ("Frances Bentley", "f", "United Kingdom"),
    ("Mila Pie", "f", "Russia"),
]

COOKIE_NAME = "zt_auth"


def field(obj: dict[str, Any], *keys: str) -> str:
    for key in keys:
        value = obj.get(key)
        if value is not None and value != "":
            return str(value)
    return ""


def sha256_hex(value: str) -> str:
    return hashlib.sha256(value.encode()).hexdigest()


class ZobTubeClient:
    def __init__(self, base_url: str, *, verbose: bool = False) -> None:
        self.base_url = base_url.rstrip("/")
        self.verbose = verbose
        self._jar = CookieJar()
        self._opener = urllib.request.build_opener(
            urllib.request.HTTPCookieProcessor(self._jar)
        )

    def request(
        self,
        method: str,
        path: str,
        *,
        data: dict[str, str] | None = None,
    ) -> tuple[int, dict[str, Any] | None, str]:
        url = self.base_url + path
        body: bytes | None = None
        headers: dict[str, str] = {}
        if data is not None:
            body = urllib.parse.urlencode(data).encode()
            headers["Content-Type"] = "application/x-www-form-urlencoded"

        req = urllib.request.Request(url, data=body, headers=headers, method=method)
        if self.verbose:
            print(f">>> {method} {url}", file=sys.stderr)
            if data:
                print(f">>> body: {data}", file=sys.stderr)

        try:
            with self._opener.open(req) as resp:
                raw = resp.read().decode()
                status = resp.status
        except urllib.error.HTTPError as exc:
            raw = exc.read().decode()
            status = exc.code

        parsed: dict[str, Any] | None = None
        if raw:
            try:
                parsed = json.loads(raw)
            except json.JSONDecodeError:
                pass

        if self.verbose:
            print(f"<<< {status} {raw}", file=sys.stderr)

        return status, parsed, raw

    def bootstrap(self) -> dict[str, Any]:
        status, body, raw = self.request("GET", "/api/bootstrap")
        if status != 200 or not body:
            raise RuntimeError(f"bootstrap failed ({status}): {raw}")
        return body

    def login(self, username: str, password: str) -> None:
        session_id = None
        for cookie in self._jar:
            if cookie.name == COOKIE_NAME:
                session_id = cookie.value
                break
        if not session_id:
            raise RuntimeError(f"missing {COOKIE_NAME} cookie after bootstrap")

        challenge = sha256_hex(session_id + sha256_hex(password))
        status, body, raw = self.request(
            "POST",
            "/auth/login",
            data={"username": username, "password": challenge},
        )
        if status != 200:
            detail = (body or {}).get("error", raw)
            raise RuntimeError(f"login failed ({status}): {detail}")

    def list_actors(self) -> list[dict[str, Any]]:
        status, body, raw = self.request("GET", "/api/actor")
        if status != 200 or not body:
            raise RuntimeError(f"list actors failed ({status}): {raw}")
        items = body.get("items")
        if not isinstance(items, list):
            raise RuntimeError(f"list actors: unexpected response: {raw}")
        return items

    @staticmethod
    def actor_index(items: list[dict[str, Any]]) -> dict[str, str]:
        """Map actor name -> id (first match when duplicates exist)."""
        index: dict[str, str] = {}
        for actor in items:
            name = field(actor, "Name", "name")
            actor_id = field(actor, "ID", "id")
            if name and actor_id and name not in index:
                index[name] = actor_id
        return index

    def create_actor(self, name: str, sex: str) -> str:
        status, body, raw = self.request(
            "POST",
            "/api/actor/",
            data={"name": name, "sex": sex},
        )
        if status != 200 or not body:
            raise RuntimeError(f"create actor {name!r} failed ({status}): {raw}")
        actor_id = body.get("result")
        if not actor_id:
            raise RuntimeError(f"create actor {name!r}: missing result id: {raw}")
        return str(actor_id)

    def ensure_actor(
        self, name: str, sex: str, index: dict[str, str]
    ) -> tuple[str, bool]:
        """Return (actor_id, created). Updates index in place."""
        if name in index:
            return index[name], False
        actor_id = self.create_actor(name, sex)
        index[name] = actor_id
        return actor_id, True

    def list_categories(self) -> list[dict[str, Any]]:
        status, body, raw = self.request("GET", "/api/category")
        if status != 200 or not body:
            raise RuntimeError(f"list categories failed ({status}): {raw}")
        items = body.get("items")
        if not isinstance(items, list):
            raise RuntimeError(f"list categories: unexpected response: {raw}")
        return items

    def create_category(self, name: str) -> None:
        status, body, raw = self.request(
            "POST",
            "/api/category",
            data={"Name": name},
        )
        if status != 200:
            detail = (body or {}).get("error", raw)
            raise RuntimeError(f"create category {name!r} failed ({status}): {detail}")

    def create_category_sub(self, name: str, parent_id: str) -> None:
        status, body, raw = self.request(
            "POST",
            "/api/category-sub",
            data={"Name": name, "Parent": parent_id},
        )
        if status != 200:
            detail = (body or {}).get("error", raw)
            raise RuntimeError(
                f"create sub-category {name!r} failed ({status}): {detail}"
            )

    @staticmethod
    def category_index(
        items: list[dict[str, Any]],
    ) -> dict[str, dict[str, Any]]:
        """Map parent category name -> {id, subs: {sub_name: sub_id}}."""
        index: dict[str, dict[str, Any]] = {}
        for cat in items:
            parent_name = field(cat, "Name", "name")
            parent_id = field(cat, "ID", "id")
            subs: dict[str, str] = {}
            for sub in cat.get("Sub") or cat.get("sub") or []:
                if not isinstance(sub, dict):
                    continue
                sub_name = field(sub, "Name", "name")
                sub_id = field(sub, "ID", "id")
                if sub_name and sub_id:
                    subs[sub_name] = sub_id
            if parent_name:
                index[parent_name] = {"id": parent_id, "subs": subs}
        return index

    def ensure_categories(
        self, categories: list[tuple[str, str]]
    ) -> dict[str, str]:
        """Create missing parent categories and sub-categories; return sub_name -> id."""
        wanted: dict[str, set[str]] = {}
        for parent, sub in categories:
            wanted.setdefault(parent, set()).add(sub)

        index = self.category_index(self.list_categories())
        for parent in wanted:
            if parent not in index:
                self.create_category(parent)
                index = self.category_index(self.list_categories())

        for parent, sub_names in wanted.items():
            parent_id = index[parent]["id"]
            if not parent_id:
                raise RuntimeError(f"category {parent!r} has no id")
            existing = index[parent]["subs"]
            for sub_name in sorted(sub_names):
                if sub_name not in existing:
                    self.create_category_sub(sub_name, parent_id)

        index = self.category_index(self.list_categories())
        sub_ids: dict[str, str] = {}
        for parent, sub_names in wanted.items():
            existing = index[parent]["subs"]
            for sub_name in sub_names:
                sub_id = existing.get(sub_name)
                if not sub_id:
                    raise RuntimeError(
                        f"sub-category {sub_name!r} under {parent!r} not found after create"
                    )
                sub_ids[sub_name] = sub_id
        return sub_ids

    def assign_actor_category(self, actor_id: str, category_sub_id: str) -> None:
        path = f"/api/actor/{actor_id}/category/{category_sub_id}"
        status, body, raw = self.request("PUT", path)
        if status != 200:
            detail = (body or {}).get("error", raw)
            raise RuntimeError(
                f"assign category to actor {actor_id} failed ({status}): {detail}"
            )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--url",
        default=os.environ.get("ZT_INSTANCE_URL", "http://localhost:8069"),
        help="ZobTube base URL (default: $ZT_INSTANCE_URL or http://localhost:8069)",
    )
    parser.add_argument(
        "--username",
        default=os.environ.get("ZT_USERNAME", "validation"),
        help="Admin username when authentication is enabled",
    )
    parser.add_argument(
        "--password",
        default=os.environ.get("ZT_PASSWORD", "validation"),
        help="Admin password when authentication is enabled",
    )
    parser.add_argument("-v", "--verbose", action="store_true", help="Log HTTP traffic")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    client = ZobTubeClient(args.url, verbose=args.verbose)

    bootstrap = client.bootstrap()
    if bootstrap.get("auth_enabled"):
        client.login(args.username, args.password)

    country_ids = client.ensure_categories(CATEGORIES)
    print(f"categories ready: {len(country_ids)} countries")

    actors_by_name = client.actor_index(client.list_actors())
    created = 0
    for name, sex, country in ACTORS:
        actor_id, is_new = client.ensure_actor(name, sex, actors_by_name)
        country_id = country_ids[country]
        client.assign_actor_category(actor_id, country_id)
        verb = "created" if is_new else "exists"
        print(f"{verb} actor {name!r} ({actor_id}) -> {country}")
        if is_new:
            created += 1

    print(f"done: {created} new actors, {len(ACTORS)} total with countries")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except (urllib.error.URLError, RuntimeError) as exc:
        print(f"error: {exc}", file=sys.stderr)
        raise SystemExit(1) from exc
