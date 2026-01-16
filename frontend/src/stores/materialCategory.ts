import { defineStore } from 'pinia'
import { ref } from 'vue'
import { materialCategoryApi } from '@/api/materialCategory'
import type { MaterialCategoryConfig, MaterialCategoryConfigRequest } from '@/types'

export const useMaterialCategoryStore = defineStore('materialCategory', () => {
  const categories = ref<MaterialCategoryConfig[]>([])
  const activeCategories = ref<MaterialCategoryConfig[]>([])
  const loading = ref(false)

  const fetchAllCategories = async () => {
    loading.value = true
    try {
      const response = await materialCategoryApi.list(false)
      if (response.code === 0 && response.data) {
        categories.value = response.data
      }
      return response
    } finally {
      loading.value = false
    }
  }

  const fetchActiveCategories = async () => {
    loading.value = true
    try {
      const response = await materialCategoryApi.list(true)
      if (response.code === 0 && response.data) {
        activeCategories.value = response.data
      }
      return response
    } finally {
      loading.value = false
    }
  }

  const createCategory = async (data: MaterialCategoryConfigRequest) => {
    const response = await materialCategoryApi.create(data)
    if (response.code === 0 && response.data) {
      categories.value.push(response.data)
      if (response.data.is_active) {
        activeCategories.value.push(response.data)
      }
    }
    return response
  }

  const updateCategory = async (id: number, data: MaterialCategoryConfigRequest) => {
    const response = await materialCategoryApi.update(id, data)
    if (response.code === 0 && response.data) {
      const index = categories.value.findIndex(c => c.id === id)
      if (index !== -1) {
        categories.value[index] = response.data
      }

      const activeIndex = activeCategories.value.findIndex(c => c.id === id)
      if (response.data.is_active) {
        if (activeIndex !== -1) {
          activeCategories.value[activeIndex] = response.data
        } else {
          activeCategories.value.push(response.data)
        }
      } else if (activeIndex !== -1) {
        activeCategories.value.splice(activeIndex, 1)
      }
    }

    return response
  }

  const deleteCategory = async (id: number) => {
    const response = await materialCategoryApi.delete(id)
    if (response.code === 0) {
      categories.value = categories.value.filter(c => c.id !== id)
      activeCategories.value = activeCategories.value.filter(c => c.id !== id)
    }
    return response
  }

  const toggleStatus = async (id: number) => {
    const response = await materialCategoryApi.toggleStatus(id)
    if (response.code === 0 && response.data) {
      const index = categories.value.findIndex(c => c.id === id)
      if (index !== -1) {
        categories.value[index] = response.data
      }

      const activeIndex = activeCategories.value.findIndex(c => c.id === id)
      if (response.data.is_active) {
        if (activeIndex === -1) {
          activeCategories.value.push(response.data)
        }
      } else if (activeIndex !== -1) {
        activeCategories.value.splice(activeIndex, 1)
      }
    }

    return response
  }

  const getCategoryByCode = (code: string): MaterialCategoryConfig | undefined => {
    return activeCategories.value.find(c => c.code === code)
  }

  const getCategoryName = (code: string): string => {
    const category = getCategoryByCode(code)
    return category?.name_zh || category?.name || code
  }

  const reset = () => {
    categories.value = []
    activeCategories.value = []
    loading.value = false
  }

  return {
    categories,
    activeCategories,
    loading,
    fetchAllCategories,
    fetchActiveCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    toggleStatus,
    getCategoryByCode,
    getCategoryName,
    reset
  }
})
