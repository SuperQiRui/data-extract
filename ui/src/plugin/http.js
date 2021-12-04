import conf from '@/settings'
import axios from 'axios'
import qs from 'qs'
import { ElMessage } from 'element-plus'

// 创建axios实例
const service = axios.create({
  baseURL: conf.domain,
  withCredentials: false,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded'
  }
})

export const get = async(api, params) => {
  return new Promise((resolve, reject) => {
    service.get(api, {
      params,
      paramsSerializer: function(params) {
        return qs.stringify(params, { arrayFormat: 'repeat', indices: false })
      }
    }).then(res => {
      if (res.status === 200) {
        if (res.data.Status === 1) {
          if (res.data.Msg) {
            ElMessage.success(res.data.Msg)
          }
          resolve(res.data.Result)
        } else {
          if (res.data.Msg) {
            ElMessage.error(res.data.Msg)
          }
          reject()
        }
      }
    }).catch((e) => {
      reject(e)
    })
  })
}

export const post = async(api, params) => {
  return new Promise((resolve, reject) => {
    service.post(api, params, { indices: false }).then(res => {
      if (res.status === 200) {
        if (res.data.Status === 1) {
          if (res.data.Msg) {
            ElMessage.success(res.data.Msg)
          }
          resolve(res.data.Result)
        } else {
          if (res.data.Msg) {
            ElMessage.error(res.data.Msg)
          }
          reject()
        }
      } else {
        reject()
      }
    }).catch((e) => {
      reject(e)
    })
  })
}

export const download = async(api, filename, params) => {
  return await service.post(api, params, {
    responseType: 'blob'
  }).then(res => {
    if (res.status === 200) {
      const blob = new Blob([res.data])
      if (res.headers['content-type'] === 'application/octet-stream') {
        saveAs(blob, filename)
      } else {
        blob.text().then(text => {
          const msg = JSON.parse(text)
          ElMessage.error(msg.Msg)
        })
      }
    }
  }).catch((e) => {
    console.log(e)
  })
}

export const saveAs = (blob, filename) => {
  const force_saveable_type = 'application/octet-stream'
  if (blob.type && blob.type !== force_saveable_type) { // 强制下载，而非在浏览器中打开
    const slice = blob.slice || blob.webkitSlice || blob.mozSlice
    blob = slice.call(blob, 0, blob.size, force_saveable_type)
  }
  const url = window.URL.createObjectURL(blob)
  const save_link = document.createElementNS('http://www.w3.org/1999/xhtml', 'a')
  save_link.href = url
  save_link.download = filename
  save_link.dispatchEvent(new Event('click'))
  save_link.click()
  window.URL.revokeObjectURL(url)
}

