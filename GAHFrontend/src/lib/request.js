import { ElMessage } from 'element-plus'

const BASE_URL = import.meta.env.VITE_API_BASE || '/api/v1'

/**
 * 异步请求封装
 * @param {Object} options
 * @param {string} options.url      接口路径（相对 BASE_URL，如 '/info'）
 * @param {string} [options.method] 请求方法，默认 GET
 * @param {Object} [options.params]  query 参数（GET）
 * @param {Object} [options.data]    请求体（POST/PUT/DELETE 自动 JSON 序列化）
 * @param {boolean} [options.showError=true] 失败时是否自动弹出错误提示
 * @param {boolean} [options.autoRedirect401=true] 401 时是否自动跳转登录页
 * @returns {Promise<Object>} resolve 后端完整响应体（{ code, data, ... }）
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

  // 拼接 query string
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
    // 携带 cookie（admin_token 等鉴权依赖）
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
    if (showError) ElMessage.error('网络请求失败，请检查网络连接')
    throw e
  }

  // 401 未登录/登录过期：跳回登录页
  if (res.status === 401 && autoRedirect401) {
    if (showError) ElMessage.error('登录已过期，请重新登录')
    window.location.href = '/'
    throw new Error('Unauthorized')
  }

  // 解析响应体（后端始终返回 JSON）
  let body
  try {
    body = await res.json()
  } catch (e) {
    if (showError) ElMessage.error('响应解析失败')
    throw new Error('Invalid JSON response')
  }

  // 后端约定：code === 200 为成功
  if (!res.ok || (body.code !== undefined && body.code !== 200)) {
    const errMsg = body.error || body.message || `请求失败（${res.status}）`
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
