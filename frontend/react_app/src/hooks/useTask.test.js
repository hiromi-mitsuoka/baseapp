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
});