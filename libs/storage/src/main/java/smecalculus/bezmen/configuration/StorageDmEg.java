package smecalculus.bezmen.configuration;

import static smecalculus.bezmen.configuration.StorageDm.MappingMode.SPRING_DATA;
import static smecalculus.bezmen.configuration.StorageDm.ProtocolMode.H2;

import smecalculus.bezmen.configuration.StorageDm.H2Props;
import smecalculus.bezmen.configuration.StorageDm.MappingMode;
import smecalculus.bezmen.configuration.StorageDm.MappingProps;
import smecalculus.bezmen.configuration.StorageDm.PostgresProps;
import smecalculus.bezmen.configuration.StorageDm.ProtocolMode;
import smecalculus.bezmen.configuration.StorageDm.ProtocolProps;
import smecalculus.bezmen.configuration.StorageDm.StorageProps;

public abstract class StorageDmEg {
    public static StorageProps.Builder storageProps() {
        return StorageProps.builder()
                .protocolProps(protocolProps().build())
                .mappingProps(mappingProps().build());
    }

    public static StorageProps.Builder storageProps(MappingMode mappingMode, ProtocolMode protocolMode) {
        return storageProps()
                .protocolProps(protocolProps(protocolMode).build())
                .mappingProps(mappingProps(mappingMode).build());
    }

    public static MappingProps.Builder mappingProps() {
        return MappingProps.builder().mappingMode(SPRING_DATA);
    }

    public static MappingProps.Builder mappingProps(MappingMode mode) {
        return mappingProps().mappingMode(mode);
    }

    public static ProtocolProps.Builder protocolProps() {
        return ProtocolProps.builder()
                .protocolMode(H2)
                .h2Props(H2Props.builder()
                        .url("jdbc:h2:mem:bezmen;DB_CLOSE_DELAY=-1")
                        .username("sa")
                        .password("sa")
                        .build())
                .postgresProps(PostgresProps.builder()
                        .url("jdbc:postgresql://localhost:5432/bezmen")
                        .username("bezmen")
                        .password("bezmen")
                        .build());
    }

    public static ProtocolProps.Builder protocolProps(ProtocolMode mode) {
        return protocolProps().protocolMode(mode);
    }
}
