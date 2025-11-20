// store/reviews.js
import { defineStore } from 'pinia'
import { reviewsAPI } from '@/api/reviews'

export const useReviewsStore = defineStore('reviews', {
  state: () => ({
    reviews: [],
    loading: false,
    error: null,

    filters: {
      property: 'all',
      channel: 'all',
      rating: 'all',
      dateRange: null,
    },

    sortBy: 'date', // 'date' or 'rating'
  }),

  getters: {
    // Filtered and sorted reviews for dashboard
    filteredReviews: (state) => {
      let filtered = state.reviews.filter(r => {
        const matchesProperty =
          state.filters.property === 'all' || r.PropertyName === state.filters.property

        const matchesChannel =
          state.filters.channel === 'all' || r.Channel === state.filters.channel

        const matchesRating =
          state.filters.rating === 'all' ||
          (() => {
            const [min, max] = state.filters.rating.split('-').map(Number)
            return (r.Rating || 0) >= min && (r.Rating || 0) <= max
          })()

        return matchesProperty && matchesChannel && matchesRating
      })

      // Sorting
      filtered.sort((a, b) => {
        if (state.sortBy === 'date') {
          return new Date(b.SubmittedAt || 0) - new Date(a.SubmittedAt || 0)
        }
        if (state.sortBy === 'rating') {
          return (b.Rating || 0) - (a.Rating || 0)
        }
        return 0
      })

      return filtered
    },

    // ✅ FIXED: Now accepts propertyId parameter
    // Returns approved reviews for a specific property
    approvedReviews: (state) => (propertyId) => {
      return state.reviews.filter(r => {
        const isApproved = r.ApprovalStatus?.IsApproved
        
        // If no propertyId provided, return all approved reviews
        if (!propertyId) return isApproved
        
        // Match by PropertyID or PropertyName
        const matchesProperty = 
          r.PropertyID === propertyId || 
          r.PropertyName === propertyId ||
          r.ListingID === propertyId
        
        return isApproved && matchesProperty
      })
    },

    // For dashboard filters
    properties: (state) => [...new Set(state.reviews.map(r => r.PropertyName))],
    channels: (state) => [...new Set(state.reviews.map(r => r.Channel))],

    statistics: (state) => {
      const total = state.reviews.length
      const approved = state.reviews.filter(r => r.ApprovalStatus?.IsApproved).length
      const avgRating = total
        ? state.reviews.reduce((sum, r) => sum + (r.Rating || 0), 0) / total
        : 0

      const propertyStats = {}
      state.reviews.forEach(review => {
        const name = review.PropertyName
        if (!propertyStats[name]) propertyStats[name] = { count: 0, totalRating: 0, approved: 0 }
        propertyStats[name].count++
        propertyStats[name].totalRating += review.Rating || 0
        if (review.ApprovalStatus?.IsApproved) propertyStats[name].approved++
      })

      return {
        total,
        approved,
        pending: total - approved,
        avgRating: Number(avgRating.toFixed(1)),
        propertyStats,
      }
    },
  },

  actions: {
    // ✅ FIXED: Made propertyId optional
    async fetchReviews(propertyId = null) {
      this.loading = true
      this.error = null

      try {
        const data = await reviewsAPI.fetchHostawayReviews(propertyId)
        console.log('API RESPONSE:', data)

        // Normalize API response: handle various backend formats
        let allReviews = []
        if (Array.isArray(data)) allReviews = data
        else if (Array.isArray(data.result)) allReviews = data.result
        else if (Array.isArray(data.reviews)) allReviews = data.reviews
        else if (Array.isArray(data.data)) allReviews = data.data
        else {
          console.error('Unexpected reviews API format:', data)
          this.reviews = []
          return []
        }

        // If fetching for a specific property, only update those reviews
        if (propertyId) {
          // Remove old reviews for this property and add new ones
          this.reviews = [
            ...this.reviews.filter(r => 
              r.PropertyID !== propertyId && 
              r.PropertyName !== propertyId &&
              r.ListingID !== propertyId
            ),
            ...allReviews
          ]
        } else {
          // Replace all reviews
          this.reviews = allReviews
        }

        return this.reviews
      } catch (error) {
        this.error = error.message || 'Unknown error'
        console.error('Error fetching reviews:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async toggleApproval(reviewId) {
      const review = this.reviews.find(r => r.ID === reviewId)
      if (!review) return

      const newStatus = !review.ApprovalStatus?.IsApproved

      try {
        await reviewsAPI.updateApproval(reviewId, newStatus)

        if (!review.ApprovalStatus) review.ApprovalStatus = {}

        review.ApprovalStatus.IsApproved = newStatus
        review.ApprovalStatus.IsRejected = !newStatus
      } catch (error) {
        console.error('Error updating approval:', error)
        throw error
      }
    },

    setFilter(filterType, value) {
      if (filterType in this.filters) {
        this.filters[filterType] = value
      }
    },

    setSortBy(sortType) {
      if (['date', 'rating'].includes(sortType)) this.sortBy = sortType
    },
  },
})