package smecalculus.rolevod.construction

import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Import
import org.springframework.jdbc.datasource.embedded.EmbeddedDatabaseBuilder
import org.springframework.jdbc.datasource.embedded.EmbeddedDatabaseType
import smecalculus.rolevod.configuration.StorageDm
import smecalculus.rolevod.configuration.StorageDm.ProtocolMode.H2
import smecalculus.rolevod.configuration.StorageDm.ProtocolMode.POSTGRES
import smecalculus.rolevod.storage.SepulkaDao
import smecalculus.rolevod.storage.SepulkaDaoMyBatis
import smecalculus.rolevod.storage.SepulkaDaoSpringData
import smecalculus.rolevod.storage.SepulkaStateMapper
import smecalculus.rolevod.storage.SepulkaStateMapperImpl
import smecalculus.rolevod.storage.mybatis.SepulkaSqlMapper
import smecalculus.rolevod.storage.springdata.SepulkaRepository
import smecalculus.rolevod.testing.Constants.CREATE_SQL
import smecalculus.rolevod.testing.Constants.DROP_SQL
import javax.sql.DataSource

class SepulkaDaoBeans {
    @Import(MappingSpringDataBeans::class)
    class SpringData {
        @Bean
        fun sepulkaDao(
            mapper: SepulkaStateMapper,
            repository: SepulkaRepository,
        ): SepulkaDao {
            return SepulkaDaoSpringData(mapper, repository)
        }
    }

    @Import(MappingMyBatisBeans::class)
    class MyBatis {
        @Bean
        fun sepulkaDao(
            stateMapper: SepulkaStateMapper,
            sqlMapper: SepulkaSqlMapper,
        ): SepulkaDao {
            return SepulkaDaoMyBatis(stateMapper, sqlMapper)
        }
    }

    class Anyone {
        @Bean
        fun sepulkaStateMapper(): SepulkaStateMapper {
            return SepulkaStateMapperImpl()
        }

        @Bean
        fun dataSource(storageProps: StorageDm.StorageProps): DataSource {
            val common = listOf("DB_CLOSE_DELAY=-1")
            val specific: List<String> =
                when (storageProps.protocolProps.protocolMode) {
                    H2 -> listOf("MODE=STRICT")
                    POSTGRES ->
                        listOf(
                            "MODE=PostgreSQL",
                            "DATABASE_TO_LOWER=TRUE",
                            "DEFAULT_NULL_ORDERING=HIGH",
                        )

                    else -> throw IllegalStateException("Unrepresentable state")
                }

            val nameWithSettings = listOf("test") + common + specific
            return EmbeddedDatabaseBuilder()
                .setType(EmbeddedDatabaseType.H2)
                .setName(nameWithSettings.joinToString(";"))
                .addScript(DROP_SQL)
                .addScript(CREATE_SQL)
                .build()
        }
    }
}
