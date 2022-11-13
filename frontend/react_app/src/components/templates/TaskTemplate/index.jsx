/**
 * TaskTemplate
 *
 * @package components
 */
import { useState, useMemo } from "react";
import { InputForm } from "../../atoms/InputForm";
import { AddTask } from "../../organisms/AddTask";
import { TaskList } from "../../organisms/TaskList";
import { INIT_TASK_LIST, INIT_UNIQUE_ID } from "../../../constants/data.js";
import styles from "./style.module.scss";

export const TaskTemplate = () => {
  const [originTaskList, setOriginTaskList] = useState(INIT_TASK_LIST);
  const [addInputValue, setAddInputValue] = useState("");
  const [uniqueId, setUniqueId] = useState(INIT_UNIQUE_ID);
  const [searchKeyword, setSearchKeyword] = useState("");

  const showTaskList = useMemo(() => {
    return originTaskList.filter((task) => {
      // 検索キーワードに前方一致したTaskだけ一覧表示
      const regexp = new RegExp("^" + searchKeyword, "i");
      return task.title.match(regexp);
    });
    // originTaskListとsearchKeywordの値が変更される度に，filterの検索処理が実行
    // useMemoを使用することで，結果が前回と同じならキャッシュを返却して処理は実行されない
  }, [originTaskList, searchKeyword])

  const onChangeAddInputValue = (e) => setAddInputValue(e.target.value);

  /**
   * Task新規登録処理
   * @param {*} e
   */
  const handleAddTask = (e) => {
    if (e.key === "Enter" && addInputValue !== "") {
      const nextUniqueID = uniqueId + 1;

      const newTaskList = [
        ...originTaskList,
        {
          id: nextUniqueID,
          title: addInputValue,
        },
      ];
      setOriginTaskList(newTaskList);

      // 採番IDを更新
      setUniqueId(nextUniqueID);

      setAddInputValue("");
    }
  };

  /**
   * Task削除処理
   * @param {*} targetId
   * @param {*} targetTitle
   */
  const handleDeleteTask = (targetId, targetTitle) => {
    if (window.confirm(`「${targetTitle}」のTaskを削除しますか？`)) {
      const newTaskList = originTaskList.filter((task) => task.id !== targetId);
      setOriginTaskList(newTaskList);
    }
  };

  /**
   * 検索キーワード更新処理
   * @param {*} e
   */
  const handleChangeSearchKeyword = (e) => setSearchKeyword(e.target.value);

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Task List</h1>
      <section className={styles.common}>
        <AddTask
          addInputValue={addInputValue}
          onChangeTask={onChangeAddInputValue}
          handleAddTask={handleAddTask}
        />
      </section>
      <section className={styles.common}>
        <InputForm
          InputValue={searchKeyword}
          placeholder={"Search Keyword"}
          handleChangeValue={handleChangeSearchKeyword}
        />
      </section>
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