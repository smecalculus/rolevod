package smecalculus.bezmen.construction;

import javax.sql.DataSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import org.springframework.jdbc.datasource.DriverManagerDataSource;
import smecalculus.bezmen.configuration.StorageDm.H2Props;
import smecalculus.bezmen.configuration.StorageDm.PostgresProps;
import smecalculus.bezmen.configuration.StorageDm.StorageProps;

@Import({StorageConfigBeans.class, MappingMyBatisBeans.class, MappingSpringDataBeans.class})
@Configuration(proxyBeanMethods = false)
public class StorageBeans {

    @Bean
    public DataSource dataSource(StorageProps storageProps) {
        var dataSource = new DriverManagerDataSource();
        var protocolProps = storageProps.protocolProps();
        switch (protocolProps.protocolMode()) {
            case H2 -> configure(dataSource, protocolProps.h2Props());
            case POSTGRES -> configure(dataSource, protocolProps.postgresProps());
        }
        return dataSource;
    }

    private void configure(DriverManagerDataSource dataSource, H2Props props) {
        dataSource.setDriverClassName("org.h2.Driver");
        dataSource.setUrl(props.url());
        dataSource.setUsername(props.username());
        dataSource.setPassword(props.password());
    }

    private void configure(DriverManagerDataSource dataSource, PostgresProps props) {
        dataSource.setDriverClassName("org.postgresql.Driver");
        dataSource.setUrl(props.url());
        dataSource.setUsername(props.username());
        dataSource.setPassword(props.password());
    }
}
