package smecalculus.rolevod.core;

import org.mapstruct.Mapper;
import smecalculus.rolevod.core.SepulkaMessageDm.RegistrationRequest;
import smecalculus.rolevod.core.SepulkaMessageDm.RegistrationResponse;
import smecalculus.rolevod.core.SepulkaStateDm.AggregateRoot;

@Mapper
public interface SepulkaMapper {
    AggregateRoot.Builder toState(RegistrationRequest request);

    RegistrationResponse.Builder toMessage(AggregateRoot state);
}
