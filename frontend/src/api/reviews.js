import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

export const reviewsAPI = {
  // Fetch all normalized reviews from Hostaway
  async fetchHostawayReviews() {
    const response = await api.get('/reviews/hostaway')
    return response.data
  },

  // Update review approval status
  async updateApproval(reviewId, approved) {
    const response = await api.patch(`/reviews/${reviewId}/approve`, { approved })
    return response.data
  },

  // Get reviews by property
  async getReviewsByProperty(propertyId) {
    const response = await api.get(`/reviews/property/${propertyId}`)
    return response.data
  },

  // Get approved reviews for public display
  async getApprovedReviews(propertyId = null) {
    const url = propertyId 
      ? `/reviews/approved?propertyId=${propertyId}`
      : '/reviews/approved'
    const response = await api.get(url)
    return response.data
  },
}

export default api