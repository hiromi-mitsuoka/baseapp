import { createSlice } from "@reduxjs/toolkit";
import { INIT_TASK_LIST, INIT_UNIQUE_ID } from "../constants/data";

export const taskListSlice = createSlice({
  name: "task",
  initialState: {
    tasks: INIT_TASK_LIST,
    uniqueId: INIT_UNIQUE_ID,
  },
  // Reducer: actionとstateから，新しいstateを作成して返すメソッド
  // NOTE: 引数のstateを更新することはせず，新しいstateのオブジェクトを作成して返す
  reducers: {
    addTask: (state, action) => {
      const nextUniqueId = state.uniqueId + 1;
      state.tasks = [
        ...state.tasks,
        {
          id: nextUniqueId,
          title: action.payload,
        },
      ];
      state.uniqueId = nextUniqueId;
    },
    deleteTask: (state, action) => {
      state.tasks = state.tasks.filter((task) => task.id !== action.payload);
    },
  },
});

export const { addTask, deleteTask } = taskListSlice.actions;

export default taskListSlice.reducer;