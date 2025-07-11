<template>
  <div>
    <h1 class="header">Todos</h1>
    <p class="doneCount">Total tasks completed: {{ doneCount }}</p>
    <div class="wrapper">
      <input id="queryFilter" class="icon" type="text" autocomplete="off" v-model="searchQuery" placeholder="Search tasks..." title="Filter tasks" />
    </div>
    <div class="wrapper">
      <label class="filter" for="show-dropped">Show dropped tasks
        <input type="checkbox" id="show-dropped" v-model="showDropped" />
      </label>
      <label class="filter" for="show-tentative">Show tentative tasks
        <input type="checkbox" id="show-tentative" v-model="showTentative" />
      </label>
      <label class="filter" for="show-low-priority">Show low priority tasks
        <input type="checkbox" id="show-low-priority" v-model="showLowPriority" />
      </label>
      <label class="filter" for="show-done">Show done tasks
        <input type="checkbox" id="show-done" v-model="showDone" />
      </label>
      <label class="filter" for="show-only-high-priority">Show only high priority tasks
        <input type="checkbox" id="show-only-high-priority" v-model="showOnlyHighPriority" />
      </label>
    </div>
    <div class="wrapper">
      <button class="btn" @click="fetchTodos">Refresh</button>
      <button class="btn" @click="expandAll">Expand All</button>
      <button class="btn" @click="collapseAll">Collapse All</button>
      <button class="btn" @click="getRandomTodo">Get Random Task</button>
      <button class="btn" @click="resetFilters">Reset filters</button>
    </div>
    <div class="random-todo-wrapper random-todo-container" v-if="randomTodo">
      <h2 class="header">Random Task</h2>
      <div class="random-todo-container">
        <TodoItem :item="randomTodo" />
        <button class="btn" @click="randomTodo = null">Clear</button>
      </div>
    </div>
    <div class="todo_list_container">
      <ul class="todo-list">
        <li v-for="(todo, index) in filteredTodos" :key="todo.id">
          <TodoItem :item="todo" :key="index" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import TodoItem, { Todo, TodoStatus, TodoPriority } from './TodoItem.vue';

// TODO: Auto-refresh via websocket + filewatcher
// TODO: global "Collapse" button
// TODO: sort either by date or by progress
// TODO: "show all" button
// TODO: error handling for API calls

const statuses = ["not started", "in progress", "done", "dropped"] as const;
const priorities = { [TodoPriority.LOW]: "Low", [TodoPriority.MEDIUM]: "Medium", [TodoPriority.HIGH]: "High" } as const;

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
  const priority: TodoPriority = todo.priority ?? TodoPriority.MEDIUM;
  const status = progress === 100 && todo.status != TodoStatus.DROPPED ? TodoStatus.DONE : todo.status;
  return {
    ...todo,
    status: status,
    statusText: statuses[status],
    progress: progress,
    tentative: todo.tentative ?? false,
    priority: priority,
    priorityText: todo.content !== undefined ? priorities[priority] : undefined,
    showSubtasks: expand,
    visible: true,
    subtasks: todo.subtasks?.map((x: Todo) => initExistingTodo(x, false)),
    date: todo.date == undefined ? "unknown date" : new Date(todo.date).toLocaleDateString(),
  };
};

type Filter = (todo: Todo) => boolean;
const combineFilters = (a: Filter, b: Filter): Filter => (todo: Todo) => a(todo) && b(todo);
const hideRecursively = (todos: Todo[], filter: Filter): void => todos.forEach((todo: Todo) => {
  todo.visible = filter(todo);
  hideRecursively(todo.subtasks || [], filter);
});

// originally modelled as a filter that was mutating the todo tree, but this was recursively
// triggering Vue's reactivity system and causing infinite loops
const checkHighPriorityRecursively = (todo: Todo, showDone: boolean): boolean => {
  if (!showDone && isDone(todo)) {
    todo.visible = false;
    return false;
  }

  const highPrioSubTask = todo.subtasks?.filter((todo: Todo) => checkHighPriorityRecursively(todo, showDone)).length || false;
  if (todo.priority === TodoPriority.HIGH || highPrioSubTask) {
    todo.visible = true;
    return true;
  }

  todo.visible = false;
  return false;
};

const checkInProgressRecursively = checkHighPriorityRecursively; // TODO

const expandSubTasks = (todo: Todo, setting: boolean): void => {
  todo.showSubtasks = setting;
  todo.subtasks?.forEach((todo: Todo) => expandSubTasks(todo, setting));
};

const hideRecursivelyShallow = (todos: Todo[], filter: Filter): void => todos.forEach((todo: Todo) => {
  console.log(`Checking visibility for todo: ${todo.content}`, todo.visible, filter(todo), todo.subtasks?.some(filter));
  if (todo.visible && filter(todo)) {
    return true;
  }

  todo.visible = getAllSubTasksRecursively(todo).some(filter);
  hideRecursivelyShallow(todo.subtasks || [], filter);
});

