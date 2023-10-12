package io.github.kprasad99.example.configuration;

import javax.sql.DataSource;

import org.springframework.beans.BeansException;
import org.springframework.beans.factory.config.BeanPostProcessor;
import org.springframework.boot.actuate.autoconfigure.tracing.ConditionalOnEnabledTracing;
import org.springframework.context.annotation.Configuration;

import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.instrumentation.jdbc.datasource.OpenTelemetryDataSource;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Configuration
@ConditionalOnEnabledTracing
@AllArgsConstructor
@Slf4j
public class DatabaseSourceOtelConfiguration implements BeanPostProcessor {

	private final OpenTelemetry telemetry;

	
	@Override
	public Object postProcessBeforeInitialization(Object bean, String beanName) throws BeansException {
		if (bean instanceof DataSource b) {
			log.info("Bean type is {}", bean.getClass().getName());
			return new OpenTelemetryDataSource(b, telemetry);
		}
		return bean;
	}

}
