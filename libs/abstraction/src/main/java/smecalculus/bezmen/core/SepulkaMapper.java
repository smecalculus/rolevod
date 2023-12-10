package smecalculus.bezmen.core;

import org.mapstruct.Mapper;
import smecalculus.bezmen.core.SepulkaMessageDm.RegistrationRequest;
import smecalculus.bezmen.core.SepulkaMessageDm.RegistrationResponse;
import smecalculus.bezmen.core.SepulkaStateDm.AggregateRoot;

@Mapper
public interface SepulkaMapper {
    AggregateRoot.Builder toState(RegistrationRequest request);

    RegistrationResponse.Builder toMessage(AggregateRoot state);
}