const constructSearchQuery = (searchQuery: string, todos: Todo[]): () => void => {
  if (!searchQuery) {
    return () => {};
  }

  const query = searchQuery.toLowerCase();
  const includes = (x: string | undefined) => x?.toLowerCase().includes(query) || false;
  const includes_todo = (x: Todo) => includes(x.content) || includes(x.source) || includes(x.date);
  return () => hideRecursivelyShallow(todos, includes_todo);
};

export default defineComponent({
  name: 'TodoTree',
  components: { TodoItem },
  computed: {
    filters(): Filter[] {
      const filters: Filter[] = [];
      if (!this.showDropped) { filters.push((todo: Todo) => todo.status !== TodoStatus.DROPPED); }
      if (!this.showTentative) { filters.push((todo: Todo) => !todo.tentative); }
      if (!this.showLowPriority) { filters.push((todo: Todo) => todo.priority !== TodoPriority.LOW); }
      if (!this.showDone) { filters.push((todo: Todo) => todo.status !== TodoStatus.DONE); }
      return filters;
    },
    filteredTodos(): Todo[] {
      const searchQuery = constructSearchQuery(this.searchQuery, this.todos);
      if (this.showOnlyHighPriority) {
        this.todos.forEach((todo: Todo) => checkHighPriorityRecursively(todo, this.showDone));
        searchQuery();
        return this.todos;
      }
      if (this.showOnlyInProgress) {
        this.todos.forEach((todo: Todo) => checkInProgressRecursively(todo, this.showDone));
        searchQuery();
        return this.todos;
      }
      const filter = this.filters.reduce(combineFilters, () => true);
      hideRecursively(this.todos, filter);
      searchQuery();
      return this.todos;
    },
    numberedTodoMap(): Map<number, Todo> {
      const map = new Map<number, Todo>();
      let counter = 0;
      const traverse = (todo: Todo) => {
        if (!todo.visible) {
          return;
        }
        todo.id = counter++;
        map.set(todo.id, todo);
        todo.subtasks?.forEach(traverse);
      };
      this.filteredTodos.forEach(traverse);
      return map;
    },
    todos(): Todo[] {
      return Array.from(this.todoMap.values());
    },
    doneCount(): number {
      return this.todos.map(countDone).reduce(sum, 0);
    }
  },
  data() {
    return {
      todoMap: new Map<string, Todo>(),
      randomTodo: null as Todo | null,
      showDropped: false as boolean,
      showTentative: false as boolean,
      showLowPriority: false as boolean,
      showDone: false as boolean,
      showOnlyHighPriority: false as boolean,
      showOnlyInProgress: false as boolean,
      searchQuery: '' as string
    };
  },
  methods: {
    async fetchTodos() {
      fetch('/api/todos')
        .then(res => res.json())
        .then(data => {
          const todos = data.map(initExistingTodo);
          todos.sort((a: Todo, b: Todo) => (b.progress as number) - (a.progress as number))
          this.todoMap.clear();
          todos.forEach((todo: Todo) => {
            this.todoMap.set(todo.source || todo.content as string, todo);
          });
        });
    },
    getRandomTodo() {
      const amount = this.numberedTodoMap.size;
      const randomIndex = Math.floor(Math.random() * amount);
      this.randomTodo = this.numberedTodoMap.get(randomIndex) ?? null;
      if (!this.randomTodo?.content) {
        this.getRandomTodo(); // ensure we get a valid todo with content
      }
    },
    resetFilters() {
      this.showDropped = false;
      this.showTentative = false;
      this.showLowPriority = false;
      this.showDone = false;
      this.showOnlyHighPriority = false;
      this.showOnlyInProgress = false;
    },
    expandAll() {
      this.todos.forEach((todo: Todo) => expandSubTasks(todo, true));
    },
    collapseAll() {
      this.todos.forEach((todo: Todo) => expandSubTasks(todo, false));
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
  align-items: center;
}

label:has(+#queryFilter) {
  margin-right: 0.5rem;
}

#queryFilter {
  padding: 0.25rem;
  border-radius: 5px;
  border: 1px solid #ccc;
  width: 30%;
  height: 2rem;
  background: url("find.png") no-repeat 2.5% 50%;
  background-size: 9%;
  text-indent: 11%;
}

.filter {
  margin-left: 1rem;
}

.doneCount {
  font-size: 1.2rem;
  color: #49f360e7;
  text-align: center;
  margin-top: 0.1rem;
  margin-bottom: 1rem;
  text-shadow: 1px 1px 20px #0770178e;
}

.random-todo-wrapper {
  border: 2px solid #ccc;
  border-radius: 5px;
  margin: 1rem auto;
  padding: 0.5rem;
}

.random-todo-wrapper > .header {
  text-wrap: nowrap;
  text-align: center;
  font-size: 1rem;
  margin-bottom: 0;
}

.random-todo-container {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.random-todo-container > .btn {
  margin-bottom: 0;
  min-width: unset;
  max-width: unset;
  width: unset;
}

</style>
