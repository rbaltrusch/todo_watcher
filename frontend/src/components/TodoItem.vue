<template>
  <div class="todo-item" v-if="item.visible" :class="{ dropped: item.statusText === 'dropped', outer: item.source }">
    <div class="wrapper" :class="[{ tentative: item.tentative }, item.priorityText?.toLowerCase() + '-priority']">
      <p class="flex-item date" v-if="item.date === 'unknown date' && item.source"></p>
      <p class="flex-item date" v-if="item.date !== 'unknown date'">{{ item.date }}</p>
      <p class="flex-item source" v-if="item.source">{{ item.source }}</p>
      <p class="flex-item content" v-if="item.content">{{ item.content + (item.tentative ? '?' : '') }}</p>
      <div class="flex-item wrapper status-priority">
        <p class="flex-item priority" v-if="item.content && item.priorityText?.toLowerCase() !== 'medium'">{{
          item.priorityText }} priority</p>
        <p class="flex-item priority" v-if="item.priorityText?.toLowerCase() === 'medium'"></p> <!-- kind of a layout hack to have status always at end of right hand side -->
        <p class="flex-item status" v-if="item.progress !== undefined" :class="item.statusText?.replaceAll(' ', '')">{{
          item.progress }}% done</p>
      </div>
      <button class="flex-item btn" v-if="item.source" @click="openFile(item.source)">Open</button>
      <button class="flex-item btn" v-if="item.subtasks?.filter(x => x.visible).length" @click="item.showSubtasks = !item.showSubtasks">{{
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

export enum TodoPriority {
  LOW = -1,
  MEDIUM = 0,
  HIGH = 1
}

export type Todo = {
  status: TodoStatus;
  statusText?: string; // set by frontend based on status
  progress?: number; // 0-100, set by frontend based on status and status of subtasks
  priority?: TodoPriority;
  priorityText?: string; // set by frontend based on priority
  tentative?: boolean;
  date?: string;
  source?: string;
  content?: string;
  subtasks?: Todo[];
  visible?: boolean; // set by frontend
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

.todo-item:not(.outer):hover {
  transform: scale(1.015);
  border-radius: 5px;
  border: 1px solid #33373c;
  background: linear-gradient(to right, rgba(255, 255, 255, 0.025), rgba(255, 255, 255, 0.05));
  font-size: 1.05rem;
  color: #eeeeeec4;
  transition: transform 0.2s ease-in-out, background-color 0.2s ease-in-out , color 0.2s ease-in-out;
}

.todo-item > .wrapper {
  border-top: 1px solid #333;
}

.todo-item:not(.outer):hover > .wrapper:not(.high-priority) {
  border-top: none;
}

.random-todo-container > .todo-item:hover {
  transform: scale(1.005, 1) !important;
}

.wrapper {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  border-radius: 5%;
}

.wrapper.status-priority {
  width: 30%;
}

.flex-item {
  margin: 0.25rem 0;
  padding-left: 0.25rem;
  padding-right: 0.25rem;
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

.priority {
  width: 10%;
  text-align: center;
  text-wrap: nowrap;
}

.btn {
  padding: 0.25rem 0.5rem;
  margin: 0.25rem 0.25rem;
  max-height: unset;
  min-height: unset;
}

.notstarted {
  color: gray;
}

.inprogress {
  color: rgb(72, 72, 220);
  text-shadow: 1px 1px 4px rgba(72, 72, 220, 0.5);
}

.done {
  color: #49f360;
  text-shadow: 1px 1px 4px #077016;
}

.todo-item.dropped {
  opacity: 0.5;
}

.wrapper.tentative {
  opacity: 0.85;
  color: rgb(65, 66, 63);
  border-radius: 3px;
}

.wrapper.low-priority {
  opacity: 0.5;
  /* background: linear-gradient(to right, rgb(100, 106, 88), rgb(112, 120, 95)); */
  background: linear-gradient(45deg, var(--color-background), rgba(112, 120, 95, 0.1));
  border-radius: 3px;
}

.wrapper.high-priority {
  border: 1px solid #f3495d;
  box-shadow: 0 0 8px rgba(243, 73, 93, 0.4);
  filter: brightness(1.1);
  border-radius: 3px;
}

.wrapper.high-priority:has(> .status-priority > .status.done) {
  opacity: 0.8;
  border: 1px solid #4cd964;
  box-shadow: 0 0 8px #285e30;
  border-radius: 3px;
}
</style>
