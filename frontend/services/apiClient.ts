import type { ApiErrorResponse } from "@/types/auth";

const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export class ApiClientError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiClientError";
    this.status = status;
  }
}

interface RequestOptions<TBody> {
  body?: TBody;
  token?: string;
}

export async function postJson<TResponse, TBody>(
  path: string,
  options: RequestOptions<TBody>,
): Promise<TResponse> {
  const response = await fetch(`${apiBaseUrl}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {}),
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  if (!response.ok) {
    const errorBody = (await response.json().catch(() => null)) as ApiErrorResponse | null;
    throw new ApiClientError(errorBody?.message ?? "Request failed.", response.status);
  }

  return response.json() as Promise<TResponse>;
}

export async function getJson<TResponse>(
  path: string,
  params: Record<string, string | number | undefined> = {},
): Promise<TResponse> {
  const url = new URL(`${apiBaseUrl}${path}`);
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined) {
      url.searchParams.set(key, String(value));
    }
  });

  const response = await fetch(url);

  if (!response.ok) {
    const errorBody = (await response.json().catch(() => null)) as ApiErrorResponse | null;
    throw new ApiClientError(errorBody?.message ?? "Request failed.", response.status);
  }

  return response.json() as Promise<TResponse>;
}
