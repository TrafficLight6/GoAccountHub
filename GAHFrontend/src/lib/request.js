import { ElMessage } from 'element-plus'

const BASE_URL = import.meta.env.VITE_API_BASE || '/api/v1'

/**
 * Async request wrapper
 * @param {Object} options
 * @param {string} options.url      API path (relative to BASE_URL, e.g. '/info')
 * @param {string} [options.method] HTTP method, defaults to GET
 * @param {Object} [options.params]  query parameters (GET)
 * @param {Object} [options.data]    request body (auto JSON-serialized for POST/PUT/DELETE)
 * @param {boolean} [options.showError=true] whether to show an error message automatically on failure
 * @param {boolean} [options.autoRedirect401=true] whether to redirect to the login page on 401
 * @returns {Promise<Object>} resolves with the full backend response body ({ code, data, ... })
 */
export async function request(options = {}) {
  const {
    url,
    method = 'GET',
    params,
    data,
    showError = true,
    autoRedirect401 = true,
  } = options

  if (!url) throw new Error('request: url is required')

  // Build the query string
  let fullUrl = BASE_URL + url
  if (params && Object.keys(params).length > 0) {
    const search = new URLSearchParams(
      Object.entries(params).reduce((acc, [k, v]) => {
        if (v !== undefined && v !== null && v !== '') acc[k] = String(v)
        return acc
      }, {})
    ).toString()
    if (search) fullUrl += (fullUrl.includes('?') ? '&' : '?') + search
  }

  const fetchOptions = {
    method: method.toUpperCase(),
    // Include cookies (authentication such as admin_token relies on it)
    credentials: 'include',
    headers: {},
  }

  if (data !== undefined && data !== null) {
    fetchOptions.headers['Content-Type'] = 'application/json'
    fetchOptions.body = JSON.stringify(data)
  }

  let res
  try {
    res = await fetch(fullUrl, fetchOptions)
  } catch (e) {
    if (showError) ElMessage.error('Network request failed, please check your network connection')
    throw e
  }

  // 401 not logged in / session expired: redirect back to the login page
  if (res.status === 401 && autoRedirect401) {
    if (showError) ElMessage.error('Session expired, please log in again')
    window.location.href = '/'
    throw new Error('Unauthorized')
  }

  // Parse the response body (the backend always returns JSON)
  let body
  try {
    body = await res.json()
  } catch (e) {
    if (showError) ElMessage.error('Failed to parse response')
    throw new Error('Invalid JSON response')
  }

  // Backend convention: code === 200 means success
  if (!res.ok || (body.code !== undefined && body.code !== 200)) {
    const errMsg = body.error || body.message || `Request failed (${res.status})`
    if (showError) ElMessage.error(errMsg)
    const err = new Error(errMsg)
    err.code = body.code ?? res.status
    err.body = body
    throw err
  }

  return body
}

export const get = (url, params, options = {}) =>
  request({ ...options, url, method: 'GET', params })

export const post = (url, data, options = {}) =>
  request({ ...options, url, method: 'POST', data })

export const put = (url, data, options = {}) =>
  request({ ...options, url, method: 'PUT', data })

export const del = (url, data, options = {}) =>
  request({ ...options, url, method: 'DELETE', data })

export default request
