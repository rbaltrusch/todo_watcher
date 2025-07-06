<script setup lang="ts">
import TodoTree from './components/TodoTree.vue';
</script>

<template>
  <main>
    <div class="wrapper">
      <TodoTree />
    </div>
  </main>
</template>

<style scoped>
.wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.main {
  padding: 1rem;
  background-color: #f0f0f0;
}
</style>
