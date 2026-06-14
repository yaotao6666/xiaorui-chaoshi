import request, { unwrapApiResponse } from '@/utils/request'
import type {
  AmountTrendData,
  AnnouncementFormData,
  AnnouncementItem,
  AnnouncementListResponse,
  DashboardData,
  MerchantCategory,
  MerchantCategoryPayload,
  MerchantDetail,
  MerchantDistributionData,
  MerchantListResponse,
  MerchantPickupPoint,
  MerchantPickupPointPayload,
  MerchantProduct,
  MerchantProductListResponse,
  MerchantProductPayload,
  OrderAnalyticsData,
  SpOrder,
  SpOrderListResponse,
  ServiceProviderLoginRequest,
  ServiceProviderLoginResponse,
  SpMerchantFormData,
  TopMerchantRanking,
  UpdateSpMerchantFormData,
} from '@/types/sp'

export function spLogin(data: ServiceProviderLoginRequest) {
  return request.post('/api/v1/admin/auth/login', data).then(unwrapApiResponse<ServiceProviderLoginResponse>)
}

export function spLogout() {
  return request.post('/api/v1/admin/auth/logout', {}).then(unwrapApiResponse<{ message: string }>)
}

export function getDashboard() {
  return request.get('/api/v1/admin/dashboard').then(unwrapApiResponse<DashboardData>)
}

export function getMerchantList(params?: Record<string, unknown>) {
  return request.get('/api/v1/admin/stores/list', { params }).then(unwrapApiResponse<MerchantListResponse>)
}

export function getMerchantDetail(merchantId: number) {
  return request.get(`/api/v1/admin/stores/${merchantId}`).then(unwrapApiResponse<MerchantDetail>)
}

export function createSpMerchant(data: SpMerchantFormData) {
  return request.post('/api/v1/admin/stores', data).then(unwrapApiResponse<MerchantDetail>)
}

export function updateSpMerchant(merchantId: number, data: UpdateSpMerchantFormData) {
  return request.put(`/api/v1/admin/stores/${merchantId}`, data).then(unwrapApiResponse<MerchantDetail>)
}

export function resetSpMerchantAdminPassword(merchantId: number, data: { new_password: string }) {
  return request.post(`/api/v1/admin/stores/${merchantId}/admin/reset-password`, data).then(
    unwrapApiResponse<{ staff_id: number; username: string; message: string }>
  )
}

export function updateSpMerchantAssets(merchantId: number, data: { logo?: string; cover_image?: string }) {
  return request.put(`/api/v1/admin/stores/${merchantId}/assets`, data).then(unwrapApiResponse<MerchantDetail>)
}

export function getMerchantPickupPoints(merchantId: number) {
  return request
    .get(`/api/v1/admin/stores/${merchantId}/pickup-points`)
    .then(unwrapApiResponse<MerchantPickupPoint[]>)
}

export function createMerchantPickupPoint(merchantId: number, data: MerchantPickupPointPayload) {
  return request
    .post(`/api/v1/admin/stores/${merchantId}/pickup-points`, data)
    .then(unwrapApiResponse<MerchantPickupPoint>)
}

export function updateMerchantPickupPoint(merchantId: number, pickupPointId: number, data: MerchantPickupPointPayload) {
  return request
    .put(`/api/v1/admin/stores/${merchantId}/pickup-points/${pickupPointId}`, data)
    .then(unwrapApiResponse<MerchantPickupPoint>)
}

export function deleteMerchantPickupPoint(merchantId: number, pickupPointId: number) {
  return request
    .delete(`/api/v1/admin/stores/${merchantId}/pickup-points/${pickupPointId}`)
    .then(unwrapApiResponse<{ message: string }>)
}

export function getMerchantCategories(merchantId: number) {
  return request
    .get(`/api/v1/admin/stores/${merchantId}/categories`)
    .then(unwrapApiResponse<MerchantCategory[]>)
}

export function createMerchantCategory(merchantId: number, data: MerchantCategoryPayload) {
  return request
    .post(`/api/v1/admin/stores/${merchantId}/categories`, data)
    .then(unwrapApiResponse<MerchantCategory>)
}

export function updateMerchantCategory(merchantId: number, categoryId: number, data: MerchantCategoryPayload) {
  return request
    .put(`/api/v1/admin/stores/${merchantId}/categories/${categoryId}`, data)
    .then(unwrapApiResponse<MerchantCategory>)
}

export function deleteMerchantCategory(merchantId: number, categoryId: number) {
  return request
    .delete(`/api/v1/admin/stores/${merchantId}/categories/${categoryId}`)
    .then(unwrapApiResponse<{ message: string }>)
}

