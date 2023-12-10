package smecalculus.bezmen.configuration;

import static smecalculus.bezmen.configuration.MessagingDm.MappingMode.SPRING_MVC;
import static smecalculus.bezmen.configuration.MessagingDm.ProtocolMode.HTTP;

import java.util.Set;
import smecalculus.bezmen.configuration.MessagingDm.MappingProps;
import smecalculus.bezmen.configuration.MessagingDm.MessagingProps;
import smecalculus.bezmen.configuration.MessagingDm.ProtocolProps;

public abstract class MessagingDmEg {
    public static MessagingProps.Builder messagingProps() {
        return MessagingProps.builder()
                .protocolProps(protocolProps().build())
                .mappingProps(mappingProps().build());
    }

    public static ProtocolProps.Builder protocolProps() {
        return ProtocolProps.builder().protocolModes(Set.of(HTTP));
    }

    public static MappingProps.Builder mappingProps() {
        return MappingProps.builder().mappingModes(Set.of(SPRING_MVC));
    }
}
