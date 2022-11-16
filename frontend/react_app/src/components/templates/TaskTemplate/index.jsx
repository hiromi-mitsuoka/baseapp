/**
 * TaskTemplate
 *
 * @package components
 */
import { InputForm } from "../../atoms/InputForm";
import { AddTask } from "../../organisms/AddTask";
import { TaskList } from "../../organisms/TaskList";
import styles from "./style.module.scss";
import { useTask } from "../../../hooks/useTask";

/**
 * TaskTemplate
 * @returns {JSX.Element}
 * @constructor
 */
export const TaskTemplate = () => {
  // カスタムフックから状態とロジックを呼び出して，コンポーネントにあてがう
  const [
    { addInputValue, searchKeyword, showTaskList },
    {
      onChangeAddInputValue,
      handleAddTask,
      handleDeleteTask,
      handleChangeSearchKeyword,
    },
  ] = useTask();

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Task List</h1>
      {/* Task追加エリア */}
      <section className={styles.common}>
        <AddTask
          addInputValue={addInputValue}
          onChangeTask={onChangeAddInputValue}
          handleAddTask={handleAddTask}
        />
      </section>
      {/* Task検索フォームエリア */}
      <section className={styles.common}>
        <InputForm
          inputValue={searchKeyword}
          placeholder={"Search Keyword"}
          handleChangeValue={handleChangeSearchKeyword}
        />
      </section>
      {/* Taskリスト一覧表示 */}
      <section className={styles.common}>
        {showTaskList.length > 0 && (
          <TaskList
            taskList={showTaskList}
            handleDeleteTask={handleDeleteTask}
          />
        )}
      </section>
    </div>
  );
};