<template>
  <div class="todo-item" :class="{dropped: item.statusText === 'dropped'}">
    <div class="wrapper">
      <p class="flex-item date" v-if="item.date === 'unknown date' && item.source"></p>
      <p class="flex-item date" v-if="item.date !== 'unknown date'">{{ item.date }}</p>
      <p class="flex-item source" v-if="item.source">{{ item.source }}</p>
      <p class="flex-item content" v-if="item.content">{{ item.content }}</p>
      <p class="flex-item status" v-if="item.progress !== undefined" :class="item.statusText?.replaceAll(' ', '')">{{
        item.progress }}% done</p>
      <button class="flex-item btn" v-if="item.source" @click="openFile(item.source)">Open</button>
      <button class="flex-item btn" v-if="item.subtasks?.length" @click="item.showSubtasks = !item.showSubtasks">{{
        item.showSubtasks ?
          "Collapse" : "Expand" }}</button>
    </div>
    <div class="sub-items" v-if="item.subtasks?.length && item.showSubtasks">
      <TodoItem v-for="(child, index) in item.subtasks" :key="index" :item="child" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import type { PropType } from 'vue';

export enum TodoStatus {
  NOT_STARTED = 0,
  IN_PROGRESS = 1,
  DONE = 2,
  DROPPED = 3
}

export interface Todo {
  status: TodoStatus;
  statusText?: string; // set by frontend based on status
  progress?: number; // 0-100, set by frontend based on status and status of subtasks
  date?: string;
  source?: string;
  content?: string;
  subtasks?: Todo[];
  showSubtasks?: boolean; // set by frontend
}

export default defineComponent({
  name: 'TodoItem',
  props: {
    item: {
      type: Object as PropType<Todo>,
      required: true
    }
  },
  components: {
    TodoItem: defineComponent() // this is what enables recursion
  },
  methods: {
    async openFile(file: string) {
      fetch(`/api/open?file=${encodeURIComponent(file)}`)
        .then(res => res.json())
        .then(data => {
          console.log('File opened:', data);
        })
        .catch(err => {
          console.error('Error opening file:', err);
        });
    }
  }
});
</script>

<style scoped>
.todo-item {
  margin-left: 1rem;
  padding: 0.5rem;
  border-left: 2px solid #ccc;
  min-width: 50vw;
}

.wrapper {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  border-top: 1px solid #333;
  border-radius: 5%;
}

.flex-item {
  margin: 0.25rem 0;
}

.source {
  width: 50%;
}

.date {
  width: 10%;
  padding-left: 0.5rem;
}

.content {
  padding: 0.25rem 0.5rem;
}

.status {
  text-wrap: nowrap;
}

.btn {
  background-color: #76889b;
  color: white;
  border-radius: 5px;
  min-width: 10%;
  text-align: center;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
}

.notstarted {
  color: gray;
}

.inprogress {
  color: rgb(72, 72, 220);
}

.done {
  color: rgb(100, 211, 100);
}

.todo-item.dropped {
  opacity: 0.5;
}
</style>
