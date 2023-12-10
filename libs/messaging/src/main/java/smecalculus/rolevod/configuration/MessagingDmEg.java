package smecalculus.rolevod.configuration;

import static smecalculus.rolevod.configuration.MessagingDm.MappingMode.SPRING_MVC;
import static smecalculus.rolevod.configuration.MessagingDm.ProtocolMode.HTTP;

import java.util.Set;
import smecalculus.rolevod.configuration.MessagingDm.MappingProps;
import smecalculus.rolevod.configuration.MessagingDm.MessagingProps;
import smecalculus.rolevod.configuration.MessagingDm.ProtocolProps;

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
