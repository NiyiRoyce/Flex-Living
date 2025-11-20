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
    filteredReviews: (state) => {
      // Combine all filters for performance
      let filtered = state.reviews.filter(r => {
        const matchesProperty = state.filters.property === 'all' || r.PropertyName === state.filters.property
        const matchesChannel = state.filters.channel === 'all' || r.Channel === state.filters.channel
        const matchesRating = state.filters.rating === 'all' || (() => {
          const [min, max] = state.filters.rating.split('-').map(Number)
          return (r.Rating || 0) >= min && (r.Rating || 0) <= max
        })()
        return matchesProperty && matchesChannel && matchesRating
      })

      // Sorting
      filtered.sort((a, b) => {
        if (state.sortBy === 'date') {
          return new Date(b.SubmittedAt || 0).getTime() - new Date(a.SubmittedAt || 0).getTime()
        } else if (state.sortBy === 'rating') {
          return (b.Rating || 0) - (a.Rating || 0)
        }
        return 0
      })

      return filtered
    },

    approvedReviews: (state) => state.reviews.filter(r => r.ApprovalStatus?.IsApproved),

    properties: (state) => [...new Set(state.reviews.map(r => r.PropertyName))],

    channels: (state) => [...new Set(state.reviews.map(r => r.Channel))],

    statistics: (state) => {
      const total = state.reviews.length
      const approved = state.reviews.filter(r => r.ApprovalStatus?.IsApproved).length
      const avgRating = total > 0
        ? state.reviews.reduce((sum, r) => sum + (r.Rating || 0), 0) / total
        : 0

      const propertyStats = {}
      state.reviews.forEach(review => {
        const name = review.PropertyName
        if (!propertyStats[name]) {
          propertyStats[name] = { count: 0, totalRating: 0, approved: 0 }
        }
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
    async fetchReviews() {
      this.loading = true
      this.error = null

      try {
        const data = await reviewsAPI.fetchHostawayReviews()
        console.log("API RESPONSE:", data) // optional debug

        if (Array.isArray(data)) {
          this.reviews = data
        } else if (data && Array.isArray(data.result)) {
          this.reviews = data.result
        } else if (data && Array.isArray(data.reviews)) {
          this.reviews = data.reviews
        } else {
          console.error("Unexpected reviews API format:", data)
          this.reviews = []
        }

      } catch (error) {
        this.error = error.message || 'Unknown error'
        console.error('Error fetching reviews:', error)
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
      if (sortType === 'date' || sortType === 'rating') {
        this.sortBy = sortType
      }
    },
  },
})
