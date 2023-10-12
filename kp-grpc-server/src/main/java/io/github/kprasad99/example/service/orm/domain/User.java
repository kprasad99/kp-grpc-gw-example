package io.github.kprasad99.example.service.orm.domain;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import lombok.Data;

@Data
@Entity(name="tbl_user")
public class User {
	@Id
	private int id;
	private String name;
	private int age;
}
