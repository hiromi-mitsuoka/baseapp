/**
 * AddTask
 *
 * @package components
 */
import { InputForm } from "../../atoms/InputForm";
import styles from "./style.module.scss";

/**
 * AddTask
 * @param {*} props
 * @returns
 */
export const AddTask = (props) => {
  // props
  const { addInputValue, onChangeTask, handleAddTask } = props;

  return (
    <>
      <h2 className={styles.subTitle}>{"ADD TASK"}</h2>
      <InputForm
        inputValue={addInputValue}
        placeholder={"New Task"}
        handleChangeValue={onChangeTask}
        handleKeyDown={handleAddTask}
      />
    </>
  );
};