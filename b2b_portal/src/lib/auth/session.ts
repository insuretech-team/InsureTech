import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { getSession } from "./session-store";

export const SESSION_COOKIE_NAME = "session_token";

export async function getServerSession() {
  const cookieStore = await cookies();
  return getSession(cookieStore.get(SESSION_COOKIE_NAME)?.value);
}

export async function requireServerSession() {
  const session = await getServerSession();
  if (!session) {
    redirect("/login");
  }
  return session;
}
