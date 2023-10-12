package io.github.kprasad99.example.service.orm.dao;

import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import io.github.kprasad99.example.service.orm.domain.User;

@Repository
public interface UserServiceDao extends CrudRepository<User, Integer> {

}
