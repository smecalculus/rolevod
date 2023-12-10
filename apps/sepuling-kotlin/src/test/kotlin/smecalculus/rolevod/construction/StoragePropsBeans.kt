package smecalculus.rolevod.construction

import org.springframework.context.annotation.Bean
import smecalculus.rolevod.configuration.StorageDm
import smecalculus.rolevod.configuration.StorageDm.MappingMode.MY_BATIS
import smecalculus.rolevod.configuration.StorageDm.MappingMode.SPRING_DATA
import smecalculus.rolevod.configuration.StorageDm.ProtocolMode.H2
import smecalculus.rolevod.configuration.StorageDm.ProtocolMode.POSTGRES
import smecalculus.rolevod.configuration.StorageDmEg

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
