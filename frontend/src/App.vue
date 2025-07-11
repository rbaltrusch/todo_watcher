<script setup lang="ts">
import TodoTree from './components/TodoTree.vue';
</script>

<template>
  <main>
    <div class="wrapper">
      <TodoTree />
    </div>
  </main>
  <footer>
    <a href="https://www.flaticon.com/free-icons/zoom" class="attributation" title="zoom icons">Zoom icons created by Ghozi Muhtarom - Flaticon</a>
  </footer>
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
