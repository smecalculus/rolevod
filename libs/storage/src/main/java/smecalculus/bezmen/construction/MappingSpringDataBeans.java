package smecalculus.bezmen.construction;

import static smecalculus.bezmen.configuration.StorageDm.MappingMode.SPRING_DATA;

import java.util.Optional;
import javax.sql.DataSource;
import lombok.NonNull;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.jdbc.core.convert.JdbcCustomConversions;
import org.springframework.data.jdbc.core.dialect.JdbcPostgresDialect;
import org.springframework.data.jdbc.core.mapping.JdbcMappingContext;
import org.springframework.data.jdbc.repository.config.AbstractJdbcConfiguration;
import org.springframework.data.jdbc.repository.config.EnableJdbcRepositories;
import org.springframework.data.relational.RelationalManagedTypes;
import org.springframework.data.relational.core.dialect.Dialect;
import org.springframework.data.relational.core.dialect.H2Dialect;
import org.springframework.data.relational.core.mapping.NamingStrategy;
import org.springframework.jdbc.core.namedparam.NamedParameterJdbcOperations;
import org.springframework.jdbc.core.namedparam.NamedParameterJdbcTemplate;
import org.springframework.jdbc.datasource.DataSourceTransactionManager;
import org.springframework.transaction.PlatformTransactionManager;
import smecalculus.bezmen.configuration.StorageDm.StorageProps;

@ConditionalOnStorageMappingMode(SPRING_DATA)
@EnableJdbcRepositories("smecalculus.bezmen.storage.springdata")
@Configuration(proxyBeanMethods = false)
public class MappingSpringDataBeans extends AbstractJdbcConfiguration {

    @Autowired
    private StorageProps storageProps;

    @Bean
    public NamedParameterJdbcTemplate namedParameterJdbcTemplate(DataSource dataSource) {
        return new NamedParameterJdbcTemplate(dataSource);
    }

    @Bean
    public PlatformTransactionManager transactionManager(DataSource dataSource) {
        return new DataSourceTransactionManager(dataSource);
    }

    @Bean
    @Override
    public @NonNull Dialect jdbcDialect(@NonNull NamedParameterJdbcOperations operations) {
        return switch (storageProps.protocolProps().protocolMode()) {
            case H2 -> H2Dialect.INSTANCE;
            case POSTGRES -> JdbcPostgresDialect.INSTANCE;
        };
    }

    @Bean
    @Override
    public @NonNull JdbcMappingContext jdbcMappingContext(
            @NonNull Optional<NamingStrategy> namingStrategy,
            @NonNull JdbcCustomConversions customConversions,
            @NonNull RelationalManagedTypes jdbcManagedTypes) {
        var mappingContext = super.jdbcMappingContext(namingStrategy, customConversions, jdbcManagedTypes);
        mappingContext.setForceQuote(false);
        return mappingContext;
    }
}
