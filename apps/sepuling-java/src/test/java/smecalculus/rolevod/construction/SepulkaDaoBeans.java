package smecalculus.rolevod.construction;

import static java.util.stream.Collectors.joining;
import static smecalculus.rolevod.testing.Constants.CREATE_SQL;
import static smecalculus.rolevod.testing.Constants.DROP_SQL;

import java.util.Collection;
import java.util.List;
import java.util.stream.Stream;
import javax.sql.DataSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Import;
import org.springframework.jdbc.datasource.embedded.EmbeddedDatabaseBuilder;
import org.springframework.jdbc.datasource.embedded.EmbeddedDatabaseType;
import smecalculus.rolevod.configuration.StorageDm.StorageProps;
import smecalculus.rolevod.storage.SepulkaDao;
import smecalculus.rolevod.storage.SepulkaDaoMyBatis;
import smecalculus.rolevod.storage.SepulkaDaoSpringData;
import smecalculus.rolevod.storage.SepulkaStateMapper;
import smecalculus.rolevod.storage.SepulkaStateMapperImpl;
import smecalculus.rolevod.storage.mybatis.SepulkaSqlMapper;
import smecalculus.rolevod.storage.springdata.SepulkaRepository;

public class SepulkaDaoBeans {

    @Import(MappingSpringDataBeans.class)
    public static class SpringData {
        @Bean
        public SepulkaDao sepulkaDao(SepulkaStateMapper mapper, SepulkaRepository repository) {
            return new SepulkaDaoSpringData(mapper, repository);
        }
    }

    @Import(MappingMyBatisBeans.class)
    public static class MyBatis {
        @Bean
        public SepulkaDao sepulkaDao(SepulkaStateMapper stateMapper, SepulkaSqlMapper sqlMapper) {
            return new SepulkaDaoMyBatis(stateMapper, sqlMapper);
        }
    }

    public static class Anyone {
        @Bean
        public SepulkaStateMapper sepulkaStateMapper() {
            return new SepulkaStateMapperImpl();
        }

        @Bean
        public DataSource dataSource(StorageProps storageProps) {
            var common = List.of("DB_CLOSE_DELAY=-1");
            var specific =
                    switch (storageProps.protocolProps().protocolMode()) {
                        case H2 -> List.of("MODE=STRICT");
                        case POSTGRES -> List.of(
                                "MODE=PostgreSQL", "DATABASE_TO_LOWER=TRUE", "DEFAULT_NULL_ORDERING=HIGH");
                    };
            var nameWithSettings = Stream.of(List.of("test"), common, specific)
                    .flatMap(Collection::stream)
                    .collect(joining(";"));
            return new EmbeddedDatabaseBuilder()
                    .setType(EmbeddedDatabaseType.H2)
                    .setName(nameWithSettings)
                    .addScript(DROP_SQL)
                    .addScript(CREATE_SQL)
                    .build();
        }
    }
}
