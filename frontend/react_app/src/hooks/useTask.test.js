import { renderHook, act } from "@testing-library/react";
import { useTask } from "./useTask.js";
import { INIT_TASK_LIST } from "../constants/data.js";

describe("[Hooksテスト] useApp test", () => {
  describe("[関数テスト] onChangeAddInputValue", () => {
    test("[正常系] addInputValueを更新できること", () => {
      const expectValue = "test";
      const eventObject = {
        target: {
          value: expectValue,
        },
      };
      // hooks呼び出し
      // renderHook: カスタムフックを呼ぶための関数．カスタムフックはReact.FCの中以外で呼ぶことができず，そのへんをいい感じにしてくれる
      // https://zenn.dev/bom_shibuya/articles/5c3ae7745c5e94#%E2%96%A0-renderhook
      const { result } = renderHook(() => useTask());
      // current: renderHookを読んだカスタムフックの返り値の現在の値が格納されている
      expect(result.current[0].addInputValue).toBe("");
      // hooks関数の実行
      act(() => result.current[1].onChangeAddInputValue(eventObject));
      expect(result.current[0].addInputValue).toBe(expectValue);
    });
  });

  describe("[関数テスト] handleAddTask", () => {
    // 予測値
    let expectTaskList = [];
    // 引数
    let eventObject = {
      target: {
        value: "test",
      },
      key: "Enter",
    };
    /**
     * beforeEach
     * test関数が実行される前に毎回実行する
     */
    beforeEach(() => {
      // 引数の初期化
      eventObject = {
        target: {
          value: "test",
        },
        key: "Enter",
      };
    });

    test("[正常系] taskList, uniqueIdが更新されること, addInputValueがリセットされること", () => {
      // 予測値
      const expectTaskTitle = "Task3";
      expectTaskList = INIT_TASK_LIST.concat({
        id: 3,
        title: expectTaskTitle,
      });
      // 引数
      eventObject.target.value = expectTaskTitle;

      // hooks呼び出し
      const { result } = renderHook(() => useTask());
      expect(result.current[0].addInputValue).toBe("");
      // hooks関数の実行（addInputValueを更新）
      act(() => result.current[1].onChangeAddInputValue(eventObject));
      expect(result.current[0].addInputValue).toBe(expectTaskTitle);

      // hooks関数の実行: handleAddTaskの実行
      act(() => result.current[1].handleAddTask(eventObject));
      // 表示用TaskListが予想通り更新されたこと
      expect(result.current[0].showTaskList).toEqual(expectTaskList);
      // 入力値（addInputValue）がリセットされたこと
      expect(result.current[0].addInputValue).toBe("");
    });

    test("[正常系] エンターキーを押していない場合，処理が発生しないこと", () => {
      // 予測値
      const expectTaskTitle = "Task4"
      expectTaskList = INIT_TASK_LIST.concat({
        id: 3,
        title: expectTaskTitle,
      });
      // 引数
      eventObject.target.value = expectTaskTitle;
      eventObject.key = ""; // NOTE: エンターキーを押さない

      // hooks呼び出し
      const { result } = renderHook(() => useTask());
      expect(result.current[0].addInputValue).toBe("");
      // hooks関数の実行（addInputValueを更新）
      act(() => result.current[1].onChangeAddInputValue(eventObject));
      expect(result.current[0].addInputValue).toBe(expectTaskTitle);

      // hooks関数の実行: handleAddTaskの実行
      act(() => result.current[1].handleAddTask(eventObject))
      // 表示用TaskListが予想通り更新されないこと
      expect(result.current[0].showTaskList).not.toEqual(expectTaskList);
      // 入力値（addInputValue）がリセットされない
      expect(result.current[0].addInputValue).not.toBe("");
    });

    test("[正常系] 入力値がない場合，処理が発生しないこと", () => {
      // 予測値
      const expectTaskTitle = "Task5";
      expectTaskList = INIT_TASK_LIST.concat({
        id: 3,
        title: expectTaskTitle,
      });
      // 引数
      eventObject.target.value = ""
      eventObject.key = "";
      // hooks呼び出し
      const { result } = renderHook(() => useTask());
      expect(result.current[0].addInputValue).toBe("");
      // hooks関数の実行（addInputValueを更新）
      act(() => result.current[1].onChangeAddInputValue(eventObject));
      expect(result.current[0].addInputValue).toBe("");
      // hooks関数の実行: handleAddTaskの実行
      act(() => result.current[1].handleAddTask(eventObject));
      // 表示用TaskListが予想通り更新されないこと
      expect(result.current[0].showTaskList).not.toEqual(expectTaskList);
    });

    // test("[正常系] 検索キーワードがある場合", () => {

    // });

    describe("[関数テスト] handleDeleteTask", () => {
      // // 予測値
      // let expectTaskList = [];

      // beforeEach(() => {
      //   // 予測値を初期化
      //   expectTaskList = [];
      // });

      // test("[正常系] taskが削除されること", () => {

      // });

      // test("[正常系] confirmでキャンセルをクリックした場合, taskが削除されないこと", () => {

      // });

      // test("[正常系] 検索キーワードがある場合", () => {

      // });
    });
  });

  describe("[関数テスト] handleSearchTask", () => {
    // test("[正常系] 検索ワードがある場合，検索された結果が反映される", () => {

    // });

    // test("[正常系] 検索ワードがない場合, 元のTaskListが反映される", () => {

    // });
  });
});