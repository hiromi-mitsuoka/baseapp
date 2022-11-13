/**
 * InputForm
 *
 * @package components
 */
import styles from './style.module.scss'

/**
 * InputForm
 * @param {*} props
 * @returns
 */
export const InputForm = (props) => {
  // props
  const {inputValue, placeholder, handleChangeValue, handleKeyDown} = props;

  return (
    <input
      className={styles.input}
      type="text"
      placeholder={placeholder}
      value={inputValue}
      onChange={handleChangeValue}
      onKeyDown={handleKeyDown}
    />
  );
};