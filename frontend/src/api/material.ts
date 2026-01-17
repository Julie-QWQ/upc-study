import { request } from '@/utils/request'
import type {
  ApiResponse,
  PaginateData,
  Material,
  MaterialListParams,
  CreateMaterialRequest,
  UpdateMaterialRequest,
  UploadSignature,
  UploadSignatureRequest,
  DownloadSignature,
  Favorite,
  FavoriteListParams,
  DownloadRecord,
  DownloadRecordListParams,
  Report,
  ReportListParams,
  CreateReportRequest,
  HandleReportRequest
} from '@/types'

/**
 * 资料管理相关 API
 */
export const materialApi = {
  /**
   * 获取资料列表
   */
  getMaterials(params: MaterialListParams): Promise<ApiResponse<PaginateData<Material>>> {
    return request.get('/materials', { params })
  },

  /**
   * 获取已审核资料列表（管理员）
   */
  getReviewedMaterials(params: MaterialListParams): Promise<ApiResponse<PaginateData<Material>>> {
    return request.get('/admin/materials/reviewed', { params })
  },

  /**
   * 搜索资料
   */
  searchMaterials(params: MaterialListParams): Promise<ApiResponse<PaginateData<Material>>> {
    return request.get('/materials/search', { params })
  },

  /**
   * 获取资料详情
   */
  getMaterial(id: number): Promise<ApiResponse<Material>> {
    return request.get(`/materials/${id}`)
  },

  /**
   * 创建资料（学委及以上权限）
   */
  createMaterial(data: CreateMaterialRequest): Promise<ApiResponse<Material>> {
    return request.post('/materials', data)
  },

  /**
   * 更新资料（学委及以上权限）
   */
  updateMaterial(id: number, data: UpdateMaterialRequest): Promise<ApiResponse<Material>> {
    return request.put(`/materials/${id}`, data)
  },

  /**
   * 删除资料（管理员权限）
   */
  deleteMaterial(id: number): Promise<ApiResponse<null>> {
    return request.delete(`/materials/${id}`)
  },

  /**
   * 获取上传签名（学委及以上权限）
   */
  getUploadSignature(data: UploadSignatureRequest): Promise<ApiResponse<UploadSignature>> {
    return request.post('/materials/upload-signature', data)
  },

  /**
   * 删除已上传但未创建记录的文件（学委及以上权限）
   */
  deleteUploadedFile(fileKey: string): Promise<ApiResponse<null>> {
    return request.post('/materials/delete-uploaded-file', { file_key: fileKey })
  },

  /**
   * 获取下载链接
   */
  getDownloadUrl(id: number): Promise<ApiResponse<DownloadSignature>> {
    return request.get(`/materials/${id}/download`)
  },

  /**
   * 审核资料（管理员权限）
   */
  reviewMaterial(id: number, approved: boolean, reason?: string): Promise<ApiResponse<Material>> {
    return request.post(`/materials/${id}/review`, {
      approved,
      rejection_reason: reason
    })
  },

  /**
   * 添加收藏
   */
  addFavorite(materialId: number): Promise<ApiResponse<null>> {
    return request.post(`/materials/${materialId}/favorite`)
  },

  /**
   * 取消收藏
   */
  removeFavorite(materialId: number): Promise<ApiResponse<null>> {
    return request.delete(`/materials/${materialId}/favorite`)
  },

  /**
   * 获取收藏列表
   */
  getFavorites(params: FavoriteListParams): Promise<ApiResponse<PaginateData<Material>>> {
    return request.get('/favorites', { params })
  },

  /**
   * 举报资料
   */
  createReport(materialId: number, data: CreateReportRequest): Promise<ApiResponse<Report>> {
    return request.post(`/materials/${materialId}/report`, data)
  },

  /**
   * 获取下载记录列表
   */
  getDownloadRecords(params: { page?: number; page_size?: number }): Promise<ApiResponse<PaginateData<Material>>> {
    return request.get('/downloads', { params })
  }
}

/**
 * 举报管理相关 API（管理员）
 */
export const reportApi = {
  /**
   * 获取举报列表
   */
  getReports(params: ReportListParams): Promise<ApiResponse<PaginateData<Report>>> {
    return request.get('/admin/reports', { params })
  },

  /**
   * 获取举报详情
   */
  getReport(id: number): Promise<ApiResponse<Report>> {
    return request.get(`/admin/reports/${id}`)
  },

  /**
   * 处理举报（管理员权限）
   */
  handleReport(id: number, data: HandleReportRequest): Promise<ApiResponse<Report>> {
    return request.post(`/reports/${id}/handle`, data)
  }
}
