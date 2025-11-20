<template>
  <div
    :class="[
      'bg-white rounded-lg shadow p-6 transition-shadow hover:shadow-md',
      readonly ? '' : 'cursor-pointer'
    ]"
  >
    <!-- Header: Listing, Channel, Host Review -->
    <div class="flex justify-between items-start mb-4">
      <div class="flex-1">
        <div class="flex items-center gap-3 mb-2">
          <h3 class="font-semibold text-lg text-gray-900">{{ review.listingName }}</h3>
          <span class="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-xs font-medium">
            {{ review.channel }}
          </span>
          <span
            v-if="review.type === 'host-to-guest'"
            class="px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-xs font-medium"
          >
            Host Review
          </span>
        </div>

        <div class="flex items-center gap-4 text-sm text-gray-600">
          <span>{{ review.guestName }}</span>
          <span>•</span>
          <span>{{ formatDate(review.submittedAt) }}</span>
        </div>
      </div>

      <!-- Rating & Approval -->
      <div class="flex items-center gap-4">
        <div v-if="review.rating" class="text-right">
          <div class="flex items-center gap-1">
            <Star :size="20" class="text-yellow-500 fill-yellow-500" />
            <span class="text-2xl font-bold text-gray-900">{{ review.rating.toFixed(1) }}</span>
          </div>
        </div>

        <!-- Approval Toggle (Dashboard Only) -->
        <button
          v-if="!readonly"
          @click="$emit('toggle-approval', review.id)"
          :class="[
            'px-4 py-2 rounded-lg font-medium transition-colors whitespace-nowrap',
            review.approved
              ? 'bg-green-100 text-green-700 hover:bg-green-200'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
          ]"
        >
          {{ review.approved ? 'Approved' : 'Approve' }}
        </button>
      </div>
    </div>

    <!-- Review Content -->
    <p class="text-gray-700 mb-4">{{ review.publicReview }}</p>

    <!-- Review Categories -->
    <div v-if="review.reviewCategory && review.reviewCategory.length" class="flex gap-3 flex-wrap">
      <div
        v-for="(cat, idx) in review.reviewCategory"
        :key="idx"
        class="px-3 py-1 bg-gray-100 rounded-lg text-sm"
      >
        <span class="text-gray-600 capitalize">{{ formatCategory(cat.category) }}: </span>
        <span class="font-semibold text-gray-900">{{ cat.rating }}/10</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Star } from 'lucide-vue-next'
import { format } from 'date-fns'

defineProps({
  review: {
    type: Object,
    required: true,
  },
  readonly: {
    type: Boolean,
    default: false, // true for public page, false for dashboard
  }
})

defineEmits(['toggle-approval'])

const formatDate = (dateString) => {
  return format(new Date(dateString), 'MMM dd, yyyy')
}

const formatCategory = (category) => {
  return category.replace(/_/g, ' ')
}
</script>

<style scoped>
/* Optional: subtle approved border highlight for public page */
</style>
