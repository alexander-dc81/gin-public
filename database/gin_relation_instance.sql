CREATE TABLE `gin`.`relation_instance` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `relation_id` INT NOT NULL,
  `parent_entity_instance_id` INT NOT NULL,
  `child_entity_instance_id` INT NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (`id`),
  INDEX `id` (`id` ASC) VISIBLE,
  INDEX `parent_entity_instance_entity_instance_idx` (`parent_entity_instance_id` ASC) VISIBLE,
  INDEX `child_entity_instance_entity_instance_idx` (`child_entity_instance_id` ASC) VISIBLE,
  CONSTRAINT `parent_entity_instance_entity_instance`
    FOREIGN KEY (`parent_entity_instance_id`)
    REFERENCES `gin`.`entity_instance` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `child_entity_instance_entity_instance`
    FOREIGN KEY (`child_entity_instance_id`)
    REFERENCES `gin`.`entity_instance` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);
