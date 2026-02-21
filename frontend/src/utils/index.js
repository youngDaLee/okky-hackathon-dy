// Utility functions

export const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleDateString('ko-KR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

export const formatExpiryStatus = (status) => {
  const statusMap = {
    URGENT: '긴급',
    SOON: '임박',
    NORMAL: '정상',
    NO_EXPIRY: '기한 없음'
  }
  return statusMap[status] || status
}

export const getExpiryStatusColor = (status) => {
  const colorMap = {
    URGENT: 'text-red-600 bg-red-50',
    SOON: 'text-orange-600 bg-orange-50',
    NORMAL: 'text-gray-600 bg-gray-50',
    NO_EXPIRY: 'text-gray-400 bg-gray-50'
  }
  return colorMap[status] || 'text-gray-600 bg-gray-50'
}
