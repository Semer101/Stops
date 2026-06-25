import { postJson } from "@/services/apiClient";
import type { AuthResponse, LoginPayload, RegisterPayload } from "@/types/auth";

const authTokenKey = "stops.auth.token";
const authUserKey = "stops.auth.user";

export async function register(payload: RegisterPayload): Promise<AuthResponse> {
  return postJson<AuthResponse, RegisterPayload>("/api/auth/register", { body: payload });
}

export async function login(payload: LoginPayload): Promise<AuthResponse> {
  return postJson<AuthResponse, LoginPayload>("/api/auth/login", { body: payload });
}

export function saveAuthSession(response: AuthResponse): void {
  window.localStorage.setItem(authTokenKey, response.data.token);
  window.localStorage.setItem(authUserKey, JSON.stringify(response.data.user));
}
