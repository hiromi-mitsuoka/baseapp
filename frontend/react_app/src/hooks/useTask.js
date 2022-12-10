/**
 * useTask
 *
 * @package hooks
 */
import { useState, useMemo } from "react";
import { useSelector, useDispatch } from "react-redux";
import { addTask, deleteTask } from "../store/task";

/**
 * useTask
 */
export const useTask = () => {
  /* store */
  // useSelector: actionがdispatchされると実行される
  // Allows you to extract data from the Redux store state, using a selector function.
  const taskList = useSelector((state) => state.task.tasks);
  // useDispatch: actionを実行して，storeへ変更を伝える
  const dispatch = useDispatch();

  /* local state */
  /* add input title */
  const [addInputValue, setAddInputValue] = useState("");
  /* 検索キーワード */
  const [searchKeyword, setSearchKeyword] = useState("");
  /* 表示用TaskList */
  const showTaskList = useMemo(() => {
    return taskList.filter((task) => {
      // 検索キーワードに前方一致
      const regexp = new RegExp("^" + searchKeyword, "i");
      return task.title.match(regexp);
    });
    // originTaskListとsearchKeywordの値が変更される度に，filterの検索処理が実行
    // useMemoを使用することで，結果が前回と同じならキャッシュを返却して処理は実行されない
  }, [taskList, searchKeyword])

  /* actions */
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
      // Task追加処理
      dispatch(addTask(addInputValue));
      // 入力値リセット
      setAddInputValue("");
    }
  };

  /**
   * Task削除処理
   * @param { number } targetId
   * @param { string } targetTitle
   */
  const handleDeleteTask = (targetId, targetTitle) => {
    if (window.confirm(`「${targetTitle}」のtaskを削除しますか？`)) {
      dispatch(deleteTask(targetId));
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