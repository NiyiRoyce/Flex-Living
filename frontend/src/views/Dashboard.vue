<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-6 py-4">
        <div class="flex justify-between items-center">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Flex Living Reviews Dashboard</h1>
            <p class="text-gray-600 text-sm mt-1">
              Manage and monitor property reviews across all channels
            </p>
          </div>
          <router-link
            to="/property/demo"
            class="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            <Eye :size="18" />
            Preview Public Page
          </router-link>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-6 py-8">
      <!-- Loading State -->
      <div v-if="reviewsStore.loading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p class="text-gray-600 mt-4">Loading reviews...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="reviewsStore.error" class="bg-red-50 border border-red-200 rounded-lg p-4">
        <p class="text-red-800">Error loading reviews: {{ reviewsStore.error }}</p>
      </div>

      <!-- Dashboard Content -->
      <template v-else>
        <!-- Stats Overview -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div class="bg-white rounded-lg shadow p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-gray-600 text-sm">Total Reviews</p>
                <p class="text-3xl font-bold text-gray-900 mt-1">
                  {{ reviewsStore.statistics.total }}
                </p>
              </div>
              <div class="bg-blue-100 p-3 rounded-lg">
                <Star class="text-blue-600" :size="24" />
              </div>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-gray-600 text-sm">Average Rating</p>
                <p class="text-3xl font-bold text-gray-900 mt-1">
                  {{ reviewsStore.statistics.avgRating }}
                </p>
              </div>
              <div class="bg-green-100 p-3 rounded-lg">
                <TrendingUp class="text-green-600" :size="24" />
              </div>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-gray-600 text-sm">Approved</p>
                <p class="text-3xl font-bold text-gray-900 mt-1">
                  {{ reviewsStore.statistics.approved }}
                </p>
              </div>
              <div class="bg-green-100 p-3 rounded-lg">
                <CheckCircle class="text-green-600" :size="24" />
              </div>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-gray-600 text-sm">Pending</p>
                <p class="text-3xl font-bold text-gray-900 mt-1">
                  {{ reviewsStore.statistics.pending }}
                </p>
              </div>
              <div class="bg-orange-100 p-3 rounded-lg">
                <XCircle class="text-orange-600" :size="24" />
              </div>
            </div>
          </div>
        </div>

        <!-- Charts -->
        <RatingSummary :reviews="reviewsStore.reviews" class="mb-8" />

        <!-- Filters -->
        <div class="bg-white rounded-lg shadow p-6 mb-6">
          <div class="flex items-center gap-4 flex-wrap">
            <div class="flex items-center gap-2">
              <Filter :size="18" class="text-gray-600" />
              <span class="text-sm font-medium text-gray-700">Filters:</span>
            </div>

            <select
              v-model="reviewsStore.filters.property"
              class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="all">All Properties</option>
              <option v-for="prop in reviewsStore.properties" :key="prop" :value="prop">
                {{ prop }}
              </option>
            </select>

            <select
              v-model="reviewsStore.filters.channel"
              class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="all">All Channels</option>
              <option v-for="channel in reviewsStore.channels" :key="channel" :value="channel">
                {{ channel }}
              </option>
            </select>

            <select
              v-model="reviewsStore.sortBy"
              class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="date">Sort by Date</option>
              <option value="rating">Sort by Rating</option>
            </select>
          </div>
        </div>

        <!-- Reviews List -->
        <div class="space-y-4">
          <ReviewCard
            v-for="review in reviewsStore.filteredReviews"
            :key="review.id"
            :review="review"
            @toggle-approval="reviewsStore.toggleApproval"
          />
        </div>
      </template>
    </main>
  </div>
</template>

<script setup>
import { Eye, Star, TrendingUp, CheckCircle, XCircle, Filter } from 'lucide-vue-next'
import { useReviewsStore } from '@/store/reviews'
import ReviewCard from '@/components/ReviewCard.vue'
import RatingSummary from '@/components/RatingSummary.vue'

const reviewsStore = useReviewsStore()
</script>