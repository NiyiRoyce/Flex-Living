<template>
  <div class="min-h-screen bg-white text-gray-800">
    <!-- Header -->
    <header class="border-b border-gray-200 bg-white">
      <div class="max-w-7xl mx-auto px-6 py-6">
        <div class="flex justify-between items-center">
          <div>
            <h1 class="text-3xl font-semibold text-gray-900 tracking-tight">
              Flex Living
            </h1>
            <p class="text-gray-500 text-sm mt-1">
              Reviews Dashboard
            </p>
          </div>

          <router-link
            to="/property/demo"
            class="flex items-center gap-2 px-6 py-3 bg-emerald-500 text-white 
                   rounded-lg hover:bg-emerald-600 transition-all duration-200 
                   font-medium text-sm shadow-sm"
          >
            <Eye :size="18" />
            Preview Public Page
          </router-link>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-6 py-10">
      <!-- Loading -->
      <div v-if="reviewsStore.loading" class="text-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-2 border-gray-300 border-t-emerald-500 mx-auto"></div>
        <p class="text-gray-500 mt-4">Loading reviews...</p>
      </div>

      <!-- Error -->
      <div v-else-if="reviewsStore.error" class="bg-red-50 border border-red-200 rounded-xl p-6">
        <p class="text-red-700">Error loading reviews: {{ reviewsStore.error }}</p>
      </div>

      <!-- Dashboard -->
      <template v-else>
        <!-- Stats Overview -->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-10">
          <!-- Total Reviews -->
          <div class="bg-white rounded-xl p-6 border border-gray-200 hover:border-emerald-300 transition-all group shadow-sm">
            <div class="flex items-start justify-between mb-4">
              <div class="bg-emerald-100 p-3 rounded-lg group-hover:bg-emerald-200 transition-colors">
                <Star class="text-emerald-600" :size="24" />
              </div>
            </div>
            <p class="text-gray-500 text-xs font-medium uppercase tracking-wider mb-2">Total Reviews</p>
            <p class="text-4xl font-bold text-gray-900">{{ reviewsStore.statistics.total }}</p>
          </div>

          <!-- Average Rating -->
          <div class="bg-white rounded-xl p-6 border border-gray-200 hover:border-emerald-300 transition-all group shadow-sm">
            <div class="flex items-start justify-between mb-4">
              <div class="bg-emerald-100 p-3 rounded-lg group-hover:bg-emerald-200 transition-colors">
                <TrendingUp class="text-emerald-600" :size="24" />
              </div>
            </div>
            <p class="text-gray-500 text-xs font-medium uppercase tracking-wider mb-2">Average Rating</p>
            <div class="flex items-baseline gap-1">
              <p class="text-4xl font-bold text-gray-900">{{ reviewsStore.statistics.avgRating }}</p>
              <span class="text-gray-500 text-xl">/5</span>
            </div>
          </div>

          <!-- Approved -->
          <div class="bg-emerald-50 rounded-xl p-6 border border-emerald-100 hover:border-emerald-300 transition-all group shadow-sm">
            <div class="flex items-start justify-between mb-4">
              <div class="bg-emerald-100 p-3 rounded-lg group-hover:bg-emerald-200 transition-colors">
                <CheckCircle class="text-emerald-600" :size="24" />
              </div>
            </div>
            <p class="text-emerald-700 text-xs font-medium uppercase tracking-wider mb-2">Approved</p>
            <p class="text-4xl font-bold text-emerald-700">{{ reviewsStore.statistics.approved }}</p>
          </div>

          <!-- Pending -->
          <div class="bg-white rounded-xl p-6 border border-gray-200 hover:border-gray-300 transition-all group shadow-sm">
            <div class="flex items-start justify-between mb-4">
              <div class="bg-gray-100 p-3 rounded-lg group-hover:bg-gray-200 transition-colors">
                <XCircle class="text-gray-400" :size="24" />
              </div>
            </div>
            <p class="text-gray-500 text-xs font-medium uppercase tracking-wider mb-2">Pending</p>
            <p class="text-4xl font-bold text-gray-700">{{ reviewsStore.statistics.pending }}</p>
          </div>
        </div>

        <!-- Charts -->
        <div class="mb-10">
          <RatingSummary :reviews="reviewsStore.reviews" />
        </div>

        <!-- Filters -->
        <div class="bg-white rounded-xl p-6 mb-8 border border-gray-200 shadow-sm">
          <div class="flex items-center gap-4 flex-wrap">
            <div class="flex items-center gap-2 mr-2">
              <div class="bg-emerald-100 p-2 rounded-lg">
                <Filter :size="18" class="text-emerald-600" />
              </div>
              <span class="text-sm font-medium text-gray-900">Filters</span>
            </div>

            <select
              v-model="reviewsStore.filters.property"
              class="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 
                     focus:ring-emerald-500 focus:border-emerald-500 bg-white 
                     text-gray-700 text-sm hover:border-emerald-300 transition-colors"
            >
              <option value="all">All Properties</option>
              <option v-for="prop in reviewsStore.properties" :key="prop" :value="prop">{{ prop }}</option>
            </select>

            <select
              v-model="reviewsStore.filters.channel"
              class="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 
                     focus:ring-emerald-500 focus:border-emerald-500 bg-white 
                     text-gray-700 text-sm hover:border-emerald-300 transition-colors"
            >
              <option value="all">All Channels</option>
              <option v-for="channel in reviewsStore.channels" :key="channel" :value="channel">{{ channel }}</option>
            </select>

            <select
              v-model="reviewsStore.sortBy"
              class="px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 
                     focus:ring-emerald-500 focus:border-emerald-500 bg-white 
                     text-gray-700 text-sm hover:border-emerald-300 transition-colors"
            >
              <option value="date">Sort by Date</option>
              <option value="rating">Sort by Rating</option>
            </select>
          </div>
        </div>

        <!-- Reviews List -->
        <div class="space-y-5">
          <ReviewCard
            v-for="review in reviewsStore.filteredReviews"
            :key="review.id"
            :review="review"
            :readonly="false"
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
