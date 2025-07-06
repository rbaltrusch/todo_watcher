<template>
  <div>
    <h1 class="header">Todos</h1>
    <div class="wrapper">
      <button class="btn" @click="fetchTodos">Refresh</button>
    </div>
    <p class="doneCount">Total tasks completed: {{ done }}</p>
    <div class="todo_list_container">
      <ul class="todo-list">
        <li v-for="(todo, index) in todos" :key="todo.id">
          <TodoItem :item="todo" :key="index" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import TodoItem, { Todo, TodoStatus } from './TodoItem.vue';

// TODO: "Get Random Item" button
// TODO: Auto-refresh via websocket + filewatcher
// TODO: "Collapse All" and "Expand All" buttons
// TODO: global "Collapse" button
// TODO: filter by status

const statuses = ["not started", "in progress", "done", "dropped"] as const;

const sum = (acc: number, cur: number) => acc + cur
const isDropped = (todo: Todo): boolean => todo.status === TodoStatus.DROPPED;
const isDone = (todo: Todo): boolean => todo.status === TodoStatus.DONE;

function getAllSubTasksRecursively(todo: Todo): Todo[] {
  if (!todo.subtasks?.length) {
    return [];
  }

  return todo.subtasks.flatMap((x: Todo) => [x, ...getAllSubTasksRecursively(x)]);
}

function determineProgress(todo: Todo): number {
  if (isDone(todo) || isDropped(todo)) {
    return 100;
  }

  const allSubTasks = getAllSubTasksRecursively(todo);
  if (!allSubTasks?.length) {
    return todo.status === TodoStatus.IN_PROGRESS ? 50 : 0;
  }

  const completed = allSubTasks.filter((x: Todo) => !isDropped(x)).map(determineProgress).reduce(sum, 0);
  return Math.round(completed / allSubTasks.length);
}

function countDone(todo: Todo): number {
  if (!todo.subtasks?.length) {
    return isDone(todo) ? 1 : 0;
  }

  return todo.subtasks.map(countDone).reduce(sum, 0) || 0;
}

function initExistingTodo(todo: Todo, expand: boolean = true): Todo {
  const progress = determineProgress(todo);
  const status = progress === 100 && todo.status != TodoStatus.DROPPED ? TodoStatus.DONE : todo.status;
  return {
    ...todo,
    status: status,
    statusText: statuses[status],
    progress: progress,
    showSubtasks: expand,
    subtasks: todo.subtasks?.map((x: Todo) => initExistingTodo(x, false)),
    date: todo.date == undefined ? "unknown date" : new Date(todo.date).toLocaleDateString(),
  };
};

export default defineComponent({
  name: 'TodoTree',
  components: { TodoItem },
  data() {
    return {
      todos: [] as Todo[],
      done: 0,
    };
  },
  methods: {
    async fetchTodos() {
      fetch('/api/todos')
        .then(res => res.json())
        .then(data => {
          const todos = data.map(initExistingTodo);
          todos.sort((a: Todo, b: Todo) => b.progress - a.progress)
          this.todos = todos;
          this.done = todos.map(countDone).reduce(sum, 0);
        });
    }
  },
  mounted() {
    this.fetchTodos();
  }
});
</script>

<style scoped>
.header {
  text-align: center;
  margin-bottom: 1rem;
}

.todo-list-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.todo-list {
  list-style-type: none;
  padding: 0;
}

.wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 1rem;
}

.btn {
  background-color: #76889b;
  color: white;
  border-radius: 5px;
  padding: 0.5rem 1rem;
  margin-bottom: 1rem;
}

.doneCount {
  color: rgb(100, 211, 100);
  text-align: center;
  margin-bottom: 1rem;
}
</style>
