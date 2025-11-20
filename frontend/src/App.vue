<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useReviewsStore } from '@/store/reviews'

const reviewsStore = useReviewsStore()

onMounted(() => {
  reviewsStore.fetchReviews()
})
</script>