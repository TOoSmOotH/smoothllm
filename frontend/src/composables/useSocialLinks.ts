import { ref, computed } from 'vue'
import { toast } from 'vue-sonner'
import {
  SocialLink,
  CreateSocialLinkRequest,
  UpdateSocialLinkRequest,
  ReorderSocialLinksRequest,
  getSocialLinks,
  createSocialLink,
  updateSocialLink,
  deleteSocialLink,
  reorderSocialLinks,
} from '@/api/social'

export function useSocialLinks() {
  const socialLinks = ref<SocialLink[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const visibleLinks = computed(() => socialLinks.value.filter((link) => link.visible))
  const sortedLinks = computed(() => {
    return [...socialLinks.value].sort((a, b) => a.order - b.order)
  })

  const fetchSocialLinks = async () => {
    loading.value = true
    error.value = null
    try {
      socialLinks.value = await getSocialLinks()
    } catch (err: any) {
      error.value = err.message || 'Failed to load social links'
      toast.error(error.value)
    } finally {
      loading.value = false
    }
  }

  const addSocialLink = async (data: CreateSocialLinkRequest) => {
    loading.value = true
    error.value = null
    try {
      const newLink = await createSocialLink(data)
      socialLinks.value.push(newLink)
      toast.success('Social link added successfully')
      return newLink
    } catch (err: any) {
      error.value = err.message || 'Failed to add social link'
      toast.error(error.value)
      throw err
    } finally {
      loading.value = false
    }
  }

  const editSocialLink = async (id: number, data: UpdateSocialLinkRequest) => {
    loading.value = true
    error.value = null
    try {
      const updatedLink = await updateSocialLink(id, data)
      const index = socialLinks.value.findIndex((link) => link.id === id)
      if (index !== -1) {
        socialLinks.value[index] = updatedLink
      }
      toast.success('Social link updated successfully')
      return updatedLink
    } catch (err: any) {
      error.value = err.message || 'Failed to update social link'
      toast.error(error.value)
      throw err
    } finally {
      loading.value = false
    }
  }

  const removeSocialLink = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      await deleteSocialLink(id)
      socialLinks.value = socialLinks.value.filter((link) => link.id !== id)
      toast.success('Social link removed successfully')
    } catch (err: any) {
      error.value = err.message || 'Failed to remove social link'
      toast.error(error.value)
      throw err
    } finally {
      loading.value = false
    }
  }

  const reorderSocialLinksItems = async (updates: ReorderSocialLinksRequest) => {
    loading.value = true
    error.value = null
    try {
      const reorderedLinks = await reorderSocialLinks(updates)
      socialLinks.value = reorderedLinks
      toast.success('Social links reordered successfully')
      return reorderedLinks
    } catch (err: any) {
      error.value = err.message || 'Failed to reorder social links'
      toast.error(error.value)
      throw err
    } finally {
      loading.value = false
    }
  }

  const toggleLinkVisibility = async (id: number) => {
    const link = socialLinks.value.find((l) => l.id === id)
    if (!link) return

    try {
      await editSocialLink(id, { visible: !link.visible })
    } catch (err) {
      throw err
    }
  }

  return {
    socialLinks,
    loading,
    error,
    visibleLinks,
    sortedLinks,
    fetchSocialLinks,
    addSocialLink,
    editSocialLink,
    removeSocialLink,
    reorderSocialLinksItems,
    toggleLinkVisibility,
  }
}
