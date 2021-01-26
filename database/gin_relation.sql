CREATE TABLE `gin`.`relation` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_entity_id` INT NOT NULL,
  `child_entity_id` INT NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `id` (`id` ASC) VISIBLE,
  INDEX `parent_entity_id_idx` (`parent_entity` ASC) VISIBLE,
  INDEX `child_entity_id_idx` (`child_entity` ASC) VISIBLE,
  CONSTRAINT `parent_entity_id`
    FOREIGN KEY (`parent_entity`)
    REFERENCES `gin`.`entity` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `child_entity_id`
    FOREIGN KEY (`child_entity`)
    REFERENCES `gin`.`entity` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);