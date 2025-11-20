<template>
  <div>
    <h1 class="text-4xl font-bold mb-2">{{ property.name }}</h1>
    <p class="text-gray-300">{{ property.location }}</p>
  </div>
</template>

<script setup>
defineProps({
  property: {
    type: Object,
    required: true,
  },
})
</script>