CREATE TABLE `users`
(
  -- http://dbinfo.sakura.ne.jp/?contents_id=34
  -- UNSIGNED: 負数が使えなくなる代わりにその分扱える正の値が増える
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの識別子',
  `name` VARCHAR(20) NOT NULL COMMENT 'ユーザーの名前',
  `password` VARCHAR(80) NOT NULL COMMENT 'パスワードハッシュ',
  `role` VARCHAR(80) NOT NULL COMMENT 'ロール',
  `created` DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
  `modified` DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
  PRIMARY KEY (`id`),
  -- https://qiita.com/kiyodori/items/f66a545a47dc59dd8839
  -- BTREE: B-Treeアルゴリズムを使用
  UNIQUE KEY `uix_name` (`name`) USING BTREE
  -- https://dev.mysql.com/doc/refman/5.6/ja/innodb-default-se.html
  -- InnoDB: デフォルトのMysqlのストレージエンジン
  -- InnoDB を使用すれば、デフォルト設定でもユーザーが RDBMS から期待する利点 (ACID トランザクション、参照整合性、およびクラッシュリカバリ) が得られます
  -- https://dev.mysql.com/doc/refman/5.6/ja/charset-unicode-utf8mb4.html#:~:text=8%20Unicode%20%E3%82%A8%E3%83%B3%E3%82%B3%E3%83%BC%E3%83%87%E3%82%A3%E3%83%B3%E3%82%B0)-,10.1.10.7%20utf8mb4%20%E6%96%87%E5%AD%97%E3%82%BB%E3%83%83%E3%83%88%20(4%20%E3%83%90%E3%82%A4%E3%83%88%E3%81%AE%20UTF%2D,%E6%96%87%E5%AD%97%E3%82%92%E3%82%B5%E3%83%9D%E3%83%BC%E3%83%88%E3%81%97%E3%81%BE%E3%81%99%E3%80%82
  -- https://penpen-dev.com/blog/mysql-utf8-utf8mb4/
  -- utf8: 1~3バイトまで, utf8mb4: 1~4バイトまで対応
  -- 絵文字などを使用したい場合は，utf8mb4の指定が必要
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザー';

CREATE TABLE `tasks`
(
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'タスクの識別子',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'タスクを作成したユーザーの識別子 ',
  `title` VARCHAR(128) NOT NULL COMMENT 'タスクのタイトル',
  `status` VARCHAR(20) NOT NULL COMMENT 'タスクの状態',
  `created` DATETIME(6) NOT NULL COMMENT 'レコード作成日時',
  `modified` DATETIME(6) NOT NULL COMMENT 'レコード修正日時',
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_user_id`
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
      -- https://qiita.com/suin/items/21fe6c5a78c1505b19cb
      -- RESTRICT: update, delete どちらもエラーになる
      ON DELETE RESTRICT ON UPDATE RESTRICT
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='タスク';