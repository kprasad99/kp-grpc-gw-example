package io.github.kprasad99.example;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest(properties = {"management.tracing.enabled=false"})
class KpGrpcGwExampleApplicationTests {

	@Test
	void contextLoads() {
	}

}
