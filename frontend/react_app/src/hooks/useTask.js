/**
 * useTask
 *
 * @package hooks
 */
import { useState, useMemo } from "react";
import { INIT_TASK_LIST, INIT_UNIQUE_ID } from "../constants/data";

/**
 * useTask
 */
export const useTask = () => {
  // states
  const [originTaskList, setOriginTaskList] = useState(INIT_TASK_LIST);
  const [addInputValue, setAddInputValue] = useState("");
  const [uniqueId, setUniqueId] = useState(INIT_UNIQUE_ID);
  const [searchKeyword, setSearchKeyword] = useState("");

  const showTaskList = useMemo(() => {
    return originTaskList.filter((task) => {
      // 検索キーワードに前方一致
      const regexp = new RegExp("^" + searchKeyword, "i");
      return task.title.match(regexp);
    });
    // originTaskListとsearchKeywordの値が変更される度に，filterの検索処理が実行
    // useMemoを使用することで，結果が前回と同じならキャッシュを返却して処理は実行されない
  }, [originTaskList, searchKeyword])

  // actions
  /**
   * addInputValueの変更処理
   * @param {e} e
   */
  const onChangeAddInputValue = (e) => setAddInputValue(e.target.value);

  /**
   * Task新規登録処理
   * @params {*} e
   */
  const handleAddTask = (e) => {
    if (e.key === "Enter" && addInputValue !== "") {
      const nextUniqueId = uniqueId + 1;

      const newTaskList = [
        ...originTaskList,
        {
          id: nextUniqueId,
          title: addInputValue,
        },
      ];
      setOriginTaskList(newTaskList);

      // 採番IDを更新
      setUniqueId(nextUniqueId);
      setAddInputValue("");
    }
  };

  /**
   * Task削除処理
   * @param {*} targetId
   * @param {*} targetTitle
   */
  const handleDeleteTask = (targetId, targetTitle) => {
    if (window.confirm(`「${targetTitle}」のtaskを削除しますか？`)) {
      const newTaskList = originTaskList.filter((task) => task.id !== targetId);
      setOriginTaskList(newTaskList);
    }
  };

  /**
   * 検索キーワード更新処理
   * @param {*} e
   */
  const handleChangeSearchKeyword = (e) => setSearchKeyword(e.target.value);

  const states = {
    addInputValue,
    searchKeyword,
    showTaskList,
  };

  const actions = {
    onChangeAddInputValue,
    handleAddTask,
    handleDeleteTask,
    handleChangeSearchKeyword,
  };

  return [states, actions];
};