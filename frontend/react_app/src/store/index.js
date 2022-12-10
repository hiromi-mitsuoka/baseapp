import { configureStore } from "@reduxjs/toolkit";
import taskListReducer from "./task";

export default configureStore({
  reducer: {
    task: taskListReducer,
  },
});