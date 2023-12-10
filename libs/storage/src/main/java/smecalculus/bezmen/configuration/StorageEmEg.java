package smecalculus.bezmen.configuration;

import static smecalculus.bezmen.configuration.StorageDm.MappingMode.SPRING_DATA;
import static smecalculus.bezmen.configuration.StorageDm.ProtocolMode.H2;

import smecalculus.bezmen.configuration.StorageEm.MappingProps;
import smecalculus.bezmen.configuration.StorageEm.ProtocolProps;
import smecalculus.bezmen.configuration.StorageEm.StorageProps;

public abstract class StorageEmEg {
    public static StorageProps storageProps() {
        var propsEdge = new StorageProps();
        propsEdge.setMapping(mappingProps());
        propsEdge.setProtocol(protocolProps());
        return propsEdge;
    }

    public static MappingProps mappingProps() {
        var propsEdge = new MappingProps();
        propsEdge.setMode(SPRING_DATA.name());
        return propsEdge;
    }

    public static ProtocolProps protocolProps() {
        var propsEdge = new ProtocolProps();
        propsEdge.setMode(H2.name());
        return propsEdge;
    }
}
