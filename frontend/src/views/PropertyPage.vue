<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Hero Section -->
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-6 py-12">
        <button @click="goBack" class="flex items-center gap-2 text-gray-600 hover:text-emerald-600 transition-colors mb-6 group">
          <ArrowLeft :size="20" class="group-hover:-translate-x-1 transition-transform" />
          <span class="font-medium">Back to Dashboard</span>
        </button>

        <div class="flex items-center justify-between mb-6">
          <div>
            <h1 class="text-4xl font-bold text-gray-900 mb-2">{{ propertyName }}</h1>
            <p class="text-gray-600">Flexible living spaces for modern lifestyles</p>
          </div>

          <div v-if="approvedReviewsList.length > 0" class="flex items-center gap-2 bg-emerald-50 px-6 py-4 rounded-xl border border-emerald-200">
            <Star :size="24" class="text-emerald-600 fill-emerald-600" />
            <div>
              <p class="text-3xl font-bold text-emerald-700">{{ averageRating }}</p>
              <p class="text-sm text-emerald-600">{{ approvedReviewsList.length }} {{ approvedReviewsList.length === 1 ? 'review' : 'reviews' }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Reviews Section -->
    <div class="max-w-7xl mx-auto px-6 py-12">
      <section v-if="approvedReviewsList.length > 0">
        <h2 class="text-3xl font-bold text-gray-900 mb-6">Guest Reviews</h2>

        <div class="space-y-6">
          <div v-for="review in approvedReviewsList" :key="review.id" class="bg-white rounded-xl p-6 border border-gray-200 shadow-sm hover:shadow-md transition-shadow">
            <div class="flex items-start justify-between mb-4">
              <div class="flex items-start gap-4">
                <div class="w-12 h-12 bg-emerald-100 rounded-full flex items-center justify-center flex-shrink-0">
                  <User :size="24" class="text-emerald-600" />
                </div>
                <div>
                  <h3 class="font-semibold text-gray-900">{{ review.guestName || 'Anonymous Guest' }}</h3>
                  <p class="text-sm text-gray-500">{{ formatDate(review.submittedAt) }}</p>
                </div>
              </div>

              <div class="flex items-center gap-1">
                <Star v-for="i in 5" :key="i" :size="18"
                  :class="i <= Math.round(review.rating || 0) ? 'text-emerald-500 fill-emerald-500' : 'text-gray-300'" />
              </div>
            </div>

            <p class="text-gray-700 leading-relaxed mb-3">{{ review.publicReview }}</p>

            <div class="flex items-center gap-4 text-sm text-gray-500">
              <div v-if="review.channel" class="flex items-center gap-1">
                <Globe :size="16" />
                <span>{{ review.channel }}</span>
              </div>
              <div class="flex items-center gap-1 text-emerald-600">
                <CheckCircle :size="16" />
                <span>Verified Stay</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section v-else class="text-center py-16">
        <div class="bg-white rounded-xl p-12 border border-gray-200 shadow-sm">
          <div class="bg-gray-100 w-20 h-20 rounded-full flex items-center justify-center mx-auto mb-4">
            <MessageSquare :size="36" class="text-gray-400" />
          </div>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">No Reviews Yet</h3>
          <p class="text-gray-600">Be the first to share your experience at this property!</p>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useReviewsStore } from '@/store/reviews'
import { Star, MessageSquare, User, CheckCircle, Globe, ArrowLeft } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const reviewsStore = useReviewsStore()

const propertyId = computed(() => route.params.id || 'demo')

const propertyName = computed(() => {
  const properties = { demo: 'Flex Living Downtown', '123': 'Flex Living Uptown' }
  return properties[propertyId.value] || 'Flex Living Property'
})

const approvedReviewsList = computed(() => {
  return reviewsStore.approvedReviews(propertyId.value).map(r => ({
    id: r.id,
    guestName: r.guestName || r.GuestName || 'Anonymous Guest',
    rating: r.rating || r.Rating || 0,
    publicReview: r.publicReview || r.ReviewText || r.Comment || '',
    submittedAt: r.submittedAt || r.SubmittedAt || '',
    channel: r.channel || r.Channel || '',
  }))
})

const averageRating = computed(() => {
  if (!approvedReviewsList.value.length) return 0
  const sum = approvedReviewsList.value.reduce((acc, r) => acc + (r.rating || 0), 0)
  return (sum / approvedReviewsList.value.length).toFixed(1)
})

const formatDate = (date) => {
  if (!date) return 'Recently'
  return new Date(date).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
}

const goBack = () => router.push('/')

onMounted(async () => {
  await reviewsStore.fetchReviews(propertyId.value)
})
</script>
