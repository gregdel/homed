/**
 * API utility functions for making HTTP requests
 * Handles common patterns: error checking, JSON parsing, response validation
 */

/**
 * Base request function
 * @param {string} url - API endpoint
 * @param {object} options - Fetch options
 * @returns {Promise<any>} Response data
 * @throws {Error} On HTTP error or API error status
 */
const apiRequest = async (url, options = {}) => {
  const response = await fetch(url, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const data = await response.json();

  if (data.status === "error") {
    throw new Error(data.data || "API error");
  }

  return data;
};

/**
 * GET request
 */
export const apiGet = async (url) => {
  return apiRequest(url, { method: "GET" });
};

/**
 * POST request
 */
export const apiPost = async (url, data) => {
  return apiRequest(url, {
    method: "POST",
    body: JSON.stringify(data),
  });
};

/**
 * PUT request
 */
export const apiPut = async (url, data) => {
  return apiRequest(url, {
    method: "PUT",
    body: JSON.stringify(data),
  });
};

/**
 * DELETE request
 */
export const apiDelete = async (url) => {
  return apiRequest(url, { method: "DELETE" });
};
