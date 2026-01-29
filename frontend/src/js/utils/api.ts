import type { APIResponse } from "../types";

/**
 * Base request function
 * @param url - API endpoint
 * @param options - Fetch options
 * @returns Response data
 * @throws Error on HTTP error or API error status
 */
const apiRequest = async <T>(
  url: string,
  options: RequestInit = {},
): Promise<APIResponse<T>> => {
  const response = await fetch(url, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const data = (await response.json()) as APIResponse<T>;

  if (data.status === "error") {
    throw new Error("data" in data ? String(data.data) : "API error");
  }

  return data;
};

/**
 * GET request
 */
export const apiGet = async <T>(url: string): Promise<APIResponse<T>> => {
  return apiRequest<T>(url, { method: "GET" });
};

/**
 * POST request
 */
export const apiPost = async <T>(
  url: string,
  data: unknown,
): Promise<APIResponse<T>> => {
  return apiRequest<T>(url, {
    method: "POST",
    body: JSON.stringify(data),
  });
};

/**
 * PUT request
 */
export const apiPut = async <T>(
  url: string,
  data: unknown,
): Promise<APIResponse<T>> => {
  return apiRequest<T>(url, {
    method: "PUT",
    body: JSON.stringify(data),
  });
};

/**
 * DELETE request
 */
export const apiDelete = async <T>(url: string): Promise<APIResponse<T>> => {
  return apiRequest<T>(url, { method: "DELETE" });
};