export function getMerchantProducts(merchantId: number, params?: Record<string, unknown>) {
  return request
    .get(`/api/v1/admin/stores/${merchantId}/products`, { params })
    .then(unwrapApiResponse<MerchantProductListResponse>)
}

export function getMerchantProduct(merchantId: number, productId: number) {
  return request
    .get(`/api/v1/admin/stores/${merchantId}/products/${productId}`)
    .then(unwrapApiResponse<MerchantProduct>)
}

export function createMerchantProduct(merchantId: number, data: MerchantProductPayload) {
  return request
    .post(`/api/v1/admin/stores/${merchantId}/products`, data)
    .then(unwrapApiResponse<MerchantProduct>)
}

export function updateMerchantProduct(merchantId: number, productId: number, data: MerchantProductPayload) {
  return request
    .put(`/api/v1/admin/stores/${merchantId}/products/${productId}`, data)
    .then(unwrapApiResponse<MerchantProduct>)
}

export function merchantProductOnSale(merchantId: number, productId: number) {
  return request
    .post(`/api/v1/admin/stores/${merchantId}/products/${productId}/on-sale`)
    .then(unwrapApiResponse<{ message: string }>)
}

export function merchantProductOffSale(merchantId: number, productId: number) {
  return request
    .post(`/api/v1/admin/stores/${merchantId}/products/${productId}/off-sale`)
    .then(unwrapApiResponse<{ message: string }>)
}

export function batchUpdateMerchantProductStatus(merchantId: number, productIds: number[], status: number) {
  return request
    .post(`/api/v1/admin/stores/${merchantId}/products/batch-status`, {
      product_ids: productIds,
      status,
    })
    .then(unwrapApiResponse<{ message: string }>)
}

export function deleteMerchantProduct(merchantId: number, productId: number) {
  return request
    .delete(`/api/v1/admin/stores/${merchantId}/products/${productId}`)
    .then(unwrapApiResponse<{ message: string }>)
}

export function getMerchantDistribution(params?: Record<string, unknown>) {
  return request.get('/api/v1/admin/stores/analytics/distribution', { params }).then(unwrapApiResponse<MerchantDistributionData>)
}

export function getOrderAnalytics(params?: Record<string, unknown>) {
  return request.get('/api/v1/admin/orders/analytics', { params }).then(unwrapApiResponse<OrderAnalyticsData>)
}

export function getSpOrders(params?: Record<string, unknown>) {
  return request.get('/api/v1/admin/orders', { params }).then(unwrapApiResponse<SpOrderListResponse>)
}

export function getSpOrderDetail(orderId: number) {
  return request.get(`/api/v1/admin/orders/${orderId}`).then(unwrapApiResponse<SpOrder>)
}

export function getAmountAnalytics(params?: Record<string, unknown>) {
  return request.get('/api/v1/admin/amount/analytics', { params }).then(unwrapApiResponse<AmountTrendData>)
}

export function getTopMerchants(params?: Record<string, unknown>) {
  return request.get('/api/v1/admin/amount/top-stores', { params }).then(unwrapApiResponse<TopMerchantRanking[]>)
}

export function changeSpPassword(data: { old_password: string; new_password: string }) {
  return request.post('/api/v1/admin/account/change-password', data).then(unwrapApiResponse<{ message: string }>)
}

export function getMerchantQRCode(merchantId: number) {
  return request.get(`/api/v1/admin/stores/${merchantId}/qrcode`).then(
    unwrapApiResponse<{ merchant_id: number; merchant_name: string; qrcode_url: string; page_path: string }>
  )
}

export function getAnnouncements(params?: Record<string, unknown>) {
  return request.get('/api/v1/admin/announcements', { params }).then(unwrapApiResponse<AnnouncementListResponse>)
}

export function getAnnouncementDetail(announcementId: number) {
  return request.get(`/api/v1/admin/announcements/${announcementId}`).then(unwrapApiResponse<AnnouncementItem>)
}

export function createAnnouncement(data: AnnouncementFormData) {
  return request.post('/api/v1/admin/announcements', data).then(
    unwrapApiResponse<{ id: number; message: string }>
  )
}

export function updateAnnouncement(announcementId: number, data: AnnouncementFormData) {
  return request.put(`/api/v1/admin/announcements/${announcementId}`, data).then(unwrapApiResponse<AnnouncementItem>)
}

export function deleteAnnouncement(announcementId: number) {
  return request.delete(`/api/v1/admin/announcements/${announcementId}`).then(
    unwrapApiResponse<{ message: string }>
  )
}
