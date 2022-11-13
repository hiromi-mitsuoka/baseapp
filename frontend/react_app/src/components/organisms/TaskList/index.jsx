/**
 * TaskList
 *
 * @package components
 */
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
// styles
import styles from "./style.module.scss";

export const TaskList = (props) => {
  // props
  const { taskList, handleDeleteTask } = props;

  return (
    <ul className={styles.list}>
      {taskList.map((task) => (
        <li key={task.id} className={styles.task}>
          <span className={styles.task_title}>{task.title}</span>
          <div className={styles.far}>
            {/* https://www.digitalocean.com/community/tutorials/how-to-use-font-awesome-5-with-react-ja */}
            <FontAwesomeIcon
              icon={faTrashAlt}
              size="lg"
              onClick={() => handleDeleteTask(task.id, task.title)}
            />
          </div>
        </li>
      ))}
    </ul>
  );
};