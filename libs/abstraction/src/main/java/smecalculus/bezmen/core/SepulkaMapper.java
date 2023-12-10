package smecalculus.bezmen.core;

import org.mapstruct.Mapper;
import smecalculus.bezmen.core.MessageDm.RegistrationRequest;
import smecalculus.bezmen.core.MessageDm.RegistrationResponse;
import smecalculus.bezmen.core.StateDm.AggregateState;

@Mapper
public interface SepulkaMapper {
    AggregateState.Builder toState(RegistrationRequest request);

    RegistrationResponse.Builder toMessage(AggregateState state);
}
