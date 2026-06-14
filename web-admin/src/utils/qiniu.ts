export async function uploadSpImage(file: File): Promise<{ url: string; key: string }> {
  const formData = new FormData()
  formData.append('file', file)
  const token = window.localStorage.getItem('admin_token')

  const response = await fetch('/api/v1/upload/file', {
    method: 'POST',
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    body: formData,
  })

  if (!response.ok) {
    throw new Error('上传图片失败')
  }

  const result = await response.json()
  if (result?.code !== 0 || !result?.data?.path) {
    throw new Error(result?.message || '上传图片失败')
  }

  return {
    key: result.data.path,
    url: result.data.url,
  }
}
