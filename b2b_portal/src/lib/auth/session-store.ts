import crypto from "node:crypto";

import { create } from "@bufbuild/protobuf";

import {
  DeviceType,
  MoneySchema,
  SessionSchema,
  SessionType,
  UserSchema,
  UserStatus,
  UserType,
  type User,
} from "@lib/proto";
import type { PortalPrincipal, PortalSession } from "@lib/types/auth";

const SESSION_TTL_MS = 1000 * 60 * 60 * 12;

const sessionStore = new Map<string, PortalSession>();

function toTimestamp(milliseconds: number) {
  return {
    seconds: BigInt(Math.floor(milliseconds / 1000)),
    nanos: (milliseconds % 1000) * 1_000_000,
  };
}

export function createSession(principal: PortalPrincipal): PortalSession {
  // Ensure organisationName is always present (backwards compat for callers that omit it)
  if (principal.organisationName === undefined) {
    (principal as PortalPrincipal).organisationName = "";
  }
  const now = Date.now();
  const sessionId = crypto.randomUUID();
  const expiresAt = now + SESSION_TTL_MS;
  const csrfToken = crypto.randomBytes(16).toString("hex");

  const session = {
    session: create(SessionSchema, {
      sessionId,
      userId: principal.user.userId,
      sessionType: SessionType.SERVER_SIDE,
      sessionTokenLookup: crypto.createHash("sha256").update(sessionId).digest("hex"),
      expiresAt: toTimestamp(expiresAt),
      ipAddress: "127.0.0.1",
      userAgent: "b2b_portal",
      deviceId: "web-browser",
      deviceName: "Web Browser",
      deviceType: DeviceType.WEB,
      createdAt: toTimestamp(now),
      lastActivityAt: toTimestamp(now),
      isActive: true,
      csrfToken,
    }),
    principal,
    user: principal.user,
    expiresAt,
  } satisfies PortalSession;

  sessionStore.set(sessionId, session);
  return session;
}

export function getSession(sessionId: string | undefined): PortalSession | null {
  if (!sessionId) {
    return null;
  }
  const session = sessionStore.get(sessionId);
  if (!session) {
    return null;
  }
  if (session.expiresAt <= Date.now()) {
    sessionStore.delete(sessionId);
    return null;
  }
  return session;
}

export function clearSession(sessionId: string | undefined): void {
  if (!sessionId) {
    return;
  }
  sessionStore.delete(sessionId);
}
