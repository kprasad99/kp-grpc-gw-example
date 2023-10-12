package io.github.kprasad99.example.service;

import java.util.stream.StreamSupport;

import org.lognet.springboot.grpc.GRpcService;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.stereotype.Service;

import com.google.protobuf.Empty;

import io.github.kprasad99.example.mapper.UserMapper;
import io.github.kprasad99.example.service.orm.dao.UserServiceDao;
import io.github.kprasad99.grpc.example.service.User;
import io.github.kprasad99.grpc.example.service.UserIDInput;
import io.github.kprasad99.grpc.example.service.UserServiceGrpc.UserServiceImplBase;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Service
@GRpcService
@AllArgsConstructor
public class KpUserService extends UserServiceImplBase {

	private final UserServiceDao dao;

	@Override
	public void list(Empty request, StreamObserver<User> responseObserver) {
		log.info("Listing users");
		try {
			StreamSupport.stream(dao.findAll().spliterator(), false).map(UserMapper.INSTANCE::toProto)
					.forEach(responseObserver::onNext);
			responseObserver.onCompleted();
		} catch (Exception e) {
			responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asException());
		}
	}

	@Override
	public void add(User request, StreamObserver<User> responseObserver) {
		log.info("Adding user {}", request.getName());
		try {
			var user = dao.save(UserMapper.INSTANCE.toDomain(request));
			responseObserver.onNext(UserMapper.INSTANCE.toProto(user));
			responseObserver.onCompleted();
		} catch (DataIntegrityViolationException e) {
			responseObserver.onError(Status.ALREADY_EXISTS.withDescription(e.getMessage()).withCause(e).asException());
		} catch (Exception e) {
			responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asException());
		}
	}

	@Override
	public void update(User request, StreamObserver<Empty> responseObserver) {
		log.info("updating user {}", request.getName());
		try {
			dao.save(UserMapper.INSTANCE.toDomain(request));
			responseObserver.onCompleted();
		} catch (DataIntegrityViolationException e) {
			responseObserver.onError(Status.ALREADY_EXISTS.withDescription(e.getMessage()).withCause(e).asException());
		} catch (Exception e) {
			responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asException());
		}
	}

	@Override
	public void remove(UserIDInput request, StreamObserver<Empty> responseObserver) {
		log.info("Removing user {}", request.getId());
		try {
			dao.deleteById(request.getId());
			responseObserver.onCompleted();
		} catch (DataIntegrityViolationException e) {
			responseObserver.onError(Status.ALREADY_EXISTS.withDescription(e.getMessage()).withCause(e).asException());
		} catch (Exception e) {
			responseObserver.onError(e);
		}
	}
}
