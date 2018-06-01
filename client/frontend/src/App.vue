<template>
  <div id="app">
    <el-table :data="todos" style="width: 100%" class="todo-table">

      <el-table-column label="Text" prop="text">
        <template slot-scope="scope">
          <el-input placeholder="Текст" v-model="scope.row.text" @keyup.native="update(scope.row)"></el-input>
        </template>
      </el-table-column>

      <!-- <el-table-column prop="address" label="Address" :formatter="formatter"></el-table-column> -->

      <el-table-column prop="priority" label="Priority">
        <template slot-scope="scope">
          <el-select v-model="scope.row.priority" placeholder="Select">
            <el-option v-for="priority in priorities" :key="priority" :label="priority" :value="priority">
            </el-option>
          </el-select>
        </template>
      </el-table-column>

      <el-table-column label="Created at">
        <template slot-scope="scope">
          <i class="el-icon-time"></i>
          <span style="margin-left: 10px">{{ (new Date(scope.row.created_at)).toLocaleString() }}</span>
        </template>
      </el-table-column>

      <el-table-column label="Due date">
        <template slot-scope="scope">
          <el-date-picker v-model="scope.row.due_date" type="date" placeholder="Pick a due date"></el-date-picker>
        </template>
      </el-table-column>

      <el-table-column label="Operations">
        <template slot-scope="scope">
          <el-button size="mini" @click="complete(scope.$index, scope.row)">Complete</el-button>
          <el-button size="mini" type="danger" @click="remove(scope.$index, scope.row)">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>


    <div style="margin-top: 20px">
      <el-button size="medium" @click="addDialogVisible = true" type="primary">Add</el-button>
    </div>

    <el-dialog title="Add todo" :visible.sync="addDialogVisible" width="50%" :before-close="handleClose">
      <el-form ref="form" label-position="left" :model="form" label-width="120px">
        <el-form-item label="Todo text">
          <el-input v-model="form.text"></el-input>
        </el-form-item>

        <el-form-item label="Due date">
          <el-date-picker v-model="form.due_date" type="date" placeholder="Pick a due date"></el-date-picker>
        </el-form-item>

        <el-form-item label="Priority">
          <el-select v-model="form.priority" placeholder="Select">
            <el-option v-for="priority in priorities" :key="priority" :label="priority" :value="priority">
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="create(form)">Create</el-button>
          <el-button @click="addDialogVisible = false">Cancel</el-button>
        </el-form-item>
      </el-form>
      </span>
    </el-dialog>

    <el-dialog title="Auth" :visible.sync="dialogVisible" width="40%" :before-close="handleClose">
      <el-input class="auth-input" placeholder="ID" v-model="id" clearable></el-input>
      <el-input class="auth-input" placeholder="Token" v-model="token" clearable></el-input>

      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="dialogVisible = false">Confirm</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
  export default {
    data() {
      return {
        id: "",
        token: "",
        dialogVisible: false,
        addDialogVisible: false,
        todos: [],
        priorities: ["low", "medium", "high"],
        form: {
          text: '',
          due_date: '',
          priority: ''
        }
      };
    },
    created() {
      this.$socket.onopen = () => {
        this.$options.sockets.onmessage = (msg) => {
          event = JSON.parse(msg.data)
          console.log(event)

          if (event.type === "task_get") {
            this.todos = event.data;
          }

          if (event.type === "task_add") {
            this.todos.push(event.data);
          }

          console.log(this.todos)
        }

        this.$socket.sendObj({
          type: "auth",
          data: {
            id: 2,
            token: "680727359b47992f"
          }
        })

        this.$socket.sendObj({
          type: "task_get",
          data: {}
        })
      }
    },
    methods: {
      create(el) {
        const data = {
          text: el.text,
          priority: el.priority
        }

        if (el.due_date) data.due_date = el.due_date

        this.$socket.sendObj({
          type: "task_add",
          data
        })

        this.addDialogVisible = false;
      },
      update(row) {
        console.log(row)

        this.$socket.sendObj({
          type: "task_update",
          data: {
            id: row.id,
            text: row.text,
            priority: row.priority,
            due_date: row.due_date
          }
        })
      },
      remove(index, row) {
        this.todos.splice(index, 1)
        this.$socket.sendObj({
          type: "task_delete",
          data: {
            id: row.id
          }
        })
      },
      reorder() {
        this.$socket.sendObj({
          type: "task_reorder",
          data: [5, 4, 3, 2]
        })
      },
      complete(index, row) {
        this.$socket.sendObj({
          type: "task_complete",
          data: {
            id: row.id
          }
        })
      },
    }
  }
</script>

<style>
  #app {
    font-family: Helvetica, sans-serif;
    text-align: center;
  }

  .auth-input {
    width: 50%;
    margin-bottom: 10px;
  }

  .todo-table th {
    text-align: center;
  }
</style>