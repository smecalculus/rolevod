package smecalculus.bezmen.construction

import org.springframework.context.annotation.Bean
import smecalculus.bezmen.configuration.StorageDm
import smecalculus.bezmen.configuration.StorageDm.MappingMode.MY_BATIS
import smecalculus.bezmen.configuration.StorageDm.MappingMode.SPRING_DATA
import smecalculus.bezmen.configuration.StorageDm.ProtocolMode.H2
import smecalculus.bezmen.configuration.StorageDm.ProtocolMode.POSTGRES
import smecalculus.bezmen.configuration.StorageDmEg

class StoragePropsBeans {
    class SpringDataPostgres {
        @Bean
        fun storageProps(): StorageDm.StorageProps {
            return StorageDmEg.storageProps(SPRING_DATA, POSTGRES).build()
        }
    }

    class SpringDataH2 {
        @Bean
        fun storageProps(): StorageDm.StorageProps {
            return StorageDmEg.storageProps(SPRING_DATA, H2).build()
        }
    }

    class MyBatisPostgres {
        @Bean
        fun storageProps(): StorageDm.StorageProps {
            return StorageDmEg.storageProps(MY_BATIS, POSTGRES).build()
        }
    }

    class MyBatisH2 {
        @Bean
        fun storageProps(): StorageDm.StorageProps {
            return StorageDmEg.storageProps(MY_BATIS, H2).build()
        }
    }
}
