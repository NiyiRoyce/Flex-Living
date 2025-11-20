<template>
  <div class="min-h-screen bg-white">
    <!-- Header -->
    <header class="bg-gray-900 text-white">
      <div class="max-w-6xl mx-auto px-6 py-4">
        <router-link to="/" class="text-sm text-gray-400 hover:text-white mb-4 inline-block">
          ← Back to Dashboard
        </router-link>
        <PropertyHeader :property="propertyData" />
      </div>
    </header>

    <!-- Property Details Placeholder -->
    <div class="max-w-6xl mx-auto px-6 py-12">
      <div class="grid grid-cols-3 gap-4 mb-12">
        <div class="col-span-2 bg-gray-200 h-96 rounded-lg flex items-center justify-center text-gray-500">
          Property Images Gallery
        </div>
        <div class="space-y-4">
          <div class="bg-gray-200 h-44 rounded-lg"></div>
          <div class="bg-gray-200 h-44 rounded-lg"></div>
        </div>
      </div>

      <!-- Amenities & Details -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <div class="bg-gray-50 p-6 rounded-lg">
          <h3 class="font-semibold mb-2">Bedrooms</h3>
          <p class="text-2xl font-bold">2</p>
        </div>
        <div class="bg-gray-50 p-6 rounded-lg">
          <h3 class="font-semibold mb-2">Bathrooms</h3>
          <p class="text-2xl font-bold">1</p>
        </div>
        <div class="bg-gray-50 p-6 rounded-lg">
          <h3 class="font-semibold mb-2">Max Guests</h3>
          <p class="text-2xl font-bold">4</p>
        </div>
      </div>

      <!-- Reviews Section -->
      <div class="border-t pt-12">
        <div class="flex items-center justify-between mb-8">
          <h2 class="text-3xl font-bold text-gray-900">Guest Reviews</h2>
          <div v-if="approvedReviews.length" class="flex items-center gap-2">
            <Star :size="32" class="text-yellow-500 fill-yellow-500" />
            <span class="text-4xl font-bold text-gray-900">{{ averageRating }}</span>
            <span class="text-gray-600">({{ approvedReviews.length }} reviews)</span>
          </div>
        </div>

        <div v-if="approvedReviews.length === 0" class="text-center py-12 text-gray-500">
          No approved reviews yet
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div
            v-for="review in approvedReviews"
            :key="review.id"
            class="border rounded-lg p-6 hover:shadow-lg transition-shadow"
          >
            <div class="flex items-center justify-between mb-3">
              <div>
                <p class="font-semibold text-gray-900">{{ review.guestName }}</p>
                <p class="text-sm text-gray-500">{{ formatDate(review.submittedAt) }}</p>
              </div>
              <div class="flex items-center gap-1">
                <Star :size="18" class="text-yellow-500 fill-yellow-500" />
                <span class="font-bold text-lg">{{ review.rating.toFixed(1) }}</span>
              </div>
            </div>
            <p class="text-gray-700">{{ review.publicReview }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Star } from 'lucide-vue-next'
import { useReviewsStore } from '@/store/reviews'
import PropertyHeader from '@/components/PropertyHeader.vue'
import { format } from 'date-fns'

const reviewsStore = useReviewsStore()

const approvedReviews = computed(() => reviewsStore.approvedReviews)

const averageRating = computed(() => {
  if (approvedReviews.value.length === 0) return '0.0'
  const sum = approvedReviews.value.reduce((acc, r) => acc + r.rating, 0)
  return (sum / approvedReviews.value.length).toFixed(1)
})

const propertyData = computed(() => ({
  name: '2B N1 A - 29 Shoreditch Heights',
  location: 'London, United Kingdom',
}))

const formatDate = (dateString) => {
  return format(new Date(dateString), 'MMMM yyyy')
}
</script>