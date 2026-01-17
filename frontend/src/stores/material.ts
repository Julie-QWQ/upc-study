import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { materialApi } from '@/api/material'
import type {
  Material,
  MaterialListParams,
  CreateMaterialRequest,
  UpdateMaterialRequest,
  UploadSignatureRequest
} from '@/types'

export const useMaterialStore = defineStore('material', () => {
  // ==================== 状态 ====================
  const materials = ref<Material[]>([])
  const currentMaterial = ref<Material | null>(null)
  const total = ref(0)
  const page = ref(1)
  const size = ref(20)
  const loading = ref(false)

  // 当前查询参数
  const currentParams = ref<MaterialListParams>({
    page: 1,
    size: 20
  })

  // ==================== 计算属性 ====================
  const hasMore = computed(() => {
    return page.value * size.value < total.value
  })

  const totalPages = computed(() => {
    return Math.ceil(total.value / size.value)
  })

  // ==================== 方法 ====================

  /**
   * 获取资料列表
   */
  const fetchMaterials = async (params?: MaterialListParams) => {
    loading.value = true
    try {
      const queryParams = {
        ...currentParams.value,
        ...params
      }

      const response = await materialApi.getMaterials(queryParams)

      if (response.code === 0 && response.data) {
        // 后端返回的是 materials 和 page_size
        materials.value = response.data.materials || response.data.list || []
        total.value = response.data.total
        page.value = response.data.page
        size.value = response.data.page_size || response.data.size || 20
        currentParams.value = queryParams
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取已审核资料列表（管理员）
   */
  const fetchReviewedMaterials = async (params?: MaterialListParams) => {
    loading.value = true
    try {
      const queryParams = {
        ...currentParams.value,
        ...params
      }

      const response = await materialApi.getReviewedMaterials(queryParams)

      if (response.code === 0 && response.data) {
        materials.value = response.data.materials || response.data.list || []
        total.value = response.data.total
        page.value = response.data.page
        size.value = response.data.page_size || response.data.size || 20
        currentParams.value = queryParams
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 搜索资料
   */
  const searchMaterials = async (params: MaterialListParams) => {
    loading.value = true
    try {
      const response = await materialApi.searchMaterials(params)

      if (response.code === 0 && response.data) {
        materials.value = response.data.materials || response.data.list || []
        total.value = response.data.total
        page.value = response.data.page
        size.value = response.data.page_size || response.data.size || 20
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取资料详情
   */
  const fetchMaterial = async (id: number) => {
    loading.value = true
    try {
      const response = await materialApi.getMaterial(id)

      if (response.code === 0 && response.data) {
        currentMaterial.value = response.data
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 创建资料
   */
  const createMaterial = async (data: CreateMaterialRequest) => {
    const response = await materialApi.createMaterial(data)

    if (response.code === 0) {
      // 刷新列表
      await fetchMaterials()
    }

    return response
  }

  /**
   * 更新资料
   */
  const updateMaterial = async (id: number, data: UpdateMaterialRequest) => {
    const response = await materialApi.updateMaterial(id, data)

    if (response.code === 0 && response.data) {
      // 更新列表中的数据
      const index = materials.value.findIndex(m => m.id === id)
      if (index !== -1) {
        materials.value[index] = response.data
      }

      // 更新当前资料
      if (currentMaterial.value?.id === id) {
        currentMaterial.value = response.data
      }
    }

    return response
  }

  /**
   * 删除资料
   */
  const deleteMaterial = async (id: number) => {
    const response = await materialApi.deleteMaterial(id)

    if (response.code === 0) {
      // 从列表中移除
      materials.value = materials.value.filter(m => m.id !== id)

      // 清空当前资料
      if (currentMaterial.value?.id === id) {
        currentMaterial.value = null
      }
    }

    return response
  }

  /**
   * 获取上传签名
   */
  const getUploadSignature = async (data: UploadSignatureRequest) => {
    return await materialApi.getUploadSignature(data)
  }

  /**
   * 获取下载链接
   */
  const getDownloadUrl = async (id: number) => {
    return await materialApi.getDownloadUrl(id)
  }

  /**
   * 审核资料
   */
  const reviewMaterial = async (id: number, approved: boolean, reason?: string) => {
    const response = await materialApi.reviewMaterial(id, approved, reason)

    if (response.code === 0 && response.data) {
      // 更新列表中的数据
      const index = materials.value.findIndex(m => m.id === id)
      if (index !== -1) {
        materials.value[index] = response.data
      }

      // 更新当前资料
      if (currentMaterial.value?.id === id) {
        currentMaterial.value = response.data
      }
    }

    return response
  }

  /**
   * 添加收藏
   */
  const addFavorite = async (materialId: number) => {
    // 检查是否已经收藏
    const material = materials.value.find(m => m.id === materialId)
    if (material?.is_favorited || currentMaterial.value?.is_favorited) {
      // 如果已经收藏,直接返回成功,不做任何操作
      return { code: 0, message: '已收藏', data: null }
    }

    const response = await materialApi.addFavorite(materialId)

    // code 0 表示成功, code 10006 表示已收藏(也视为成功)
    if (response.code === 0 || response.code === 10006) {
      // 更新列表中的收藏状态
      const material = materials.value.find(m => m.id === materialId)
      if (material) {
        material.is_favorited = true
        material.favorite_count++
      }

      // 更新当前资料的收藏状态
      if (currentMaterial.value?.id === materialId) {
        currentMaterial.value.is_favorited = true
        currentMaterial.value.favorite_count++
      }

      // 如果是重复收藏,返回成功消息
      if (response.code === 10006) {
        return { code: 0, message: '已收藏', data: null }
      }
    }

    return response
  }

  /**
   * 取消收藏
   */
  const removeFavorite = async (materialId: number) => {
    const response = await materialApi.removeFavorite(materialId)

    if (response.code === 0) {
      // 更新列表中的收藏状态
      const material = materials.value.find(m => m.id === materialId)
      if (material) {
        material.is_favorited = false
        material.favorite_count--
      }

      // 更新当前资料的收藏状态
      if (currentMaterial.value?.id === materialId) {
        currentMaterial.value.is_favorited = false
        currentMaterial.value.favorite_count--
      }
    }

    return response
  }

  /**
   * 获取收藏列表
   */
  const fetchFavorites = async (params?: MaterialListParams) => {
    loading.value = true
    try {
      const response = await materialApi.getFavorites(params || {})

      if (response.code === 0 && response.data) {
        materials.value = response.data.materials || response.data.list || []
        total.value = response.data.total
        page.value = response.data.page
        size.value = response.data.page_size || response.data.size || 20
      }

      return response
    } finally {
      loading.value = false
    }
  }

  /**
   * 举报资料
   */
  const createReport = async (materialId: number, reason: string, description: string) => {
    return await materialApi.createReport(materialId, { reason: reason as any, description })
  }

  /**
   * 重置状态
   */
  const reset = () => {
    materials.value = []
    currentMaterial.value = null
    total.value = 0
    page.value = 1
    size.value = 20
    loading.value = false
    currentParams.value = {
      page: 1,
      size: 20
    }
  }

  return {
    // 状态
    materials,
    currentMaterial,
    total,
    page,
    size,
    loading,
    currentParams,

    // 计算属性
    hasMore,
    totalPages,

    // 方法
    fetchMaterials,
    fetchReviewedMaterials,
    searchMaterials,
    fetchMaterial,
    createMaterial,
    updateMaterial,
    deleteMaterial,
    getUploadSignature,
    getDownloadUrl,
    reviewMaterial,
    addFavorite,
    removeFavorite,
    fetchFavorites,
    createReport,
    reset
  }
})
