<template>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <div class="bg-white rounded-lg shadow p-6">
      <h3 class="text-lg font-semibold mb-4">Property Performance</h3>
      <div class="space-y-3">
        <div
          v-for="(stat, property) in propertyStats"
          :key="property"
          class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
        >
          <div class="flex-1">
            <p class="font-medium text-gray-900 truncate">{{ property }}</p>
            <p class="text-sm text-gray-600">{{ stat.count }} reviews</p>
          </div>
          <div class="flex items-center gap-2">
            <Star :size="16" class="text-yellow-500 fill-yellow-500" />
            <span class="font-bold text-lg">{{ stat.avgRating }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow p-6">
      <h3 class="text-lg font-semibold mb-4">Rating Distribution</h3>
      <div class="space-y-3">
        <div
          v-for="range in ratingDistribution"
          :key="range.label"
          class="flex items-center gap-3"
        >
          <span class="text-sm font-medium text-gray-700 w-16">{{ range.label }}</span>
          <div class="flex-1 bg-gray-200 rounded-full h-6 overflow-hidden">
            <div
              :style="{ width: `${range.percentage}%` }"
              :class="range.color"
              class="h-full flex items-center justify-end pr-2"
            >
              <span v-if="range.count > 0" class="text-xs font-medium text-white">
                {{ range.count }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Star } from 'lucide-vue-next'

const props = defineProps({
  reviews: {
    type: Array,
    required: true,
  },
})

const propertyStats = computed(() => {
  const stats = {}
  props.reviews.forEach(review => {
    if (!stats[review.listingName]) {
      stats[review.listingName] = {
        count: 0,
        totalRating: 0,
      }
    }
    stats[review.listingName].count++
    stats[review.listingName].totalRating += review.rating
  })

  Object.keys(stats).forEach(property => {
    stats[property].avgRating = (stats[property].totalRating / stats[property].count).toFixed(1)
  })

  return stats
})

const ratingDistribution = computed(() => {
  const distribution = [
    { label: '9-10', range: [9, 10], count: 0, color: 'bg-green-500' },
    { label: '7-8.9', range: [7, 8.9], count: 0, color: 'bg-blue-500' },
    { label: '5-6.9', range: [5, 6.9], count: 0, color: 'bg-yellow-500' },
    { label: '0-4.9', range: [0, 4.9], count: 0, color: 'bg-red-500' },
  ]

  props.reviews.forEach(review => {
    const rating = review.rating
    for (const dist of distribution) {
      if (rating >= dist.range[0] && rating <= dist.range[1]) {
        dist.count++
        break
      }
    }
  })

  const total = props.reviews.length
  distribution.forEach(dist => {
    dist.percentage = total > 0 ? (dist.count / total) * 100 : 0
  })

  return distribution
})
</script>