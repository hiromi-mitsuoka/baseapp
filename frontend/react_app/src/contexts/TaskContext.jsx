/**
 * TaskContext
 *
 * @package contexts
 */
import { useContext, createContext } from "react";
import { useTask } from "../hooks/useTask.js";

/**
 * TaskContext
 */
const TaskContext = createContest({});

/**
 * TaskProvider
 * @param children
 * @constructor
 */
export const TaskProvider = ({ children }) => {
  // カスタムフックから，状態とロジックを呼び出してコンテキストプロバイダーにあてがう
  const {
    addInputValue,
    searchKeyword,
    showTaskList,
    onChangeAddInputValue,
    handleAddTask,
    handleDeleteTask,
    handleChangeSearchKeyword,
  } = useTask();

  return (
    <TaskContext.Provider
      value={{
        addInputValue,
        searchKeyword,
        showTaskList,
        onChangeAddInputValue,
        handleAddTask,
        handleDeleteTask,
        handleChangeSearchKeyword,
      }}
    >
      {children}
    </TaskContext.Provider>
  );
};

/**
 * useTaskContext
 */
export const useTaskContext = () => useContext(TaskContext);