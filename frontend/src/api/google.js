import api from './reviews'

export const googleAPI = {
  // Explore Google Places API integration
  async searchPlace(placeName) {
    const response = await api.get('/google/search', {
      params: { query: placeName }
    })
    return response.data
  },

  async getPlaceReviews(placeId) {
    const response = await api.get(`/google/reviews/${placeId}`)
    return response.data
  },
}